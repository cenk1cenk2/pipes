package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func GoLint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		Set(func(t *Task) error {
			// linted from inside each module, since the go tool resolves a "<module>/..."
			// pattern to nothing for an underscored directory.
			for _, module := range setup.C.Modules {
				lint(t).
					SetDir(module).
					AppendArgs("./...").
					AddSelfToTheTask()
			}

			if len(setup.C.Modules) == 0 {
				lint(t).
					SetDir(setup.C.Cwd).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func lint(t *Task) *Command {
	return t.CreateCommand(
		"golangci-lint",
		"run",
		"-v",
		"--timeout",
		P.Timeout.String(),
	).
		AppendEnvironment(setup.C.Env).
		SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG)
}
