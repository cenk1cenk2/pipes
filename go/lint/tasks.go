package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func GoLint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		Set(func(t *Task) error {
			// A workspace is linted one module at a time from inside it. The go tool
			// drops directories whose name starts with an underscore out of a package
			// pattern, so a "<module>/..." argument resolves to nothing at all for the
			// scaffold and it would go unlinted without a word from the linter.
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
