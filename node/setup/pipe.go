package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	Node struct {
		PackageManager string `validate:"oneof=npm yarn"`
	}

	Pipe struct {
		Node
		Ctx
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		SetupPackageManager(&TL).Job(),
	)
}
