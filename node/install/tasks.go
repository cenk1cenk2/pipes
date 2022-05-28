package install

import (
	"github.com/workanator/go-floc/v3"
	pipe "gitlab.kilic.dev/devops/pipes/node/setup"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
}

func InstallNodeDependencies(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "install").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		cmd := Command[Pipe, Ctx]{}

		cmd.New(t, pipe.P.Context.PackageManager.Exe).Set(func(c *Command[Pipe, Ctx]) error {
			if P.Pipe.NodeInstall.UseLockFile {
				c.AppendArgs(pipe.P.Context.PackageManager.Commands.InstallWithLock...)

				t.Log.Debugln("Using lockfile for installation.")
			} else {
				c.AppendArgs(pipe.P.Context.PackageManager.Commands.Install...)

				t.Log.Debugln("Installing dependencies without a lockfile.")
			}

			c.SetDir(P.Pipe.NodeInstall.Cwd)

			return nil
		})

		t.AddCommands(cmd)

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.RunCommandJobAsJobSequence()
	})
}
