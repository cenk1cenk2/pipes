package lint

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Args    string
		Timeout time.Duration
		Cache   string
	}

	Ctx struct {
		Modules []string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GoLint(tl).Job(),
			)
		})
}
