package build

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/shell"
	"gitlab.kilic.dev/devops/pipes/python/setup"
)

func build(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"uv",
				"build",
				"--out-dir",
				P.Output,
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Args != "" {
						args, err := shell.Split(P.Args)
						if err != nil {
							return err
						}

						c.AppendArgs(args...)
					}

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
