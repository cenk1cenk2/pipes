package stack

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func PulumiSelectStack(tl *TaskList, deps Deps) *Task {
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
				SetDir(deps.Tool.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
