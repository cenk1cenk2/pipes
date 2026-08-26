package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/semantic-release/release"
	"gitlab.kilic.dev/devops/pipes/semantic-release/setup"
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
				Flags:       CombineFlags(setup.EnvironmentFlags, setup.NodeFlags, setup.LoginFlags, release.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(
						CombineTaskLists(
							environment.TaskList(p, setup.Environment, setup.EnvironmentCtx),
							node.SetupTaskList(p, setup.NodeConfig, setup.NodeCtx),
							node.LoginTaskList(p, setup.Login),
							release.New(p),
						),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
