package test

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Runner   string `validate:"omitempty,oneof=ginkgo go"`
		Args     string
		Coverage string
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
				test(tl).Job(),
			)
		})
}
