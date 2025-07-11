package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func SetupPackageManager(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(t *Task) error {
			C.PackageManager = PackageManager{
				Exe:      P.Node.PackageManager,
				Commands: PackageManagers[P.Node.PackageManager],
			}

			t.Log.Infof("Using package manager: %s", P.Node.PackageManager)

			return nil
		})
}

func NodeVersion(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand(
				"node",
				"--version",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					stream := c.GetCombinedStream()

					if len(stream) == 0 {
						t.Log.Debugln("Can not fetch node.js version.")

						return nil
					}

					t.Log.Infof("node.js version: %s", stream[0])

					return nil
				}).
				AddSelfToTheTask()

			t.CreateCommand(
				C.PackageManager.Exe,
			).
				Set(func(c *Command) error {
					c.AppendArgs(C.PackageManager.Commands.Version...)

					return nil
				}).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					stream := c.GetCombinedStream()

					if len(stream) == 0 {
						t.Log.Debugln("Can not fetch package manager version.")

						return nil
					}

					t.Log.Infof("%s version: v%s", C.PackageManager.Exe, stream[0])

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
