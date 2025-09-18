package lint

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func GoLintWithTool(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		ShouldDisable(func(t *Task) bool {
			return P.Source != "tools"
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"tool",
			).
				SetDir(P.Cwd).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.AppendArgs(P.Tool)

					c.AppendArgs(strings.Split(P.Args, " ")...)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
