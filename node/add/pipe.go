package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

type (
	Add struct {
		Packages   []string
		Global     bool
		ScriptArgs string
		Cwd        string
	}

	Pipe struct {
		Add
	}

	// Deps is the package manager the packages are added with and the
	// environment the script arguments are templated against.
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
		ShouldDisable(func(tl *TaskList) bool {
			return len(P.Add.Packages) == 0
		}).
		ShouldRunBefore(func(tl *TaskList) error {
			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				AddNodeModules(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
