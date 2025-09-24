package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
)

func PulumiPlan(tl *TaskList) *Task {
	return tl.CreateTask("plan").
		Set(func(t *Task) error {
			t.CreateCommand(
				"pulumi",
				"preview",
				"--non-interactive",
				"--diff",
				"--save-plan",
				P.Plan,
			).
				SetDir(setup.P.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
