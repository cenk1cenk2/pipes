package apply

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "terraform-apply.out",
		Usage:       "Output file for terraform apply.",
		Required:    false,
		EnvVars:     []string{"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT"},
		Value:       "plan",
		Destination: &TL.Pipe.Apply.Output,
	},

	&cli.StringFlag{
		Name:        "terraform-apply.args",
		Usage:       "Additional arguments for terraform apply.",
		Required:    false,
		EnvVars:     []string{"TF_APPLY_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.Apply.Args,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
