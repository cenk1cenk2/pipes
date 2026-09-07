// Command pipe-kustomize runs Kustomize operations for CI pipelines.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/kustomize/build"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

func main() {
	plumber.NewPlumber(func(p *plumber.Plumber) *cli.Command {
		return &cli.Command{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Commands: []*cli.Command{
				{
					Name:        "build",
					Description: "Build and validate Kustomize overlays.",
					Flags:       plumber.CombineFlags(setup.Flags, build.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(plumber.CombineTaskLists(
							setup.New(p),
							build.New(p),
						))
					},
				},
			},
		}
	}).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
