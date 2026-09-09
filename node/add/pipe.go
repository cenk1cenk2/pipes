package add

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Add struct {
		Packages   []string
		Global     bool
		ScriptArgs string
		Cwd        string
	}

	Pipe struct {
		Add
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldDisable(func(tl *TaskList) bool {
			return len(P.Add.Packages) == 0
		}).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				add(tl).Job(),
			)
		})
}
