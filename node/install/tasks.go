package install

import (
	"fmt"
	"os"
	"strings"

	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func InstallNodeDependencies(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("install").
		Set(func(t *Task[Pipe]) error {
			packageManager := setup.TL.Pipe.Ctx.PackageManager

			t.CreateCommand(
				packageManager.Exe,
			).
				Set(func(c *Command[Pipe]) error {
					if TL.Pipe.NodeInstall.UseLockFile {
						c.AppendArgs(packageManager.Commands.InstallWithLock...)

						t.Log.Infoln("Using lockfile for installation.")
					} else {
						c.AppendArgs(packageManager.Commands.Install...)

						t.Log.Infoln("Installing dependencies without a lockfile.")
					}

					c.AppendArgs(strings.Split(t.Pipe.NodeInstall.Args, " ")...)

					if t.Pipe.NodeInstall.Cache {
						cacheDir := fmt.Sprintf(".%s", packageManager.Exe)
						t.Log.Infof("Setting up cache: %s", cacheDir)

						c.AppendArgs(packageManager.Commands.Cache...)
						c.AppendArgs(cacheDir)
					}

					c.SetDir(TL.Pipe.NodeInstall.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendEnvironment(environment.TL.Pipe.Ctx.EnvVars)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
