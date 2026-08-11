package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
	"gitlab.kilic.dev/devops/pipes/common/report/iac"
)

//revive:disable:line-length-limit

var Flags = CombineFlags(
	[]cli.Flag{
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

		&cli.StringFlag{
			Name: "pulumi.preview.summary-output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_SUMMARY_OUTPUT"),
			),
			Usage:       "Output file for Pulumi preview summary. Leave empty to skip summary generation.",
			Required:    false,
			Value:       "pulumi-summary.json",
			Destination: &P.Summary.Output,
		},
	},
	gitlab.NewMergeRequestReportFlags(&P.MergeRequestReportConfig),
	iac.NewMetadataFlags(&P.ReportMetadata),
)
