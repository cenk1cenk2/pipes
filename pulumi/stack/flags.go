package stack

import (
	"github.com/urfave/cli/v3"
)

const CategoryPulumiStack = "Stack"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CategoryPulumiStack,
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
