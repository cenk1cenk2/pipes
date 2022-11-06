package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/login"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func main() {
	p := Plumber{
		DocsExcludeFlags:       true,
		DocsExcludeHelpCommand: true,
	}

	p.New(
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
							return login.TL.RunJobs(
								login.TL.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:        "install",
						Description: "Install node.js dependencies with the given package manager.",
						Usage:       fmt.Sprintf("%s install", CLI_NAME),
						Flags:       p.AppendFlags(setup.Flags, login.Flags, install.Flags),
						Action: func(c *cli.Context) error {
							return install.TL.RunJobs(
								install.TL.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
									install.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "build",
						Flags: p.AppendFlags(setup.Flags, build.Flags),
						Action: func(c *cli.Context) error {
							return build.TL.RunJobs(
								build.TL.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									build.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "run",
						Flags: p.AppendFlags(setup.Flags, run.Flags),
						Action: func(c *cli.Context) error {
							return run.TL.RunJobs(
								run.TL.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									run.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},
				},
			}
		},
	).Run()
}
