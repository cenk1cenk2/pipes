package add

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func add(tl *TaskList) *Task {
	return tl.CreateTask("add").
		Set(func(_ context.Context, t *Task) error {
			packageManager := setup.NodeCtx.PackageManager

			t.CreateCommand(
				packageManager.Exe,
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Add.Global {
						c.AppendArgs(packageManager.Commands.Global...)
					}

					c.AppendArgs(packageManager.Commands.Add...)

					if P.Add.ScriptArgs != "" {
						c.AppendArgs(P.Add.ScriptArgs)
					}

					c.AppendArgs(P.Add.Packages...)

					c.SetDir(P.Add.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
