package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	pipe "gitlab.kilic.dev/devops/pipes/_template/pipe"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       CombineFlags(pipe.Flags),
				Action: func(_ context.Context, c *cli.Command) error {
					return p.RunJobs(
						JobSequence(
							pipe.New(p).Job(),
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
