package run

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/python/setup"
)

func run(tl *TaskList) *Task {
	return tl.CreateTask("run").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"uv",
				"run",
				"--no-sync",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				Set(func(_ context.Context, c *Command) error {
					c.AppendArgs(P.Command...)

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
