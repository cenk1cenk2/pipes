package plan

import (
	"time"

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

	&cli.Uint32Flag{
		Name: "terraform-plan.retry-tries",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_PLAN_RETRY_TRIES"),
		),
		Usage:       "Number of retries for terraform plan command.",
		Required:    false,
		Value:       5,
		Destination: &P.Plan.RetryTries,
	},

	&cli.DurationFlag{
		Name: "terraform-plan.retry-delay",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_PLAN_RETRY_DELAY"),
		),
		Usage:       "Delay between retries for terraform plan command.",
		Required:    false,
		Value:       60 * time.Second,
		Destination: &P.Plan.RetryDelay,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList) error {
	return nil
}
