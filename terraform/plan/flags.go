package plan

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const reportCategory = "Report"

var Flags = CombineFlags(
	[]cli.Flag{
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

		&cli.BoolFlag{
			Category: reportCategory,
			Name:     "report.enabled",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_REPORT_ENABLED"),
			),
			Usage:       "Generate Terraform report artifacts.",
			Required:    false,
			Value:       true,
			Destination: &P.Report.Enabled,
		},

		&cli.StringFlag{
			Category: reportCategory,
			Name:     "report.output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_REPORT_OUTPUT"),
			),
			Usage:       "Output file for Terraform report artifact.",
			Required:    false,
			Value:       "terraform-report.html",
			Destination: &P.Report.Output,
		},
	},
)

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList) error {
	return nil
}
