package main

import (
	"github.com/urfave/cli/v2"

	pipe "gitlab.kilic.dev/devops/pipes/_template/pipe"
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
					return pipe.P.RunJobs(
						pipe.P.JobSequence(
							pipe.New(p).Job(ctx),
						),
					)
				},
			}
		}).Run()
}
