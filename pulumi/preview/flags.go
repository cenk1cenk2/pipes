package preview

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name: "pulumi.preview.plan",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_PLAN"),
		),
		Usage:       "Output file for pulumi plan.",
		Required:    false,
		Value:       "plan.json",
		Destination: &P.Plan,
	},
}
