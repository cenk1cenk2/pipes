package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type Ctx struct {
	PackageManager
}

func SetupPackageManager(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("setup").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.PackageManager = PackageManager{
				Exe:      t.Pipe.Node.PackageManager,
				Commands: PackageManagers[t.Pipe.Node.PackageManager],
			}

			t.Log.Debugf("Using package manager: %s", t.Pipe.Node.PackageManager)

			return nil
		})
}
