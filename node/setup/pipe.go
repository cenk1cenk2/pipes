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

var P = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(a *App) *TaskList[Pipe] {
	return P.New(a).SetTasks(
		SetupPackageManager(&P).Job(),
	)
}
