package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	install "gitlab.kilic.dev/devops/pipes/helm/install"
	"gitlab.kilic.dev/devops/pipes/helm/lint"
	"gitlab.kilic.dev/devops/pipes/helm/login"
	"gitlab.kilic.dev/devops/pipes/helm/publish"
	setup "gitlab.kilic.dev/devops/pipes/helm/setup"
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
						Description: "Install Helm chart dependencies.",
						Flags: CombineFlags(
							setup.Flags,
							login.Flags,
							install.Flags,
						),
						Action: func(_ context.Context, c *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									login.New(p),
									install.New(p),
								),
							)
						},
					},

					{
						Name:        "lint",
						Description: "Lint Helm chart templates.",
						Flags: CombineFlags(
							setup.Flags,
							lint.Flags,
						),
						Action: func(_ context.Context, c *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									lint.New(p),
								),
							)
						},
					},

					{
						Name:        "publish",
						Description: "Publish Helm chart templates.",
						Flags: CombineFlags(
							setup.Flags,
							login.Flags,
							publish.Flags,
						),
						Action: func(_ context.Context, c *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									login.New(p),
									publish.New(p),
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
