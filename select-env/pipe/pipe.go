package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Environment struct {
		Conditions string
	}

	Pipe struct {
		Ctx

		Environment
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
