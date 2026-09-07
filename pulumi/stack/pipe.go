package stack

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Stack string
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
				PulumiSelectStack(tl).Job(),
			)
		})
}
