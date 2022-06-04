package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type (
	Default struct {
		Flag string
	}

	Pipe struct {
		Ctx
		Default
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			DefaultTask(&TL).Job(),
		),
	)
}
