package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

const CategoryPulumiPreview = "Preview"

//revive:disable:line-length-limit

var Flags = CombineFlags(
	[]cli.Flag{
		&cli.StringFlag{
			Category: CategoryPulumiPreview,
			Name:     "pulumi.preview.plan",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_PREVIEW_PLAN"),
				cli.EnvVar("PULUMI_PLAN"),
			),
			Usage:       "Output file for pulumi plan.",
			Required:    false,
			Value:       "plan.json",
			Destination: &P.Plan,
		},

		&cli.StringFlag{
			Category: CategoryPulumiPreview,
			Name:     "pulumi.preview.summary.output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_PREVIEW_SUMMARY_OUTPUT"),
				cli.EnvVar("PULUMI_SUMMARY_OUTPUT"),
			),
			Usage:       "Output file for pulumi preview summary. Leave empty to skip summary generation.",
			Required:    false,
			Value:       "pulumi-summary.json",
			Destination: &P.Summary.Output,
		},
	},
	gitlab.NewFlags(gitlab.Options{Destination: &P.MergeRequestReport}),
	ci.NewFlags(ci.Options{Destination: &P.ReportMetadata}),
)
