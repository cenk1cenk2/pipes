package main

import (
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/node/login"
	node "gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	"gitlab.kilic.dev/devops/pipes/semantic-release/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v5"
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
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(environment.Flags, node.Flags, login.Flags, pipe.Flags),
				Action: func(c *cli.Context) error {
					tl := &pipe.TL

					return tl.RunJobs(
						tl.JobSequence(
							environment.New(p).SetCliContext(c).Job(),
							node.New(p).SetCliContext(c).Job(),
							login.New(p).SetCliContext(c).Job(),
							pipe.New(p).SetCliContext(c).Job(),
						),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags:       true,
			ExcludeHelpCommand: true,
		}).
		Run()
}
