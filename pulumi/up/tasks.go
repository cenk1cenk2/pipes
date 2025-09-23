package up

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
)

func PulumiUp(tl *TaskList) *Task {
	return tl.CreateTask("up").
		Set(func(t *Task) error {
			t.CreateCommand(
				"pulumi",
				"up",
				"--diff",
				"--yes",
				"-f",
			).
				SetDir(setup.P.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
