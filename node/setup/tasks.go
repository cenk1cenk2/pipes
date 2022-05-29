package setup

import (
	"github.com/workanator/go-floc/v3"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
	PackageManager
}

func SetupPackageManager(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("setup").
		Set(func(t *Task[Pipe], c floc.Control) error {
			t.Pipe.Ctx.PackageManager = PackageManager{
				Exe:      t.Pipe.Node.PackageManager,
				Commands: PackageManagers[t.Pipe.Node.PackageManager],
			}

			t.Log.Debugf("Using package manager: %s", t.Pipe.Node.PackageManager)

			return nil
		})
}
