package stack

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "pulumi.stack",
		Usage: "Stack name for the pulumi to be used in the commands.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_STACK"),
		),
		Required:    true,
		Value:       "",
		Destination: &P.Stack,
	},
}
