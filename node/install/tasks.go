package install

import (
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func InstallNodeDependencies(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("install").
		Set(func(t *Task) error {
			packageManager := deps.Node.PackageManager

			t.CreateCommand(
				packageManager.Exe,
			).
				Set(func(c *Command) error {
					if P.Install.UseLockFile {
						c.AppendArgs(packageManager.Commands.InstallWithLock...)

						t.Log.Infoln("Using lockfile for installation.")
					} else {
						c.AppendArgs(packageManager.Commands.Install...)

						t.Log.Infoln("Installing dependencies without a lockfile.")
					}

					c.AppendArgs(strings.Split(P.Install.Args, " ")...)

					if P.Install.Cache {
						cacheDir := fmt.Sprintf(".%s", packageManager.Exe)
						t.Log.Infof("Setting up cache: %s", cacheDir)

						c.AppendArgs(packageManager.Commands.Cache...)
						c.AppendArgs(cacheDir)
					}

					c.SetDir(P.Install.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendEnvironment(deps.Environment.EnvVars)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
