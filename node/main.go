package main

import (
	"github.com/urfave/cli/v2"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	build "gitlab.kilic.dev/devops/gitlab-pipes/node/build"
	install "gitlab.kilic.dev/devops/gitlab-pipes/node/install"
	login "gitlab.kilic.dev/devops/gitlab-pipes/node/login"
	pipe "gitlab.kilic.dev/devops/gitlab-pipes/node/pipe"
	run "gitlab.kilic.dev/devops/gitlab-pipes/node/run"
)

func main() {
	utils.CliCreate(
		&cli.App{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Action: func(c *cli.Context) error {
				utils.CliGreet(c)

				err := cli.ShowAppHelp(c)

				if err != nil {
					return err
				}

				utils.Log.Fatalln("Need a subcommand to run!")

				return nil
			},
			Commands: cli.Commands{
				{
					Name: "login",
					Flags: append(
						append(utils.CliDefaultFlags, pipe.Flags...),
						login.Flags...),
					Action: func(c *cli.Context) error {
						utils.CliGreet(c)

						if err := pipe.Pipe.Exec(); err != nil {
							return err
						}

						return login.Pipe.Exec()
					},
				},

				{
					Name: "install",
					Flags: append(
						append(append(utils.CliDefaultFlags, pipe.Flags...), install.Flags...),
						login.Flags...),
					Action: func(c *cli.Context) error {
						utils.CliGreet(c)

						if err := pipe.Pipe.Exec(); err != nil {
							return err
						}

						if err := login.Pipe.Exec(); err != nil {
							return err
						}

						return install.Pipe.Exec()
					},
				},

				{
					Name:  "build",
					Flags: append(append(utils.CliDefaultFlags, pipe.Flags...), build.Flags...),
					Action: func(c *cli.Context) error {
						utils.CliGreet(c)

						if err := pipe.Pipe.Exec(); err != nil {
							return err
						}

						return build.Pipe.Exec()
					},
				},

				{
					Name:  "run",
					Flags: append(append(utils.CliDefaultFlags, pipe.Flags...), run.Flags...),
					Action: func(c *cli.Context) error {
						utils.CliGreet(c)

						if err := pipe.Pipe.Exec(); err != nil {
							return err
						}

						return run.Pipe.Exec(c)
					},
				},
			},
		},
	)
}
