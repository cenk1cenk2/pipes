package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Cwd       string `validate:"omitempty,dir"`
		Cache     string `validate:"omitempty,dirpath"`
		Workspace bool
	}

	// Ctx carries whether the modules are driven as a workspace, so the commands
	// that tell the two apart never probe the toolchain again.
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
		ShouldRunBefore(func(_ *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				initialize(tl).Job(),
				version(tl).Job(),
				GoEnv(tl).Job(),
				GoWorkspace(tl).Job(),
				GoModules(tl).Job(),
			)
		})
}
