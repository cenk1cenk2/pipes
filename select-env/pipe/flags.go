package pipe

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "default.flag",
		Usage:       "Some default flag.",
		Required:    false,
		EnvVars:     []string{"DEFAULT_FLAG"},
		Value:       "",
		Destination: &TL.Pipe.Default.Flag,
	},
}
