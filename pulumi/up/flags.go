package up

import (
	"github.com/urfave/cli/v3"
)

const CATEGORY_PULUMI_UP = "Up"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_PULUMI_UP,
		Name:     "pulumi.up.plan",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_UP_PLAN"),
			cli.EnvVar("PULUMI_PLAN"),
		),
		Usage:       "Input file for pulumi plan.",
		Required:    false,
		Value:       "plan.json",
		Destination: &P.Plan,
	},
}
