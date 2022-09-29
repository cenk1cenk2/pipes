package setup

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_node_package_manager = "Package Manager"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    category_node_package_manager,
		Name:        "node.package_manager",
		Usage:       "Preferred Package manager for nodejs.",
		Required:    false,
		EnvVars:     []string{"NODE_PACKAGE_MANAGER"},
		Value:       "yarn",
		Destination: &TL.Pipe.Node.PackageManager,
	},
}
