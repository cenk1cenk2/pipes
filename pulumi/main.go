// Command pipe-pulumi runs Pulumi actions for CI pipelines.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/pulumi/preview"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
	"gitlab.kilic.dev/devops/pipes/pulumi/up"
)

func newCommand(p *plumber.Plumber) *cli.Command {
	return &cli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Commands: []*cli.Command{
			{
				Name:        "preview",
				Description: "Preview the Pulumi changes.",
				Flags:       plumber.CombineFlags(setup.Flags, stack.Flags, preview.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						stack.New(p),
						preview.New(p),
					))
				},
			},
			{
				Name:        "up",
				Description: "Apply the Pulumi changes.",
				Flags:       plumber.CombineFlags(setup.Flags, stack.Flags, up.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						stack.New(p),
						up.New(p),
					))
				},
			},
		},
	}
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
