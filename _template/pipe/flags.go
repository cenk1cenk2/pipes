package pipe

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "default.flag",
		Usage: "Some default flag.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PIPE_DEFAULT_FLAG"),
		),
		Required:    false,
		Value:       "",
		Destination: &P.Default.Flag,
	},
}
