package login

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				DockerLoginParent(tl).Job(),
			)
		})
}
