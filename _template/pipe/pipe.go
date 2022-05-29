package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	Pipe struct {
		Ctx
	}
)

var P = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(a *App) *TaskList[Pipe] {
	return P.New(a).SetTasks(
		P.JobSequence(
			DefaultTask(&P).Job(),
		),
	)
}
