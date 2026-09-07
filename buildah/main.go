// Command pipe-buildah builds and publishes container images from CI/CD.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/build"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/buildah/setup"
)

func newCommand(p *plumber.Plumber) *cli.Command {
	return &cli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Commands: []*cli.Command{
			{
				Name:        "login",
				Description: "Login to the given container registries.",
				Flags:       plumber.CombineFlags(setup.Flags, login.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						login.New(p),
					))
				},
			},
			{
				Name:        "build",
				Description: "Build container images.",
				Flags:       plumber.CombineFlags(setup.Flags, login.Flags, build.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						login.New(p),
						build.New(p, build.Deps{Registry: login.P}),
					))
				},
			},
			{
				Name:        "manifest",
				Description: "Update manifests of the container images.",
				Flags:       plumber.CombineFlags(setup.Flags, login.Flags, manifest.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						login.New(p),
						manifest.New(p, manifest.Deps{Registry: login.P}),
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
