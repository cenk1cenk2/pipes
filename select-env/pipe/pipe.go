package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Environment struct {
		File string
	}

	Pipe struct {
		Environment
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).Set(func(tl *TaskList[Pipe]) Job {
		return tl.JobSequence(
			WriteEnvironmentFile(tl).Job(),
		)
	})
}
