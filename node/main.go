package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/login"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func main() {
	p := Plumber{}

	p.New(
		func(a *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Commands: cli.Commands{
					{
						Name:  "login",
						Flags: a.AppendFlags(setup.Flags, login.Flags),
						Action: func(c *cli.Context) error {
							return login.TL.RunJobs(
								login.TL.JobSequence(
									setup.New(a).SetCliContext(c).Job(),
									login.New(a).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "install",
						Flags: a.AppendFlags(setup.Flags, login.Flags, install.Flags),
						Action: func(c *cli.Context) error {
							return install.TL.RunJobs(
								install.TL.JobSequence(
									setup.New(a).SetCliContext(c).Job(),
									login.New(a).SetCliContext(c).Job(),
									install.New(a).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "build",
						Flags: a.AppendFlags(setup.Flags, build.Flags),
						Action: func(c *cli.Context) error {
							return build.TL.RunJobs(
								build.TL.JobSequence(
									setup.New(a).SetCliContext(c).Job(),
									build.New(a).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "run",
						Flags: a.AppendFlags(setup.Flags, run.Flags),
						Action: func(c *cli.Context) error {
							return run.TL.RunJobs(
								run.TL.JobSequence(
									setup.New(a).SetCliContext(c).Job(),
									run.New(a).SetCliContext(c).Job(),
								),
							)
						},
					},
				},
			}
		},
	).Run()
}
