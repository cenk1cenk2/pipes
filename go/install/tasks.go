package install

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func GoModVendor(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("vendor").
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				SetDir(deps.Tool.Cwd).
				Set(func(c *Command) error {
					if deps.Tool.Workspace {
						c.AppendArgs("work", "vendor")

						t.Log.Infof("Vendoring workspace: in %s", deps.Tool.Cwd)
					} else {
						c.AppendArgs("mod", "vendor")

						t.Log.Infof("Vendoring: in %s", deps.Tool.Cwd)
					}

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				AppendEnvironment(deps.Tool.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func GoModVerify(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("verify").
		ShouldDisable(func(t *Task) bool {
			return !P.Verify
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"mod",
				"verify",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				SetDir(deps.Tool.Cwd).
				Set(func(c *Command) error {
					t.Log.Infof("Verifying modules: in %s", deps.Tool.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
