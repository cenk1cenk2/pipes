package plan

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "terraform-plan.out",
		Usage:       "Output file for terraform plan.",
		Required:    false,
		EnvVars:     []string{"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT"},
		Value:       "plan",
		Destination: &TL.Pipe.Plan.Output,
	},

	&cli.StringFlag{
		Name:        "terraform-plan.args",
		Usage:       "Additional arguments for terraform plan.",
		Required:    false,
		EnvVars:     []string{"TF_PLAN_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.Plan.Args,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
