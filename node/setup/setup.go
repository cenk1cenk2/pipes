// Package setup holds the configuration the node subcommands share. The
// instances live here rather than in main so that the flags that fill them and
// the tasks that read them back are the same ones, whichever subcommand ran.
package setup

import (
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
