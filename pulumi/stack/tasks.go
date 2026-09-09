package stack

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
)

func stack(tl *TaskList) *Task {
	return tl.CreateTask("stack").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"pulumi",
				"stack",
				"select",
			).
				Set(func(_ context.Context, c *Command) error {
					c.AppendArgs(P.Stack)

					return nil
				}).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
