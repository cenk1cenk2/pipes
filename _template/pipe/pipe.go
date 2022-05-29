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

func New(p *Plumber) *TaskList[Pipe] {
	return P.New(p).SetTasks(
		P.JobSequence(
			DefaultTask(&P).Job(),
		),
	)
}
