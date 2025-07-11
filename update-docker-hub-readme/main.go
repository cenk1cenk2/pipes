package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/pipe"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       pipe.Flags,
				Action: func(_ context.Context, c *cli.Command) error {
					return p.RunJobs(
						pipe.New(p).Job(),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
