package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/node/login"
	"gitlab.kilic.dev/devops/pipes/semantic-release/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v2"
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
				Flags:       p.AppendFlags(pipe.Flags),
				Action: func(ctx *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.TL.JobSequence(
							login.New(p).Job(ctx),
							pipe.New(p).Job(ctx),
						),
					)
				},
			}
		}).Run()
}
