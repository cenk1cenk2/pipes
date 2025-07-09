package login

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	DockerRegistry struct {
		Registry string
		Username string
		Password string
	}

	Pipe struct {
		DockerRegistry
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList[Pipe]) Job {
			return JobSequence(
				DockerLogin(tl).Job(),
			)
		})
}
