package pipe

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "node.package_manager",
		Usage:       "Preferred Package manager for nodejs.",
		Required:    false,
		EnvVars:     []string{"NODE_PACKAGE_MANAGER"},
		Value:       "npm",
		Destination: &Pipe.Node.PackageManager,
	},
}
