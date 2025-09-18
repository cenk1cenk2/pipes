package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Cache string `validate:"omitempty,dirpath"`
	}

	Ctx struct {
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GoVersion(tl).Job(),
			)
		})
}
