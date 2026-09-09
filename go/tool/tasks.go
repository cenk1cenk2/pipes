package tool

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func tool(tl *TaskList) *Task {
	return tl.CreateTask("tool", P.Tool).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"tool",
			).
				SetDir(setup.C.Cwd).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Tool: %s in %s", P.Tool, setup.C.Cwd))

					c.AppendArgs(P.Tool)

					c.AppendArgs(strings.Split(P.Args, " ")...)

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
