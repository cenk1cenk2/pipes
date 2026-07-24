package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/kustomize/build"
	setup "gitlab.kilic.dev/devops/pipes/kustomize/setup"
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
						Name:        "build",
						Description: "Build and validate Kustomize overlays.",
						Flags: CombineFlags(
							setup.Flags,
							build.Flags,
						),
						Action: func(_ context.Context, c *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									build.New(p),
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
