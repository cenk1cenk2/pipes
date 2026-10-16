package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/build"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/buildah/setup"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		container := registry.NewLogin(registry.LoginOptions{
			Category: login.CategoryContainerRegistry,
			Alias:    "BUILDAH_CONTAINER_LOGIN",
			Command:  []string{"buildah", "login"},
		})

		return &cli.Command{
			Name:        "pipe-buildah",
			Version:     version,
			Description: "Builds and publishes container images from CI with buildah.io",
			Commands: []*cli.Command{
				{
					Name:        "login",
					Description: "Login to the given container registries.",
					Flags:       CombineFlags(setup.Flags, login.Flags, container.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p, container),
							container.New(p),
						))
					},
				},
				{
					Name:        "build",
					Description: "Build container images.",
					Flags:       CombineFlags(setup.Flags, login.Flags, container.Flags, build.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p, container),
							container.New(p),
							build.New(p),
						))
					},
				},
				{
					Name:        "manifest",
					Description: "Update manifests of the container images.",
					Flags:       CombineFlags(setup.Flags, login.Flags, container.Flags, manifest.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p, container),
							container.New(p),
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
