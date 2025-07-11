package install

import (
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func InstallNodeDependencies(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(t *Task) error {
			packageManager := setup.C.PackageManager

			t.CreateCommand(
				packageManager.Exe,
			).
				Set(func(c *Command) error {
					if P.NodeInstall.UseLockFile {
						c.AppendArgs(packageManager.Commands.InstallWithLock...)

						t.Log.Infoln("Using lockfile for installation.")
					} else {
						c.AppendArgs(packageManager.Commands.Install...)

						t.Log.Infoln("Installing dependencies without a lockfile.")
					}

					c.AppendArgs(strings.Split(P.NodeInstall.Args, " ")...)

					if P.NodeInstall.Cache {
						cacheDir := fmt.Sprintf(".%s", packageManager.Exe)
						t.Log.Infof("Setting up cache: %s", cacheDir)

						c.AppendArgs(packageManager.Commands.Cache...)
						c.AppendArgs(cacheDir)
					}

					c.SetDir(P.NodeInstall.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendEnvironment(environment.C.EnvVars)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
