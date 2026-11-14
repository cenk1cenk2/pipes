package up

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

const CategoryPulumiUp = "Up"

//revive:disable:line-length-limit

var Flags = CombineFlags(
	[]cli.Flag{
		&cli.StringFlag{
			Category: CategoryPulumiUp,
			Name:     "pulumi.up.plan",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_UP_PLAN"),
				cli.EnvVar("PULUMI_PLAN"),
			),
			Usage:       "Input file for pulumi plan.",
			Required:    false,
			Value:       "plan.json",
			Destination: &P.Plan,
		},
	},
	gitlab.NewFlags(gitlab.Options{Destination: &P.MergeRequestReport}),
	ci.NewFlags(ci.Options{Destination: &P.ReportMetadata}),
)
