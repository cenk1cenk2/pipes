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

// P is the registry the pipe authenticates against; the build and the manifest
// commands read its uri back to prefix the images they publish.
var P = &Pipe{}

// New is the login stage of every buildah command.
func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, _ *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				loginParent(tl).Job(),
			)
		})
}
