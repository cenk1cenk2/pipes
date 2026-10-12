package setup

import (
	"gitlab.kilic.dev/devops/pipes/internal/node"
)

var NodeFlags = node.NewFlags(node.Options{Destination: NodeConfig})
