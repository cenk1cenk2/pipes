package apply

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name: "terraform-apply.out",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_PLAN_CACHE"),
			cli.EnvVar("TF_APPLY_OUTPUT"),
			cli.EnvVar("TF_PLAN_OUTPUT"),
		),
		Usage:       "Output file for terraform apply.",
		Required:    false,
		Value:       "plan",
		Destination: &P.Apply.Output,
	},

	&cli.StringFlag{
		Name: "terraform-apply.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_APPLY_ARGS"),
		),
		Usage:       "Additional arguments for terraform apply.",
		Required:    false,
		Value:       "",
		Destination: &P.Apply.Args,
	},
}
