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
						Flags:       p.AppendFlags(setup.Flags, login.Flags),
						Action: func(ctx context.Context, command *cli.Command) error {
							tl := &login.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCli(command).Job(),
									login.New(p).SetCli(command).Job(),
								),
							)
						},
					},

					{
						Name:        "build",
						Description: "Build Docker images.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, build.Flags),
						Action: func(ctx context.Context, command *cli.Command) error {
							tl := &build.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCli(command).Job(),
									login.New(p).SetCli(command).Job(),
									build.New(p).SetCli(command).Job(),
								),
							)
						},
					},

					{
						Name:        "manifest",
						Description: "Update manifests of the Docker images.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, manifest.Flags),
						Action: func(ctx context.Context, command *cli.Command) error {
							tl := &manifest.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCli(command).Job(),
									login.New(p).SetCli(command).Job(),
									manifest.New(p).SetCli(command).Job(),
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
