package build

import (
	"context"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"
)

func build(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(_ context.Context, t *Task) error {
			t.Log.Info(fmt.Sprintf("Building container image: %s in %s",
				P.File.Name,
				P.File.Context),
			)

			// build image
			t.CreateCommand(
				"buildah",
				"build",
			).
				SetDir(P.File.Context).
				Set(func(_ context.Context, c *Command) error {
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
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func push(tl *TaskList) *Task {
	return tl.CreateTask("push").
		ShouldDisable(func(t *Task) bool {
			return !P.Image.Push
		}).
		Set(func(ctx context.Context, t *Task) error {
			for _, tag := range C.Tags {
				t.CreateSubtask(tag).
					Set(func(_ context.Context, t *Task) error {
						t.CreateCommand(
							"buildah",
							"push",
							tag,
						).
							Set(func(_ context.Context, c *Command) error {
								t.Log.Info(fmt.Sprintf("Pushing container image: %s",
									tag),
								)

								c.AppendEnvironment(map[string]string{
									"STORAGE_DRIVER": P.Image.StorageDriver,
								})

								return nil
							}).
							SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(ctx context.Context, t *Task) error {
						return t.RunCommandJobAsJobParallel(ctx)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}
