package test

import (
	"context"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/python/setup"
)

// The environment the install command synced is what the commands run in, so uv is
// told not to sync it again.
func test(tl *TaskList) *Task {
	return tl.CreateTask("test").
		Set(func(_ context.Context, t *Task) error {
			for _, command := range P.Commands {
				t.CreateCommand(
					"uv",
					"run",
					"--no-sync",
				).
					SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDebug).
					Set(func(_ context.Context, c *Command) error {
						c.AppendArgs(strings.Fields(command)...)

						return nil
					}).
					SetDir(setup.C.Cwd).
					AppendEnvironment(setup.C.Env).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
