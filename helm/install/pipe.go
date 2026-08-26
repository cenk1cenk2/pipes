package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

// Deps is the resolved helm tool: the chart directory the dependency update runs
// in.
type Deps struct {
	Tool *tool.Ctx
}

var TL = TaskList{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				HelmInstall(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		New: func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
