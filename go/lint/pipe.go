package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Cwd    string `validate:"dirpath"`
		Args   string
		Source string `validate:"oneof=tools"`
		Tool   string
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
				GoLintWithTool(tl).Job(),
			)
		})
}
