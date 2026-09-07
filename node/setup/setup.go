// Package setup holds the configuration the node subcommands share. The
// instances live here rather than in main so that the flags that fill them and
// the tasks that read them back are the same ones, whichever subcommand ran.
package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
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

// The flags are built once, since main unhides the environment enable flag on
// the slice it hands to the subcommands and a second slice would not carry it.
var (
	EnvironmentFlags = environment.NewFlags(Environment)
	NodeFlags        = node.NewFlags(NodeConfig)
	LoginFlags       = node.NewLoginFlags(Login)
)

// New resolves the package manager every node command runs through.
func New(p *plumber.Plumber) *plumber.TaskList {
	return node.SetupTaskList(p, NodeConfig, NodeCtx)
}

// NewLogin writes the npmrc the package manager reads its credentials back
// from, so the commands that reach a registry compose it after New.
func NewLogin(p *plumber.Plumber) *plumber.TaskList {
	return node.LoginTaskList(p, Login)
}

// NewEnvironment selects the environment the scripts are templated against and
// run with.
func NewEnvironment(p *plumber.Plumber) *plumber.TaskList {
	return environment.TaskList(p, Environment, EnvironmentCtx)
}
