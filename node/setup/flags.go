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
		Category:    CATEGORY_NODE_PACKAGE_MANAGER,
		Name:        "node.package_manager",
		Usage:       `Preferred Package manager for nodejs. enum("npm", "yarn", "pnpm")`,
		Required:    false,
		EnvVars:     []string{"NODE_PACKAGE_MANAGER"},
		Value:       flags.FLAG_DEFAULT_NODE_PACKAGE_MANAGER,
		Destination: &TL.Pipe.Node.PackageManager,
	},
}
