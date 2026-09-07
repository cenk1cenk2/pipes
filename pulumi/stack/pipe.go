package stack

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Pipe struct {
		Stack string
	}

	// Deps is the resolved pulumi tool: the directory the stack is selected in,
	// which every command that follows then runs against.
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
				PulumiSelectStack(tl, deps).Job(),
			)
		})
}
