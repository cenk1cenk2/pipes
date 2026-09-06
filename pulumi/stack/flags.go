package stack

import (
	"github.com/urfave/cli/v3"
)

const CATEGORY_PULUMI_STACK = "Stack"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_PULUMI_STACK,
		Name:     "pulumi.stack",
		Usage:    "Stack name for the pulumi to be used in the commands.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_STACK"),
		),
		Required:    true,
		Value:       "",
		Destination: &P.Stack,
	},
}
