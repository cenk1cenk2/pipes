package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/select-env/pipe"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       CombineFlags(setup.Flags, pipe.Flags),
				Action: func(_ context.Context, c *cli.Command) error {
					return p.RunJobs(
						CombineTaskLists(
							setup.New(p),
							pipe.New(p),
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
