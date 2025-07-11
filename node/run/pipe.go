package run

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	NodeCommand struct {
		Script string
		Cwd    string `validate:"dir"`
	}

	Pipe struct {
		NodeCommand
	}

	Ctx struct {
		Script     string
		ScriptArgs string
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
				RunNodeScript(tl).Job(),
			)
		})
}
