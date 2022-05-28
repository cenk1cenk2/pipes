package setup

import (
	"github.com/workanator/go-floc/v3"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
	PackageManager
}

func SetupPackageManager(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "setup").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		t.Context.PackageManager = PackageManager{
			Exe:      t.Pipe.Node.PackageManager,
			Commands: PackageManagers[t.Pipe.Node.PackageManager],
		}

		t.Log.Debugf("Using package manager: %s", t.Pipe.Node.PackageManager)

		return nil
	})
}
