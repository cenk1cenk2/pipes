package setup

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Pipe struct {
		Cwd     string `validate:"omitempty,dir"`
		DataDir string `validate:"required"`
	}

	// Ctx keeps the two binaries apart: Executable is the mise the image ships,
	// Binary is its copy inside the data directory that the shims link to.
	Ctx struct {
		Cwd        string
		Version    string
		Executable string
		DataDir    string
		Binary     string
		Env        map[string]string
	}
)

var P = &Pipe{}
var C = &Ctx{Env: map[string]string{}}

func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, _ *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				initialize(tl).Job(),
				version(tl).Job(),
				environment(tl).Job(),
			)
		})
}
