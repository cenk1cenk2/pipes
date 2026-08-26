package up

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func PulumiUp(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("up").
		Set(func(t *Task) error {
			t.CreateCommand(
				"pulumi",
				"up",
				"--diff",
				"--yes",
				"-f",
				"--plan",
				P.Plan,
			).
				SetDir(deps.Tool.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
