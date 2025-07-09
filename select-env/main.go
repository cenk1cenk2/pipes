package main

import (
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/select-env/pipe"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "github.com/cenk1cenk2/plumber/v6"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(setup.Flags, pipe.Flags),
				Action: func(c *cli.Context) error {
					tl := &pipe.TL

					return tl.RunJobs(
						tl.JobSequence(
							setup.New(p).SetCli(c).Job(),
							pipe.New(p).SetCli(c).Job(),
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
