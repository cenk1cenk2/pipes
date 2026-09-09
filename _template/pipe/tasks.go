package pipe

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

func entrypoint(tl *TaskList) *Task {
	return tl.CreateTask("entrypoint").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"echo",
			).
				Set(func(_ context.Context, c *Command) error {
					c.AppendArgs("hello")

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
