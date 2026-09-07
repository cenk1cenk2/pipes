package apply

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Apply struct {
		Args   string
		Output string
	}

	Pipe struct {
		Apply
	}

	// Deps is the resolved terraform tool: the directory the apply runs in and the
	// environment the registry and state steps have written into.
	Deps struct {
		Tool *tool.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				TerraformApply(tl, deps).Job(),
			)
		})
}
