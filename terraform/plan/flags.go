package plan

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name: "terraform-plan.out",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_PLAN_CACHE"),
			cli.EnvVar("TF_APPLY_OUTPUT"),
			cli.EnvVar("TF_PLAN_OUTPUT"),
		),
		Usage:       "Output file for terraform plan.",
		Required:    false,
		Value:       "plan",
		Destination: &P.Plan.Output,
	},

	&cli.StringFlag{
		Name: "terraform-plan.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_PLAN_ARGS"),
		),
		Usage:       "Additional arguments for terraform plan.",
		Required:    false,
		Value:       "",
		Destination: &P.Plan.Args,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList) error {
	return nil
}
