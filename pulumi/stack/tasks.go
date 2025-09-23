package stack

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
)

func PulumiSelectStack(tl *TaskList) *Task {
	return tl.CreateTask("stack").
		Set(func(t *Task) error {
			t.CreateCommand(
				"pulumi",
				"stack",
				"select",
			).
				Set(func(c *Command) error {
					c.AppendArgs(P.Stack)

					return nil
				}).
				SetDir(setup.P.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
