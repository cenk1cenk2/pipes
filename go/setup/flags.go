package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_SETUP = "Setup"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_SETUP,
		Name:     "go.cache",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_CACHE"),
		),
		Usage:       "Enable go cache.",
		Required:    false,
		Value:       "",
		Destination: &P.Cache,
	},
}
