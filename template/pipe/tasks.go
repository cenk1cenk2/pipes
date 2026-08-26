package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func DefaultTask(tl *TaskList) *Task {
	return tl.CreateTask("default").
		Set(func(t *Task) error {
			t.CreateCommand(
				"echo",
			).
				Set(func(c *Command) error {
					c.AppendArgs("hello")

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
