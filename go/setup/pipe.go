package setup

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Pipe struct {
		Cwd   string `validate:"omitempty,dir"`
		Cache string `validate:"omitempty,dirpath"`
	}

	// Ctx carries whether the modules are driven as a workspace, so the commands
	// that tell the two apart never look for the workspace file again.
	Ctx struct {
		Cwd       string
		Version   string
		Env       map[string]string
		Workspace bool
		// Modules are the directories the workspace drives, resolved once for lint and build.
		Modules []string
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
				env(tl).Job(),
				workspace(tl).Job(),
				modules(tl).Job(),
			)
		})
}
