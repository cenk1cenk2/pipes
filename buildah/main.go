package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/build"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/buildah/setup"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        "pipe-buildah",
			Version:     version,
			Description: "Builds and publishes container images from CI with buildah.io",
			Commands: []*cli.Command{
				{
					Name:        "login",
					Description: "Login to the given container registries.",
					Flags:       CombineFlags(setup.Flags, login.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p),
						))
					},
				},
				{
					Name:        "build",
					Description: "Build container images.",
					Flags:       CombineFlags(setup.Flags, login.Flags, build.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p),
							build.New(p),
						))
					},
				},
				{
					Name:        "manifest",
					Description: "Update manifests of the container images.",
					Flags:       CombineFlags(setup.Flags, login.Flags, manifest.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p),
							manifest.New(p),
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
