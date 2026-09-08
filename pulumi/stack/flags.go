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
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_STACK"),
		),
		Usage:       "Stack name to use for pulumi commands.",
		Required:    true,
		Value:       "",
		Destination: &P.Stack,
	},
}
