package main

import (
	"github.com/urfave/cli/v2"

	pipe "gitlab.kilic.dev/devops/pipes/_template/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(pipe.Flags),
				Action: func(c *cli.Context) error {
					tl := &pipe.TL

					return tl.RunJobs(
						tl.JobSequence(
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
