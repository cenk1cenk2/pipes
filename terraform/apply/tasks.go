package apply

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/shell"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func apply(tl *TaskList) *Task {
	return tl.CreateTask("apply").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"terraform",
				"apply",
				"-input=false",
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Apply.Output != "" {
						c.AppendArgs(P.Apply.Output)
					}

					if P.Apply.Args != "" {
						args, err := shell.Split(P.Apply.Args)
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
			return C.Report.Apply(ctx, t, t.RunCommandJobAsJobSequence)
		})
}
