package pipe

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gitlab.kilic.dev/libraries/go-utils/utils"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type Ctx struct {
	Tags []string
}

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		ShouldRunBefore(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.Tags = []string{}

			return nil
		}).
		Set(func(t *Task[Pipe]) error {
			TagAsLatestForTagsRegex := []string{}
			TagAsLatestForBranchesRegex := []string{}

			// add all the specified tags
			for _, v := range utils.RemoveDuplicateStr(utils.DeleteEmptyStringsFromSlice(t.Pipe.DockerImage.Tags.Value())) {
				if err := AddDockerTag(v); err != nil {
					return err
				}
			}

			// add tags through tags file
			if _, err := os.Stat(t.Pipe.DockerImage.TagsFile); err == nil {
				t.Log.Infof(
					"Tags file does exists, will overwrite current tags: %s",
					t.Pipe.DockerImage.TagsFile,
				)

				content, err := ioutil.ReadFile(t.Pipe.DockerImage.TagsFile)
				if err != nil {
					return err
				}

				tags := strings.Split(string(content), ",")

				t.Pipe.Ctx.Tags = []string{}

				re := regexp.MustCompile(`\r?\n`)

				for _, v := range tags {
					func(v string) {
						t.CreateSubtask("").
							Set(func(t *Task[Pipe]) error {
								return AddDockerTag(re.ReplaceAllString(v, ""))
							}).
							AddSelfToParent(func(pt, st *Task[Pipe]) {
								pt.ExtendSubtask(func(j Job) Job {
									return tl.JobParallel(j, st.Job())
								})
							})
					}(v)
				}

				if err = t.RunSubtasks(); err != nil {
					return err
				}
			} else if errors.Is(err, os.ErrNotExist) {
				if t.Pipe.DockerImage.TagsFile != "" {
					t.Log.Warnf("Tags file is set but it does not exists: %s", t.Pipe.DockerImage.TagsFile)

					t.SendExit(0)

					return nil
				} else {
					t.Log.Debugf("Tags file is not specified: %s", t.Pipe.DockerImage.TagsFile)
				}
			} else {
				t.Log.Warnf("Can not read the tags file: %s", t.Pipe.DockerImage.TagsFile)
			}

			// tag as latest for tags
			if t.Pipe.DockerImage.TagAsLatestForTagsRegex != "" {
				err := json.Unmarshal(
					[]byte(t.Pipe.DockerImage.TagAsLatestForTagsRegex),
					&TagAsLatestForTagsRegex,
				)

				if err != nil {
					return err
				}

				if t.Pipe.Git.Tag != "" {
					for _, re := range TagAsLatestForTagsRegex {
						m, err := regexp.Match(re, []byte(t.Pipe.Git.Tag))

						if err != nil {
							return err
						}

						if m {
							if err := AddDockerTag(DOCKER_LATEST_TAG); err != nil {
								return err
							}

							t.Log.Infof(
								"Will tag image as latest since tag regex matches: %s",
								re,
							)

							break
						}
					}
				}
			}

			// tag as latest for branches
			if t.Pipe.DockerImage.TagAsLatestForBranchesRegex != "" {
				err := json.Unmarshal(
					[]byte(t.Pipe.DockerImage.TagAsLatestForBranchesRegex),
					&TagAsLatestForBranchesRegex,
				)

				if err != nil {
					return err
				}

				if t.Pipe.Git.Branch != "" {
					for _, re := range TagAsLatestForBranchesRegex {
						m, err := regexp.Match(re, []byte(t.Pipe.Git.Branch))

						if err != nil {
							return err
						}

						if m {
							if err := AddDockerTag(DOCKER_LATEST_TAG); err != nil {
								return err
							}

							t.Log.Infof(
								"Will tag image as latest since branch regex matches: %s",
								re,
							)

							break
						}
					}
				}
			}

			t.Pipe.Ctx.Tags = utils.RemoveDuplicateStr(
				utils.DeleteEmptyStringsFromSlice(t.Pipe.Ctx.Tags),
			)

			t.Log.Infof(
				"Image tags: %s", strings.Join(t.Pipe.Ctx.Tags, ", "),
			)

			return nil
		})
}

func DockerVersion(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(DOCKER_EXE, "--version").
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			if t.Pipe.Docker.UseBuildx {
				t.Log.Infoln("Docker Buildx is enabled.")

				t.CreateCommand(DOCKER_EXE, "buildx", "version").
					SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}

func DockerLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		Set(func(t *Task[Pipe]) error {
			// login task
			t.CreateSubtask("login").
				ShouldDisable(func(t *Task[Pipe]) bool {
					return t.Pipe.DockerRegistry.Username == "" ||
						t.Pipe.DockerRegistry.Password == ""
				}).
				Set(func(t *Task[Pipe]) error {
					t.CreateCommand(
						DOCKER_EXE,
						"login",
						t.Pipe.DockerRegistry.Registry,
						"--username",
						t.Pipe.DockerRegistry.Username,
						"--password-stdin",
					).
						SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
						Set(func(c *Command[Pipe]) error {
							c.Command.Stdin = strings.NewReader(t.Pipe.DockerRegistry.Password)

							c.Log.Infof(
								"Logging in to Docker registry: %s",
								t.Pipe.DockerRegistry.Registry,
							)

							return nil
						}).
						AddSelfToTheTask()

					return nil
				}).
				ShouldRunAfter(func(t *Task[Pipe]) error {
					return t.RunCommandJobAsJobSequence()
				}).
				AddSelfToParent(func(pt, st *Task[Pipe]) {
					pt.ExtendSubtask(func(job Job) Job {
						return tl.JobSequence(job, st.Job())
					})
				})

				// login verify task
			t.CreateSubtask("login:verify").
				ShouldDisable(func(t *Task[Pipe]) bool {
					return t.Pipe.DockerRegistry.Username != "" &&
						t.Pipe.DockerRegistry.Password != ""
				}).
				Set(func(t *Task[Pipe]) error {
					t.CreateCommand(
						DOCKER_EXE,
						"login",
						t.Pipe.DockerRegistry.Registry,
					).
						SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
						Set(func(c *Command[Pipe]) error {
							c.Log.Debugf(
								"Will verify authentication in to Docker registry: %s",
								t.Pipe.DockerRegistry.Registry,
							)

							return nil
						}).
						AddSelfToTheTask()

					return nil
				}).
				ShouldRunAfter(func(t *Task[Pipe]) error {
					return t.RunCommandJobAsJobSequence()
				}).
				AddSelfToParent(func(pt, st *Task[Pipe]) {
					pt.ExtendSubtask(func(job Job) Job {
						return tl.JobSequence(job, st.Job())
					})
				})

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func DockerSetupBuildX(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("setup:buildx").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Docker.UseBuildx
		}).
		Set(func(t *Task[Pipe]) error {
			err := t.CreateCommand(
				DOCKER_EXE,
				"buildx",
				"create",
				"--use",
				"--name",
				"gitlab",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				Set(func(c *Command[Pipe]) error {
					c.Log.Infoln("Creating a new instance of docker buildx.")

					return nil
				}).Run()

			if err != nil {
				err := t.CreateCommand(
					DOCKER_EXE,
					"buildx",
					"use",
					"gitlab",
				).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
					Set(func(c *Command[Pipe]) error {
						c.Log.Warnln(
							"Creating a new docker buildx instance failed, trying to use the existing one.",
						)

						return nil
					}).Run()

				if err != nil {
					return err
				}
			}

			return nil
		})
}

func DockerBuildX(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Docker.UseBuildx
		}).
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				t.Pipe.DockerFile.Name,
				t.Pipe.DockerFile.Context,
			)

			t.Log.Infoln("Using Docker Buildx for building the Docker image.")

			// spawn virtual machine
			t.CreateCommand(
				DOCKER_EXE,
				"run",
				"--rm",
				"--privileged",
				"multiarch/qemu-user-static",
				"--reset",
				"-p",
				"yes",
			).
				AddSelfToTheTask()

				// check virtual machine
			t.CreateCommand(
				DOCKER_EXE,
				"buildx",
				"inspect",
				"--bootstrap",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

				// build image
			t.CreateCommand(
				DOCKER_EXE,
				"buildx",
				"build",
			).
				Set(func(c *Command[Pipe]) error {
					for _, v := range t.Pipe.DockerImage.BuildArgs.Value() {
						c.AppendArgs("--build-arg", v)
					}

					if t.Pipe.DockerImage.Pull {
						c.AppendArgs("--pull")
					}

					c.AppendArgs("--push")

					if t.Pipe.Docker.BuildxPlatforms != "" {
						c.AppendArgs("--platform", t.Pipe.Docker.BuildxPlatforms)
					}

					for _, tag := range t.Pipe.Ctx.Tags {
						c.AppendArgs("-t", tag)
					}

					c.AppendArgs(
						"--file",
						t.Pipe.DockerFile.Name,
						".",
					)

					c.SetDir(t.Pipe.DockerFile.Context)
					t.Log.Debugf("CWD set as: %s", c.Command.Dir)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerBuild(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("build").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Docker.UseBuildx
		}).
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				t.Pipe.DockerFile.Name,
				t.Pipe.DockerFile.Context,
			)

			// build image
			t.CreateCommand(
				DOCKER_EXE,
				"build",
			).
				Set(func(c *Command[Pipe]) error {
					for _, v := range t.Pipe.DockerImage.BuildArgs.Value() {
						c.AppendArgs("--build-arg", v)
					}

					if t.Pipe.DockerImage.Pull {
						c.AppendArgs("--pull")
					}

					for _, tag := range t.Pipe.Ctx.Tags {
						c.AppendArgs("-t", tag)
					}

					c.AppendArgs(
						"--file",
						t.Pipe.DockerFile.Name,
						".",
					)

					c.SetDir(t.Pipe.DockerFile.Context)
					t.Log.Debugf("CWD set as: %s", c.Command.Dir)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerPush(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("push").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Docker.UseBuildx
		}).
		Set(func(t *Task[Pipe]) error {
			for _, tag := range t.Pipe.Ctx.Tags {
				func(tag string) {
					t.CreateCommand(
						DOCKER_EXE,
						"push",
						tag,
					).
						Set(func(c *Command[Pipe]) error {
							c.Log.Infof(
								"Pushing Docker image: %s",
								tag,
							)

							return nil
						}).
						AddSelfToTheTask()
				}(tag)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}

func DockerInspect(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("push").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.DockerImage.Inspect
		}).
		Set(func(t *Task[Pipe]) error {
			for _, tag := range t.Pipe.Ctx.Tags {
				func(tag string) {
					t.CreateCommand(
						DOCKER_EXE,
						"manifest",
						"inspect",
						tag,
					).
						SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
						Set(func(c *Command[Pipe]) error {
							c.Log.Infof(
								"Inspecting Docker image: %s",
								tag,
							)

							return nil
						}).
						AddSelfToTheTask()
				}(tag)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}
