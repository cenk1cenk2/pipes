// Package setup holds the configuration the semantic-release pipe shares
// between its task lists. The instances live here rather than in main so that
// the flags that fill them and the tasks that read them back are the same ones.
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

// The flags are built once, since main unhides the environment enable flag on
// the slice it hands to the command and a second slice would not carry it.
var (
	EnvironmentFlags = environment.NewFlags(environment.Options{Destination: Environment})
	NodeFlags        = node.NewFlags(node.Options{Destination: NodeConfig})
	LoginFlags       = node.NewLoginFlags(node.LoginOptions{Destination: Login})
)

// NewEnvironment selects the environment the release runs against.
func NewEnvironment(p *Plumber) *TaskList {
	return environment.SetupTaskList(p, Environment, EnvironmentCtx)
}

// New resolves the package manager the release library is installed with.
func New(p *Plumber) *TaskList {
	return node.SetupTaskList(p, NodeConfig, NodeCtx)
}

// NewLogin writes the npmrc the package manager reads its credentials back
// from, so it composes after New.
func NewLogin(p *Plumber) *TaskList {
	return node.LoginTaskList(p, Login)
}
