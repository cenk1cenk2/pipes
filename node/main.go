package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/login"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func main() {
	OverwriteCliFlag(environment.Flags, func(f *cli.BoolFlag) bool {
		return f.Name == "environment.enable"
	}, func(f *cli.BoolFlag) *cli.BoolFlag {
		f.Hidden = false
		f.Value = false

		return f
	})

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
						Description: "Login to the given NPM registries.",
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
						Name:        "install",
						Description: "Install node.js dependencies with the given package manager.",
						Flags:       CombineFlags(setup.Flags, login.Flags, install.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								JobSequence(
									setup.New(p).Job(),
									login.New(p).Job(),
									install.New(p).Job(),
								),
							)
						},
					},

					{
						Name:  "build",
						Flags: CombineFlags(setup.Flags, environment.Flags, build.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								JobSequence(
									setup.New(p).Job(),
									environment.New(p).Job(),
									build.New(p).Job(),
								),
							)
						},
					},

					{
						Name:  "run",
						Flags: CombineFlags(setup.Flags, environment.Flags, run.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								JobSequence(
									setup.New(p).Job(),
									environment.New(p).Job(),
									run.New(p).Job(),
								),
							)
						},
					},
				},
			}
		},
	).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
