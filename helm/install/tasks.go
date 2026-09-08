package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func install(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(t *Task) error {
			t.CreateCommand(
				"helm",
				"dependency",
				"update",
			).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
