package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
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
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				WriteEnvironmentFile(tl).Job(),
			)
		})
}
