package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Cwd   string `validate:"omitempty,dir"`
		Paths []string
	}

	// Ctx carries the resolved overlays alongside what the setup resolves, so the build reads one context.
	Ctx struct {
		Cwd      string
		Version  string
		Env      map[string]string
		Overlays []string
	}
)

var P = &Pipe{}
var C = &Ctx{Env: map[string]string{}}

func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				initialize(tl).Job(),
				version(tl).Job(),
				ResolveOverlays(tl).Job(),
			)
		})
}
