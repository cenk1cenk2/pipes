package install

import (
	"os"
	"strings"

	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
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

					if t.Pipe.NodeCache.Enable {
						t.Log.Infof("Setting up cache for %s.", packageManager.Exe)

						switch packageManager.Exe {
						case "npm":
							c.AppendArgs("--cache", t.Pipe.NodeCache.NpmCacheDir, "--prefer-offline")
						case "yarn":
							c.AppendArgs("--cache-folder", t.Pipe.NodeCache.YarnCacheDir, "--prefer-offline")
						case "pnpm":
							c.AppendArgs("--store-dir", t.Pipe.NodeCache.PnpmCacheDir, "--prefer-offline")
						}
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
