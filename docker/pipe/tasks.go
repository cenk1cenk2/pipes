package pipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
	u "gitlab.kilic.dev/libraries/go-utils/utils"
)

type Ctx struct {
	Tags []string
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify"},
		Task: func(t *utils.Task) error {
			Context.Tags = []string{}
			TagAsLatestForTagsRegex := []string{}
			TagAsLatestForBranchesRegex := []string{}

			for _, v := range u.RemoveDuplicateStr(u.DeleteEmptyStringsFromSlice(Pipe.DockerImage.Tags.Value())) {
				tag, err := AddDockerTag(v)

				if err != nil {
					return err
				}

				t.Log.Debugf("Docker image tag: %s",
					tag,
				)
			}

			if Pipe.DockerImage.TagAsLatestForTagsRegex != "" {
				err := json.Unmarshal(
					[]byte(Pipe.DockerImage.TagAsLatestForTagsRegex),
					&TagAsLatestForTagsRegex,
				)

				if err != nil {
					return err
				}

				if Pipe.Git.Tag != "" {
					for _, re := range TagAsLatestForTagsRegex {
						m, err := regexp.Match(re, []byte(Pipe.Git.Tag))

						if err != nil {
							return err
						}

						if m {
							tag, err := AddDockerTag(DOCKER_LATEST_TAG)

							if err != nil {
								return err
							}

							t.Log.Infof(
								"Will tag image as latest since tag regex matches: %s",
								tag,
							)

							break
						}
					}
				}
			}

			if Pipe.DockerImage.TagAsLatestForBranchesRegex != "" {
				err := json.Unmarshal(
					[]byte(Pipe.DockerImage.TagAsLatestForBranchesRegex),
					&TagAsLatestForBranchesRegex,
				)

				if err != nil {
					return err
				}

				if Pipe.Git.Branch != "" {
					for _, re := range TagAsLatestForBranchesRegex {
						m, err := regexp.Match(re, []byte(Pipe.Git.Branch))

						if err != nil {
							return err
						}

						if m {
							tag, err := AddDockerTag(DOCKER_LATEST_TAG)

							if err != nil {
								return err
							}

							t.Log.Infof(
								"Will tag image as latest since branch regex matches: %s",
								tag,
							)

							break
						}
					}
				}
			}

			if _, err := os.Stat(Pipe.DockerImage.TagsFile); err == nil {
				t.Log.Infof(
					"Tags file does exists, will override: %s",
					Pipe.DockerImage.TagsFile,
				)

				content, err := ioutil.ReadFile(Pipe.DockerImage.TagsFile)
				if err != nil {
					return err
				}

				tags := strings.Split(string(content), ",")

				Context.Tags = []string{}

				re := regexp.MustCompile(`\r?\n`)

				errs := []error{}

				for _, v := range tags {
					tag, err := AddDockerTag(re.ReplaceAllString(v, ""))

					if err != nil {
						errs = append(errs, err)
					}

					Context.Tags = append(Context.Tags, tag)
				}

				if len(errs) > 0 {
					for _, v := range errs {
						t.Log.Errorln(v)
					}

					return fmt.Errorf("Errors encountered while injecting environment variables.")
				}
			} else if errors.Is(err, os.ErrNotExist) {
				if Pipe.DockerImage.TagsFile != "" {
					t.Log.Warnf("Tags file is set but it does not exists: %s", Pipe.DockerImage.TagsFile)

					t.Log.Warnln("Nothing to do. Exitting...")
					os.Exit(0)
				} else {
					t.Log.Debugf("Tags file does not exists: %s", Pipe.DockerImage.TagsFile)
				}
			} else {
				t.Log.Warnf("Can not read the tags file: %s", Pipe.DockerImage.TagsFile)
			}

			Context.Tags = u.RemoveDuplicateStr(
				u.DeleteEmptyStringsFromSlice(Context.Tags),
			)

			t.Log.Infof(
				"Image tags: %s", strings.Join(Context.Tags, ", "),
			)

			return nil
		},
	}
}

func DockerVersion() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "version"},
		Task: func(t *utils.Task) error {
			cmd := exec.Command(DOCKER_EXE, "--version")

			t.Commands = append(t.Commands, cmd)

			if Pipe.Docker.UseBuildx {
				t.Log.Infoln("Docker Buildx is enabled.")

				cmd := exec.Command(DOCKER_EXE, "buildx", "version")

				t.Commands = append(t.Commands, cmd)
			}

			return nil
		},
	}

}

func DockerLogin() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{
			Context:        "login",
			StdOutLogLevel: logrus.DebugLevel,
		},
		Task: func(t *utils.Task) error {
			if Pipe.DockerRegistry.Username != "" && Pipe.DockerRegistry.Password != "" {
				t.Log.Infof(
					"Logging in to Docker registry: %s", Pipe.DockerRegistry.Registry,
				)

				cmd := exec.Command(DOCKER_EXE, "login")

				cmd.Args = append(cmd.Args, Pipe.DockerRegistry.Registry)
				cmd.Args = append(
					cmd.Args,
					"--username",
					Pipe.DockerRegistry.Username,
					"--password-stdin",
				)

				cmd.Stdin = strings.NewReader(Pipe.DockerRegistry.Password)

				t.Commands = append(t.Commands, cmd)
			}

			t.Log.Debugf(
				"Will verify authentication in to Docker registry: %s",
				Pipe.DockerRegistry.Registry,
			)

			cmd := exec.Command(DOCKER_EXE, "login")

			cmd.Args = append(cmd.Args, Pipe.DockerRegistry.Registry)

			t.Commands = append(t.Commands, cmd)

			return nil
		},
	}
}

func DockerSetupBuildx() utils.Task {
	metadata := utils.TaskMetadata{Context: "setup-buildx", Skip: !Pipe.Docker.UseBuildx}

	return utils.Task{
		Metadata: metadata,
		Task: func(t *utils.Task) error {
			t.Log.Infoln("Creating a new instance of docker buildx.")

			cmd := exec.Command(DOCKER_EXE, "buildx", "create", "--use", "--name", "gitlab")

			err := cmd.Run()

			if err != nil {
				t.Log.Warnln(
					"Creating a new docker buildx instance failed, trying to use the existing one.",
				)
				cmd = exec.Command(DOCKER_EXE, "buildx", "use", "gitlab")

				err := cmd.Run()

				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func DockerBuildx() utils.Task {
	metadata := utils.TaskMetadata{Context: "buildx", Skip: !Pipe.Docker.UseBuildx}

	return utils.Task{
		Metadata: metadata,
		Task: func(t *utils.Task) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				Pipe.DockerFile.Name,
				Pipe.DockerFile.Context,
			)

			t.Log.Infoln("Using Docker Buildx for building the Docker image.")

			cmd := exec.Command(
				DOCKER_EXE,
				"run",
				"--rm",
				"--privileged",
				"multiarch/qemu-user-static",
				"--reset",
				"-p",
				"yes",
			)

			t.Commands = append(t.Commands, cmd)

			cmd = exec.Command(DOCKER_EXE, "buildx", "inspect", "--bootstrap")

			t.Commands = append(t.Commands, cmd)

			cmd = exec.Command(DOCKER_EXE, "buildx", "build")

			for _, v := range Pipe.DockerImage.BuildArgs.Value() {
				cmd.Args = append(cmd.Args, "--build-arg", v)
			}

			if Pipe.DockerImage.Pull {
				cmd.Args = append(cmd.Args, "--pull")
			}

			cmd.Args = append(cmd.Args, "--push")

			if Pipe.Docker.BuildxPlatforms != "" {
				cmd.Args = append(cmd.Args, "--platform", Pipe.Docker.BuildxPlatforms)
			}

			for _, tag := range Context.Tags {
				cmd.Args = append(cmd.Args, "-t", tag)
			}

			cmd.Dir = Pipe.DockerFile.Context
			t.Log.Debugf("CWD set as: %s", cmd.Dir)

			cmd.Args = append(
				cmd.Args,
				"--file",
				Pipe.DockerFile.Name,
				".",
			)

			t.Commands = append(t.Commands, cmd)

			return nil
		},
	}
}

func DockerBuild() utils.Task {
	return utils.Task{Metadata: utils.TaskMetadata{Context: "build", Skip: Pipe.Docker.UseBuildx},
		Task: func(t *utils.Task) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				Pipe.DockerFile.Name,
				Pipe.DockerFile.Context,
			)

			cmd := exec.Command(DOCKER_EXE, "build")

			for _, v := range Pipe.DockerImage.BuildArgs.Value() {
				cmd.Args = append(cmd.Args, "--build-arg", v)
			}

			if Pipe.DockerImage.Pull {
				cmd.Args = append(cmd.Args, "--pull")
			}

			for _, tag := range Context.Tags {
				cmd.Args = append(cmd.Args, "-t", tag)
			}

			cmd.Dir = Pipe.DockerFile.Context
			t.Log.Debugf("CWD set as: %s", cmd.Dir)

			cmd.Args = append(
				cmd.Args,
				"--file",
				Pipe.DockerFile.Name,
				".",
			)

			t.Commands = append(t.Commands, cmd)

			return nil
		}}
}

func DockerPush() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "push", Skip: Pipe.Docker.UseBuildx},
		Task: func(t *utils.Task) error {
			t.Commands = []utils.Command{}

			for _, tag := range Context.Tags {
				t.Log.Infof(
					"Pushing Docker image: %s",
					tag,
				)

				cmd := exec.Command(DOCKER_EXE, "push")

				cmd.Args = append(cmd.Args, tag)

				t.Commands = append(t.Commands, cmd)
			}

			return nil
		},
	}
}

func DockerInspect() utils.Task {
	return utils.Task{Metadata: utils.TaskMetadata{
		Context:        "inspect",
		Skip:           !Pipe.DockerImage.Inspect,
		StdOutLogLevel: logrus.DebugLevel,
	},
		Task: func(t *utils.Task) error {
			t.Commands = []utils.Command{}

			for _, tag := range Context.Tags {
				t.Log.Infof(
					"Inspecting Docker image: %s",
					tag,
				)

				cmd := exec.Command(DOCKER_EXE, "manifest", "inspect")

				cmd.Args = append(cmd.Args, tag)

				t.Commands = append(t.Commands, cmd)
			}

			return nil
		}}

}
