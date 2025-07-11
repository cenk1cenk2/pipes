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

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				WriteEnvironmentFile(tl).Job(),
			)
		})
}
