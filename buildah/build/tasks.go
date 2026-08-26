package build

import (
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
)

func ContainerBuild(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(t *Task) error {
			t.Log.Infof(
				"Building container image: %s in %s",
				P.File.Name,
				P.File.Context,
			)

			// build image
			t.CreateCommand(
				"buildah",
				"build",
			).
				SetDir(P.File.Context).
				Set(func(c *Command) error {
					c.AppendEnvironment(map[string]string{
						"STORAGE_DRIVER": P.Image.StorageDriver,
					})

					c.AppendArgs("--format", P.Image.Format)

					if P.Image.Cache != "" {
						c.AppendArgs(
							"--layers",
							"--cache-from",
							P.Image.Cache,
							"--cache-to",
							P.Image.Cache,
						)
					}

					for k, t := range P.Image.BuildArgs {
						v, err := InlineTemplate[any](t, nil)
						if err != nil {
							return fmt.Errorf("Cannot process build argument template for %s: %w", k, err)
						}

						c.AppendArgs("--build-arg", fmt.Sprintf("%s=%s", k, v))
					}

					if P.Image.Pull {
						c.AppendArgs("--pull")
					}

					for _, tag := range C.Tags {
						c.AppendArgs("-t", tag)
					}

					c.AppendArgs(
						"--file",
						P.File.Name,
						".",
					)

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
			return !P.Image.Push
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

								c.AppendEnvironment(map[string]string{
									"STORAGE_DRIVER": P.Image.StorageDriver,
								})

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
