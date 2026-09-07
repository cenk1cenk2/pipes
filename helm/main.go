// Command pipe-helm runs Helm operations for CI pipelines.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/helm/install"
	"gitlab.kilic.dev/devops/pipes/helm/lint"
	"gitlab.kilic.dev/devops/pipes/helm/login"
	"gitlab.kilic.dev/devops/pipes/helm/publish"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func newCommand(p *plumber.Plumber) *ucli.Command {
	return &ucli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Commands: []*ucli.Command{
			{
				Name:        "install",
				Description: "Install Helm chart dependencies.",
				Flags:       plumber.CombineFlags(setup.Flags, login.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						login.New(p),
						install.New(p, install.Deps{Tool: setup.C.Ctx}),
					))
				},
			},
			{
				Name:        "lint",
				Description: "Lint Helm chart templates.",
				Flags:       plumber.CombineFlags(setup.Flags, lint.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						lint.New(p, lint.Deps{Tool: setup.C.Ctx}),
					))
				},
			},
			{
				Name:        "publish",
				Description: "Publish Helm chart templates.",
				Flags:       plumber.CombineFlags(setup.Flags, login.Flags, publish.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						login.New(p),
						publish.New(p, publish.Deps{Tool: setup.C}),
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
