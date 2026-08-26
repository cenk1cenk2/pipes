package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_PACKAGE_MANAGER = "Package Manager"

	DEFAULT_PACKAGE_MANAGER = "pnpm"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_PACKAGE_MANAGER

	&cli.StringFlag{
		Category: CATEGORY_NODE_PACKAGE_MANAGER,
		Name:     "node.package_manager",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_PACKAGE_MANAGER"),
		),
		Usage:       `Preferred Package manager for nodejs. enum("npm", "yarn", "pnpm")`,
		Required:    false,
		Value:       DEFAULT_PACKAGE_MANAGER,
		Destination: &P.Node.PackageManager,
	},
}
