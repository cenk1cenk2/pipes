package login

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
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
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				DockerLogin(tl).Job(),
			)
		})
}
