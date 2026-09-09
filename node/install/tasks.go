package install

import (
	"context"
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func install(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(_ context.Context, t *Task) error {
			packageManager := setup.NodeCtx.PackageManager

			t.CreateCommand(
				packageManager.Exe,
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Install.UseLockFile {
						c.AppendArgs(packageManager.Commands.InstallWithLock...)

						t.Log.Info("Using lockfile for installation.")
					} else {
						c.AppendArgs(packageManager.Commands.Install...)

						t.Log.Info("Installing dependencies without a lockfile.")
					}

					c.AppendArgs(strings.Split(P.Install.Args, " ")...)

					if P.Install.Cache {
						cacheDir := fmt.Sprintf(".%s", packageManager.Exe)
						t.Log.Info(fmt.Sprintf("Setting up cache: %s", cacheDir))

						c.AppendArgs(packageManager.Commands.Cache...)
						c.AppendArgs(cacheDir)
					}

					c.SetDir(P.Install.Cwd)

					c.AppendDirectEnvironment(os.Environ()...)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
