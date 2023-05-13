package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func SetupPackageManager(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.PackageManager = PackageManager{
				Exe:      t.Pipe.Node.PackageManager,
				Commands: PackageManagers[t.Pipe.Node.PackageManager],
			}

			t.Log.Infof("Using package manager: %s", t.Pipe.Node.PackageManager)

			return nil
		})
}

func NodeVersion(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version").
		Set(func(t *Task[Pipe]) error {
			node := t.CreateCommand(
				"node",
				"--version",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				AddSelfToTheTask()

			pm := t.CreateCommand(
				t.Pipe.Ctx.PackageManager.Exe,
			).
				Set(func(c *Command[Pipe]) error {
					c.AppendArgs(t.Pipe.Ctx.PackageManager.Commands.Version...)

					return nil
				}).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				AddSelfToTheTask()

			if err := t.RunCommandJobAsJobParallel(); err != nil {
				return err
			}

			t.Log.Infof("node.js version: %s", node.GetStdoutStream()[0])
			t.Log.Infof("%s version: v%s", t.Pipe.Ctx.PackageManager.Exe, pm.GetStdoutStream()[0])

			return nil
		})
}
