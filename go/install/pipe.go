package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Pipe struct {
		Verify bool
		Args   string
	}

	Ctx struct {
		Workspace string
	}

	// Deps is the resolved go tool: the directory the modules are vendored in and
	// the environment the cache setup has written into.
	Deps struct {
		Tool *tool.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GoModVendor(tl, deps).Job(),
				GoModVerify(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
