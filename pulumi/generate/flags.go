package generate

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

const CategoryPulumiGenerate = "Generate"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CategoryPulumiGenerate,
		Name:     "pulumi.generate.build",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_GENERATE_BUILD"),
		),
		Usage:       "Shell command that builds the Pulumi program before the generation.",
		Required:    false,
		Value:       "",
		Destination: &P.Build,
	},

	flags.YAMLFlag(&P.Paths, &cli.StringFlag{
		Category: CategoryPulumiGenerate,
		Name:     "pulumi.generate.paths",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_GENERATE_PATHS"),
		),
		Usage:    "Paths the generation writes, relative to the working directory, removed before every run so the files it no longer generates are gone. format(yaml([]string))",
		Required: false,
		Value:    `[]`,
	}),

	&cli.StringFlag{
		Category: CategoryPulumiGenerate,
		Name:     "pulumi.generate.command",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_GENERATE_COMMAND"),
		),
		Usage:       "Shell command that generates the sources against the selected stack. Defaults to deploying the stack.",
		Required:    false,
		Value:       "",
		Destination: &P.Command,
	},
}
