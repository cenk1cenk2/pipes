// Package setup holds the configuration the node subcommands share. The
// instances live here rather than in main so that the flags that fill them and
// the tasks that read them back are the same ones, whichever subcommand ran.
package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
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

// The steps are values rather than factories for the same reason: a command
// composes the one above instead of registering a second copy of it.
var (
	// Step resolves the package manager every node command runs through.
	Step = cli.Step{Flags: NodeFlags, New: func(p *plumber.Plumber) *plumber.TaskList {
		return node.SetupTaskList(p, NodeConfig, NodeCtx)
	}}

	// LoginStep writes the npmrc the package manager reads its credentials back
	// from, so the commands that reach a registry compose it after Step.
	LoginStep = cli.Step{Flags: LoginFlags, New: func(p *plumber.Plumber) *plumber.TaskList {
		return node.LoginTaskList(p, Login)
	}}

	// EnvironmentStep selects the environment the scripts are templated against
	// and run with.
	EnvironmentStep = cli.Step{Flags: EnvironmentFlags, New: func(p *plumber.Plumber) *plumber.TaskList {
		return environment.TaskList(p, Environment, EnvironmentCtx)
	}}
)
