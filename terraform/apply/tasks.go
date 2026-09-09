package apply

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
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
						c.AppendArgs(P.Apply.Args)
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
