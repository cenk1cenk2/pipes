// Command pipe-node installs dependencies and builds Node.js applications from
// CI/CD.
package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/node/add"
	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/login"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		// the environment is opt-in here, and the flags are shared package level values,
		// so this runs before the command tree reads them.
		OverwriteCliFlag(setup.EnvironmentFlags, func(f *cli.BoolFlag) bool {
			return f.Name == "environment.enable"
		}, func(f *cli.BoolFlag) *cli.BoolFlag {
			f.Hidden = false
			f.Value = false

			return f
		})

		return &cli.Command{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Commands: []*cli.Command{
				{
					Name:        "login",
					Description: "Login to the given NPM registries.",
					Flags:       CombineFlags(setup.NodeFlags, login.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p),
						))
					},
				},
				{
					Name:        "install",
					Description: "Install node.js dependencies with the given package manager.",
					Flags:       CombineFlags(setup.NodeFlags, login.Flags, install.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							login.New(p),
							install.New(p),
						))
					},
				},
				{
					Name:        "add",
					Description: "Install node packages with the given package manager.",
					Flags:       CombineFlags(setup.NodeFlags, add.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							add.New(p),
						))
					},
				},
				{
					Name:        "build",
					Description: "",
					Flags:       CombineFlags(setup.NodeFlags, setup.EnvironmentFlags, build.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							setup.NewEnvironment(p),
							build.New(p),
						))
					},
				},
				{
					Name:        "run",
					Description: "",
					Flags:       CombineFlags(setup.NodeFlags, setup.EnvironmentFlags, run.Flags),
					Arguments:   run.Arguments,
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							setup.NewEnvironment(p),
							run.New(p),
						))
					},
				},
			},
		}
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
