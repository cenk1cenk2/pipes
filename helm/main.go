// Command pipe-helm runs Helm operations for CI pipelines.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/helm/install"
	"gitlab.kilic.dev/devops/pipes/helm/lint"
	"gitlab.kilic.dev/devops/pipes/helm/login"
	"gitlab.kilic.dev/devops/pipes/helm/publish"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func newCommand(p *plumber.Plumber) *cli.Command {
	return &cli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Commands: []*cli.Command{
			{
				Name:        "install",
				Description: "Install Helm chart dependencies.",
				Flags:       plumber.CombineFlags(setup.Flags, login.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						login.New(p),
						install.New(p),
					))
				},
			},
			{
				Name:        "lint",
				Description: "Lint Helm chart templates.",
				Flags:       plumber.CombineFlags(setup.Flags, lint.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						lint.New(p),
					))
				},
			},
			{
				Name:        "publish",
				Description: "Publish Helm chart templates.",
				Flags:       plumber.CombineFlags(setup.Flags, login.Flags, publish.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						login.New(p),
						publish.New(p),
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
