package main

import (
	"context"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/go/build"
	"gitlab.kilic.dev/devops/pipes/go/install"
	"gitlab.kilic.dev/devops/pipes/go/lint"
	"gitlab.kilic.dev/devops/pipes/go/setup"

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
						Name:        "install",
						Description: "Vendor go modules.",
						Flags:       CombineFlags(setup.Flags, install.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									install.New(p),
								),
							)
						},
					},

					{
						Name:        "lint",
						Description: "Lint the repository.",
						Flags:       CombineFlags(setup.Flags, lint.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									lint.New(p),
								),
							)
						},
					},

					{
						Name:        "build",
						Description: "Build an application.",
						Flags:       CombineFlags(setup.Flags, build.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
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
