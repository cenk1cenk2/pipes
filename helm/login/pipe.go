package login

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type Pipe struct {
	Uri      string
	Username string
	Password string
}

// P is the chart registry the pipe authenticates against.
var P = &Pipe{}

// New is the login stage of every helm command that pulls a dependency or pushes a chart.
func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, _ *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				login(tl).Job(),
			)
		})
}
