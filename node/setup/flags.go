package setup

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_PACKAGE_MANAGER = "Package Manager"
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
		Value:       flags.FLAG_DEFAULT_NODE_PACKAGE_MANAGER,
		Destination: &P.Node.PackageManager,
	},
}
