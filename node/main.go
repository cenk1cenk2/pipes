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
	p := Plumber{
		DocsExcludeFlags:       true,
		DocsExcludeHelpCommand: true,
		DeprecationNotices: []DeprecationNotice{
			{
				Flag:        []string{"--node.build_environment_files", "--node.build_environment_fallback", "--node.build_environment_conditions"},
				Environment: []string{"NODE_BUILD_ENVIRONMENT_FILES", "NODE_BUILD_ENVIRONMENT_CONDITIONS", "NODE_BUILD_ENVIRONMENT_FALLBACK"},
				Level:       LOG_LEVEL_ERROR,
				Message:     `"%s" is deprecated, please utilize the new select-env flags instead.`,
			},
		},
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
						Flags: p.AppendFlags(environment.Flags, setup.Flags, build.Flags),
						Action: func(c *cli.Context) error {
							return build.TL.RunJobs(
								build.TL.JobSequence(
									build.TL.JobIf(
										build.TL.Predicate(func(tl *TaskList[build.Pipe]) bool {
											return tl.Pipe.Environment.Enable
										}),
										environment.New(p).SetCliContext(c).Job(),
									),
									setup.New(p).SetCliContext(c).Job(),
									build.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:  "run",
						Flags: p.AppendFlags(environment.Flags, setup.Flags, run.Flags),
						Action: func(c *cli.Context) error {
							return run.TL.RunJobs(
								run.TL.JobSequence(
									run.TL.JobIf(
										run.TL.Predicate(func(tl *TaskList[run.Pipe]) bool {
											return tl.Pipe.Environment.Enable
										}),
										environment.New(p).SetCliContext(c).Job(),
									),
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
