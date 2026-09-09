package setup

import (
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

// The flags are built once, since main unhides the environment enable flag on the
// slice it hands to the subcommands and a second slice would not carry it.
var (
	EnvironmentFlags = environment.NewFlags(environment.Options{Destination: Environment})
	NodeFlags        = node.NewFlags(node.Options{Destination: NodeConfig})
)
