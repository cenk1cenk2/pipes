package install

import (
	"strings"

	pipe "gitlab.kilic.dev/devops/pipes/node/setup"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func InstallNodeDependencies(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("install").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(pipe.TL.Pipe.Ctx.PackageManager.Exe).Set(func(c *Command[Pipe]) error {
				if TL.Pipe.NodeInstall.UseLockFile {
					c.AppendArgs(pipe.TL.Pipe.Ctx.PackageManager.Commands.InstallWithLock...)

					t.Log.Debugln("Using lockfile for installation.")
				} else {
					c.AppendArgs(pipe.TL.Pipe.Ctx.PackageManager.Commands.Install...)

					t.Log.Debugln("Installing dependencies without a lockfile.")
				}

				c.AppendArgs(strings.Split(t.Pipe.NodeInstall.Args, " ")...)

				c.SetDir(TL.Pipe.NodeInstall.Cwd)

				return nil
			}).AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
