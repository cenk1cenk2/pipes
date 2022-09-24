package main

import (
	"github.com/urfave/cli/v2"

	pipe "gitlab.kilic.dev/devops/pipes/_template/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func main() {
	p := Plumber{
		DocsExcludeFlags: true,
	}

	p.New(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(pipe.Flags),
				Action: func(c *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.TL.JobSequence(
							pipe.New(p).SetCliContext(c).Job(),
						),
					)
				},
			}
		}).Run()
}
