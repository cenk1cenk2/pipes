package tool

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func tool(tl *TaskList) *Task {
	return tl.CreateTask("tool", P.Tool).
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"tool",
			).
				SetDir(setup.C.Cwd).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					t.Log.Infof("Tool: %s in %s", P.Tool, setup.C.Cwd)

					c.AppendArgs(P.Tool)

					c.AppendArgs(strings.Split(P.Args, " ")...)

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
