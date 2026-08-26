package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/template/pipe"
)

const Name = "pipe-template"

const Description = "template-cli"

var VERSION = "latest"

func main() {
	NewPlumber(
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        Name,
				Version:     VERSION,
				Usage:       Description,
				Description: Description,
				Flags:       CombineFlags(pipe.Flags),
				Action: func(_ context.Context, c *cli.Command) error {
					return p.RunJobs(
						CombineTaskLists(
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
