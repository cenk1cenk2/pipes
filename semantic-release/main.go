package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/node/login"
	"gitlab.kilic.dev/devops/pipes/semantic-release/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func main() {
	p := Plumber{}

	p.New(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(login.Flags, pipe.Flags),
				Action: func(c *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.TL.JobSequence(
							login.New(p).SetCliContext(c).Job(),
							pipe.New(p).SetCliContext(c).Job(),
						),
					)
				},
			}
		}).Run()
}
