package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

type (
	Install struct {
		Cwd         string `validate:"dir"`
		UseLockFile bool
		Args        string
		Cache       bool
	}

	Pipe struct {
		Install
	}

	// Deps is the package manager the dependencies are installed with and the
	// environment handed to it. The install command does not select an
	// environment itself, so the variables are only there for the pipes that
	// compose this step after one that does.
	Deps struct {
		Node        *node.Ctx
		Environment *environment.Ctx
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
				InstallNodeDependencies(tl, deps).Job(),
			)
		})
}
