package up

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
)

func up(tl *TaskList) *Task {
	return tl.CreateTask("up").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"pulumi",
				"up",
				"--diff",
				"--yes",
				"-f",
				"--plan",
				P.Plan,
			).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
