package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/pipes/node/setup"
)

func main() {
	OverwriteCliFlag(setup.EnvironmentFlags, func(f *cli.BoolFlag) bool {
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
						Flags:       CombineFlags(setup.NodeFlags, setup.LoginFlags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									node.SetupTaskList(p, setup.NodeConfig, setup.NodeCtx),
									node.LoginTaskList(p, setup.Login),
								),
							)
						},
					},

					{
						Name:        "install",
						Description: "Install node.js dependencies with the given package manager.",
						Flags:       CombineFlags(setup.NodeFlags, setup.LoginFlags, install.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									node.SetupTaskList(p, setup.NodeConfig, setup.NodeCtx),
									node.LoginTaskList(p, setup.Login),
									install.New(p),
								),
							)
						},
					},

					{
						Name:  "build",
						Flags: CombineFlags(setup.NodeFlags, setup.EnvironmentFlags, build.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									node.SetupTaskList(p, setup.NodeConfig, setup.NodeCtx),
									environment.TaskList(p, setup.Environment, setup.EnvironmentCtx),
									build.New(p),
								),
							)
						},
					},

					{
						Name:      "run",
						Flags:     CombineFlags(setup.NodeFlags, setup.EnvironmentFlags, run.Flags),
						Arguments: CombineArguments(run.Arguments),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									node.SetupTaskList(p, setup.NodeConfig, setup.NodeCtx),
									environment.TaskList(p, setup.Environment, setup.EnvironmentCtx),
									run.New(p),
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
