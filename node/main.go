// Command pipe-node installs dependencies and builds Node.js applications from
// CI/CD.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func newCommand(p *plumber.Plumber) *cli.Command {
	// The environment feature is opt-in for this pipe, unlike the pipes that own
	// their environment. The flags are shared package level values, so this runs
	// before the command tree reads them.
	plumber.OverwriteCliFlag(setup.EnvironmentFlags, func(f *cli.BoolFlag) bool {
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
				Flags:       plumber.CombineFlags(setup.NodeFlags, setup.LoginFlags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						setup.NewLogin(p),
					))
				},
			},
			{
				Name:        "install",
				Description: "Install node.js dependencies with the given package manager.",
				Flags:       plumber.CombineFlags(setup.NodeFlags, setup.LoginFlags, install.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						setup.NewLogin(p),
						install.New(p, install.Deps{Node: setup.NodeCtx, Environment: setup.EnvironmentCtx}),
					))
				},
			},
			{
				Name:        "build",
				Description: "",
				Flags:       plumber.CombineFlags(setup.NodeFlags, setup.EnvironmentFlags, build.Flags),
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						setup.NewEnvironment(p),
						build.New(p, build.Deps{Node: setup.NodeCtx, Environment: setup.EnvironmentCtx}),
					))
				},
			},
			{
				Name:        "run",
				Description: "",
				Flags:       plumber.CombineFlags(setup.NodeFlags, setup.EnvironmentFlags, run.Flags),
				Arguments:   run.Arguments,
				Action: func(_ context.Context, _ *cli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						setup.NewEnvironment(p),
						run.New(p, run.Deps{Node: setup.NodeCtx, Environment: setup.EnvironmentCtx}),
					))
				},
			},
		},
	}
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
