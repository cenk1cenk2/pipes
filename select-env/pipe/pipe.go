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
	return TL.New(p).Set(func(tl *TaskList[Pipe]) Job {
		return tl.JobSequence(
			DefaultTask(tl).Job(),
		)
	})
}
