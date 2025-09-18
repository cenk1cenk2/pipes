package build

import (
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/parser"
)

func ParseReferences(tl *TaskList) *Task {
	return tl.CreateTask("init", "references").
		Set(func(t *Task) error {
			C.References = parser.ParseGitReferences(P.Git.Tag, P.Git.Branch)

			t.Log.Debugf("References for environment selection: %v", C.References)

			return nil
		})
}

func ContainerBuild(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(t *Task) error {
			t.Log.Infof(
				"Building container image: %s in %s",
				P.ContainerFile.Name,
				P.ContainerFile.Context,
			)

			// build image
			t.CreateCommand(
				"buildah",
				"build",
			).
				Set(func(c *Command) error {
					c.AppendEnvironment(map[string]string{
						"STORAGE_DRIVER": P.ContainerImage.StorageDriver,
					})

					c.AppendArgs("--format", P.ContainerImage.Format)

					if P.ContainerImage.Cache != "" {
						c.AppendArgs(
							"--layers",
							"--cache-from",
							P.ContainerImage.Cache,
							"--cache-to",
							P.ContainerImage.Cache,
						)
					}

					for k, t := range P.ContainerImage.BuildArgs {
						v, err := InlineTemplate[any](t, nil)
						if err != nil {
							return fmt.Errorf("Cannot process build argument template for %s: %w", k, err)
						}

						c.AppendArgs("--build-arg", fmt.Sprintf("%s=%s", k, v))
					}

					if P.ContainerImage.Pull {
						c.AppendArgs("--pull")
					}

					for _, tag := range C.Tags {
						c.AppendArgs("-t", tag)
					}

					c.AppendArgs(
						"--file",
						P.ContainerFile.Name,
						".",
					)

					c.SetDir(P.ContainerFile.Context)
					t.Log.Debugf("CWD set as: %s", c.Command.Dir)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func ContainerPush(tl *TaskList) *Task {
	return tl.CreateTask("push").
		ShouldDisable(func(t *Task) bool {
			return !P.ContainerImage.Push
		}).
		Set(func(t *Task) error {
			for _, tag := range C.Tags {
				t.CreateSubtask(tag).
					Set(func(t *Task) error {
						t.CreateCommand(
							"buildah",
							"push",
							tag,
						).
							Set(func(c *Command) error {
								t.Log.Infof(
									"Pushing container image: %s",
									tag,
								)

								return nil
							}).
							SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task) error {
						return t.RunCommandJobAsJobParallel()
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

func ContainerInspect(tl *TaskList) *Task {
	return tl.CreateTask("inspect").
		ShouldDisable(func(t *Task) bool {
			return !P.ContainerImage.Inspect
		}).
		Set(func(t *Task) error {
			for _, tag := range C.Tags {
				t.CreateSubtask(tag).
					Set(func(t *Task) error {
						t.CreateCommand(
							"buildah",
							"manifest",
							"inspect",
							tag,
						).
							SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
							Set(func(c *Command) error {
								c.Log.Infof(
									"Inspecting container image: %s",
									tag,
								)

								return nil
							}).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task) error {
						return t.RunCommandJobAsJobParallel()
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}
