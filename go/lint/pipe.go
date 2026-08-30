package lint

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
)

type (
	Pipe struct {
		Args    string
		Timeout time.Duration
		Cache   string
	}

	Ctx struct {
		Modules []string
	}

	// Deps is the resolved go tool: the directory the modules are listed from,
	// whether they are a workspace and the environment the cache setup has
	// written into.
	Deps struct {
		Tool *setup.Ctx
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
				GoLint(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
