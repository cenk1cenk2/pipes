package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Lint struct {
		FormatCheckEnable bool
		FormatCheckArgs   string
		ValidateEnable    bool
		ValidateArgs      string
	}

	Pipe struct {
		Lint
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				TerraformLint(tl).Job(),
			)
		})
}
