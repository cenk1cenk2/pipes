package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/node/login"
	node "gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	"gitlab.kilic.dev/devops/pipes/semantic-release/pipe"
)

func main() {
	OverwriteCliFlag(environment.Flags, func(f *cli.BoolFlag) bool {
		return f.Name == "environment.enable"
	}, func(f *cli.BoolFlag) *cli.BoolFlag {
		f.Hidden = false
		f.Value = false

		return f
	})

	NewPlumber(
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       CombineFlags(environment.Flags, node.Flags, login.Flags, pipe.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(
						JobSequence(
							environment.New(p).Job(),
							node.New(p).Job(),
							login.New(p).Job(),
							pipe.New(p).Job(),
						),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
