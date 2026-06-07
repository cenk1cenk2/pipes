package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const reportCategory = "Report"

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

		&cli.BoolFlag{
			Category: reportCategory,
			Name:     "report.enabled",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_REPORT_ENABLED"),
			),
			Usage:       "Generate Pulumi report artifact.",
			Required:    false,
			Value:       true,
			Destination: &P.Report.Enabled,
		},

		&cli.StringFlag{
			Category: reportCategory,
			Name:     "report.output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_REPORT_OUTPUT"),
			),
			Usage:       "Output file for Pulumi report artifact.",
			Required:    false,
			Value:       "pulumi-report.html",
			Destination: &P.Report.Output,
		},
	},
)
