package plan

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

//revive:disable:line-length-limit

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

		&cli.BoolFlag{
			Name: "terraform-plan.preview-for-mrs",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TF_PLAN_PREVIEW_FOR_MRS"),
			),
			Usage:       "Run merge request terraform plans as previews without state locking.",
			Required:    false,
			Value:       true,
			Destination: &P.Plan.PreviewForMergeRequests,
		},

		&cli.StringFlag{
			Name: "terraform-plan.pipeline-source",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PIPELINE_SOURCE"),
			),
			Usage:       "GitLab CI pipeline source used to detect merge request pipelines.",
			Required:    false,
			Value:       "",
			Destination: &P.Plan.PipelineSource,
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

		&cli.StringFlag{
			Name: "terraform-plan.summary-output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_SUMMARY_OUTPUT"),
			),
			Usage:       "Output file for terraform plan summary. Leave empty to skip summary generation.",
			Required:    false,
			Value:       "terraform-summary.json",
			Destination: &P.Summary.Output,
		},
	},
	gitlab.NewMergeRequestReportFlags(&P.MergeRequestReport),
	ci.NewMetadataFlags(&P.ReportMetadata),
)

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList) error {
	return nil
}
