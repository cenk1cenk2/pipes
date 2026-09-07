// Package setup holds the configuration the node subcommands share. The
// instances live here rather than in main so that the flags that fill them and
// the tasks that read them back are the same ones, whichever subcommand ran.
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
)

// New resolves the package manager every node command runs through.
func New(p *Plumber) *TaskList {
	return node.SetupTaskList(p, NodeConfig, NodeCtx)
}

// NewEnvironment selects the environment the scripts are templated against and
// run with.
func NewEnvironment(p *Plumber) *TaskList {
	return environment.SetupTaskList(p, Environment, EnvironmentCtx)
}
