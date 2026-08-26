package main

import (
	"context"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/pulumi/preview"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
	"gitlab.kilic.dev/devops/pipes/pulumi/up"

	. "github.com/cenk1cenk2/plumber/v6"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Commands: []*cli.Command{
					{
						Name:        "preview",
						Description: "Preview the Pulumi changes.",
						Flags:       CombineFlags(setup.Flags, stack.Flags, preview.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									stack.New(p),
									preview.New(p),
								),
							)
						},
					},

					{
						Name:        "up",
						Description: "Apply the Pulumi changes.",
						Flags:       CombineFlags(setup.Flags, stack.Flags, up.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									stack.New(p),
									up.New(p),
								),
							)
						},
					},
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
