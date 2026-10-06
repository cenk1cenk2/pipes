package test

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Output struct {
		File string
		Type string `validate:"oneof=JUnit NUnit XUnit Sonar"`
	}

	Pipe struct {
		Files  []string
		Values []string
		Strict bool
		Output
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				test(tl).Job(),
			)
		})
}
