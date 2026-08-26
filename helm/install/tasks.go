package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func HelmInstall(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("install").
		Set(func(t *Task) error {
			t.CreateCommand(
				"helm",
				"dependency",
				"update",
			).
				SetDir(deps.Tool.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
