package apply

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

const CategoryApply = "Apply"

//revive:disable:line-length-limit

var Flags = CombineFlags(
	[]cli.Flag{
		&cli.StringFlag{
			Category: CategoryApply,
			Name:     "terraform.apply.output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_APPLY_OUTPUT"),
				cli.EnvVar("TF_PLAN_CACHE"),
				cli.EnvVar("TF_APPLY_OUTPUT"),
				cli.EnvVar("TF_PLAN_OUTPUT"),
			),
			Usage:       "Output file for terraform apply.",
			Required:    false,
			Value:       "plan",
			Destination: &P.Apply.Output,
		},

		&cli.StringFlag{
			Category: CategoryApply,
			Name:     "terraform.apply.args",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_APPLY_ARGS"),
				cli.EnvVar("TF_APPLY_ARGS"),
			),
			Usage:       "Additional arguments for terraform apply.",
			Required:    false,
			Value:       "",
			Destination: &P.Apply.Args,
		},
	},
	gitlab.NewFlags(gitlab.Options{Destination: &P.MergeRequestReport}),
	ci.NewFlags(ci.Options{Destination: &P.ReportMetadata}),
)
