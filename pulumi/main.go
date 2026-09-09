// Command pipe-pulumi runs Pulumi actions for CI pipelines.
package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/pulumi/preview"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
	"gitlab.kilic.dev/devops/pipes/pulumi/up"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        "pipe-pulumi",
			Version:     version,
			Description: "Pulumi related tasks in the pipeline.",
			Commands: []*cli.Command{
				{
					Name:        "preview",
					Description: "Preview the Pulumi changes.",
					Flags:       CombineFlags(setup.Flags, stack.Flags, preview.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							stack.New(p),
							preview.New(p),
						))
					},
				},
				{
					Name:        "up",
					Description: "Apply the Pulumi changes.",
					Flags:       CombineFlags(setup.Flags, stack.Flags, up.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							stack.New(p),
							up.New(p),
						))
					},
				},

				DocsCommand(p),
			},
		}
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
