package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	Node struct {
		PackageManager string `validate:"oneof=npm yarn"`
	}

	Pipe struct {
		Node Node
	}
)

var P = TaskList[Pipe, Ctx]{
	Pipe:    Pipe{},
	Context: Ctx{},
}

func New(a *App) *TaskList[Pipe, Ctx] {
	return P.New(a).SetTasks(
		SetupPackageManager(&P).Job(),
	)
}
