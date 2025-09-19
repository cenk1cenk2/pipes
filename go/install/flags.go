package install

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_INSTALL = "Install"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Category: CATEGORY_INSTALL,
		Name:     "go.install.verify",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_INSTALL_VERIFY"),
		),
		Usage:       "Use the sum file to verify module integrity.",
		Required:    false,
		Value:       true,
		Destination: &P.Verify,
	},

	&cli.StringFlag{
		Category: CATEGORY_INSTALL,
		Name:     "go.install.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_INSTALL_ARGS"),
		),
		Usage:       "Arguments to append to install command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},
}
