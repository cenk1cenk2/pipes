// Command pipe-kustomize runs Kustomize operations for CI pipelines.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/kustomize/build"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

func newCommand(p *plumber.Plumber) *ucli.Command {
	return &ucli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Commands: []*ucli.Command{
			{
				Name:        "build",
				Description: "Build and validate Kustomize overlays.",
				Flags:       plumber.CombineFlags(setup.Flags, build.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						build.New(p, build.Deps{Tool: setup.C}),
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
