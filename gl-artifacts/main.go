package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/gl-artifacts/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       pipe.Flags,
				Action: func(c *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.New(p).SetCliContext(c).Job(),
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
