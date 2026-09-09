// Package setup holds the configuration the semantic-release pipe shares between
// its stages, so that the flags that fill it and the tasks that read it back are
// the same instances.
package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

var (
	Environment    = &environment.Config{}
	EnvironmentCtx = &environment.Ctx{}

	NodeConfig = &node.Config{}
	NodeCtx    = &node.Ctx{}

	Login = &node.Login{}
)

// New composes every setup stage into the one task list the command wires in:
// the injected environment, the package manager and the registry credentials.
// The inner lists keep their own validation and disable hooks, which is why
// they are combined as a job rather than flattened.
func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		Set(func(_ *TaskList) Job {
			return CombineTaskLists(
				environment.SetupTaskList(p, Environment, EnvironmentCtx),
				node.SetupTaskList(p, NodeConfig, NodeCtx),
				node.LoginTaskList(p, Login),
			)
		})
}
