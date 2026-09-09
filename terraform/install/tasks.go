package install

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func install(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"terraform",
				"init",
				"-input=false",
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Install.Reconfigure {
						t.Log.Info("Will reconfigure state.")

						c.AppendArgs("-reconfigure")
					}

					if P.Install.UseLockfile {
						t.Log.Info("Using lockfile.")

						c.AppendArgs("-lockfile=readonly")
					}

					if P.Install.Args != "" {
						c.AppendArgs(P.Install.Args)
					}

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
