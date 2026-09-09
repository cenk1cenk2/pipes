package lint

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

// Only one of the lint tasks fires, since the gates below are the two halves of
// the workspace condition the setup resolved.
func lint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		SetJobWrapper(func(_ Job, t *Task) Job {
			return JobParallel(
				lintModule(tl).Job(),
				lintWorkspace(tl).Job(),
			)
		})
}

func lintModule(tl *TaskList) *Task {
	return tl.CreateTask("lint", "module").
		ShouldDisable(func(_ *Task) bool {
			return setup.C.Workspace
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"golangci-lint",
				"run",
				"-v",
				"--timeout",
				P.Timeout.String(),
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				SetDir(setup.C.Cwd).
				Set(func(c *Command) error {
					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

// Every module is linted from inside its own directory, since the go tool
// resolves a "<module>/..." pattern to nothing for an underscored directory.
func lintWorkspace(tl *TaskList) *Task {
	return tl.CreateTask("lint", "workspace").
		ShouldDisable(func(_ *Task) bool {
			return !setup.C.Workspace
		}).
		Set(func(t *Task) error {
			for _, module := range setup.C.Modules {
				t.CreateCommand(
					"golangci-lint",
					"run",
					"-v",
					"--timeout",
					P.Timeout.String(),
				).
					SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
					SetDir(module).
					Set(func(c *Command) error {
						if P.Args != "" {
							c.AppendArgs(strings.Split(P.Args, " ")...)
						}

						c.AppendArgs("./...")

						return nil
					}).
					AppendEnvironment(setup.C.Env).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
