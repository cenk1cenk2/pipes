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

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validator.Struct(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				DockerLogin(tl).Job(),
			)
		})
}
