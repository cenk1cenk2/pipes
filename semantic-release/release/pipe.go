package release

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	CI struct {
		CommitReference string
	}

	SemanticRelease struct {
		DryRun            bool
		Workspace         bool
		IsolateWorkspaces []string
	}

	Pipe struct {
		SemanticRelease
		CI
	}

	Ctx struct {
		Exe string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

// The release binary is run directly, so this stage reads nothing the stages before it resolved.
func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				release(tl).Job(),
			)
		})
}
