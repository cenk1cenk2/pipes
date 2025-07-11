package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/docker/build"
	"gitlab.kilic.dev/devops/pipes/docker/login"
	"gitlab.kilic.dev/devops/pipes/docker/manifest"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
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
						Name:        "login",
						Description: "Login to the given Docker registries.",
						Flags:       CombineFlags(setup.Flags, login.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								JobSequence(
									setup.New(p).Job(),
									login.New(p).Job(),
								),
							)
						},
					},

					{
						Name:        "build",
						Description: "Build Docker images.",
						Flags:       CombineFlags(setup.Flags, login.Flags, build.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								JobSequence(
									setup.New(p).Job(),
									login.New(p).Job(),
									build.New(p).Job(),
								),
							)
						},
					},

					{
						Name:        "manifest",
						Description: "Update manifests of the Docker images.",
						Flags:       CombineFlags(setup.Flags, login.Flags, manifest.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								JobSequence(
									setup.New(p).Job(),
									login.New(p).Job(),
									manifest.New(p).Job(),
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
