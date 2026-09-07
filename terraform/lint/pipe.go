package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Lint struct {
		FormatCheckEnable bool
		FormatCheckArgs   string
		ValidateEnable    bool
		ValidateArgs      string
	}

	Pipe struct {
		Lint
	}

	// Deps is the resolved terraform tool: the directory the lint runs in and the
	// environment the setup step has written into.
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
				TerraformLint(tl, deps).Job(),
			)
		})
}
