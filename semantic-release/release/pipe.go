package release

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
)

type (
	CI struct {
		CommitReference string
	}

	SemanticRelease struct {
		IsDryRun          bool
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

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				RunSemanticRelease(tl).Job(),
			)
		})
}

// The release binary is run directly rather than through the package manager,
// so the step reads nothing the steps before it resolved.
var Step = icli.Step{Flags: Flags, New: New}
