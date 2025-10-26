package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func GoLint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		Set(func(t *Task) error {
			t.CreateCommand(
				"golangci-lint",
				"run",
				"-v",
				"--timeout",
				P.Timeout.String(),
			).
				AppendEnvironment(setup.C.EnvVars).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
