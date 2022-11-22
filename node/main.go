package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/login"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Commands: cli.Commands{
					{
						Name:        "login",
						Description: "Login to the given NPM registries.",
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
						Name:        "install",
						Description: "Install node.js dependencies with the given package manager.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, install.Flags),
						Action: func(c *cli.Context) error {
							tl := &install.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
									install.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "build",
						Flags: p.AppendFlags(setup.Flags, environment.Flags, build.Flags),
						Before: func(ctx *cli.Context) error {
							p.SetDeprecationNotices(build.DeprecationNotices)

							return nil
						},
						Action: func(c *cli.Context) error {
							tl := &build.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									environment.New(p).SetCliContext(c).Job(),
									build.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "run",
						Flags: p.AppendFlags(setup.Flags, environment.Flags, run.Flags),
						Action: func(c *cli.Context) error {
							tl := &run.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									environment.New(p).SetCliContext(c).Job(),
									run.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},
				},
			}
		},
	).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags:       true,
			ExcludeHelpCommand: true,
		}).
		Run()
}
