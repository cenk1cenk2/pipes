package lint

import (
	"context"
	"time"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Pipe struct {
		Args    string
		Timeout time.Duration
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
				lint(tl).Job(),
			)
		})
}
