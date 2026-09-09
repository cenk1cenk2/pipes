package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

// The flags are built once, since main unhides the environment enable flag on
// this very slice before the command tree reads it.
var Flags = CombineFlags(
	environment.NewFlags(environment.Options{Destination: Environment}),
	node.NewFlags(node.Options{Destination: NodeConfig}),
	node.NewLoginFlags(node.LoginOptions{Destination: Login}),
)
