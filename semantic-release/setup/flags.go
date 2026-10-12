package setup

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

var Flags = CombineFlags(
	node.NewFlags(node.Options{Destination: NodeConfig}),
	node.NewLoginFlags(node.LoginOptions{Destination: Login}),
)
