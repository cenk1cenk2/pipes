package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func HelmVersion(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand(
				"helm",
				"version",
			).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
