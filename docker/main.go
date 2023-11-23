package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/docker/build"
	"gitlab.kilic.dev/devops/pipes/docker/login"
	"gitlab.kilic.dev/devops/pipes/docker/manifest"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Commands: cli.Commands{
					{
						Name:        "login",
						Description: "Login to the given Docker registries.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags),
						Action: func(c *cli.Context) error {
							tl := &login.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:        "build",
						Description: "Build Docker images.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, build.Flags),
						Action: func(c *cli.Context) error {
							tl := &build.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
									build.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:        "manifest",
						Description: "Update manifests of the Docker images.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, manifest.Flags),
						Action: func(c *cli.Context) error {
							tl := &manifest.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
									manifest.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags:       true,
			ExcludeHelpCommand: true,
		}).
		Run()
}
