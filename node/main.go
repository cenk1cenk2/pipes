package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func main() {
	OverwriteCliFlag(setup.EnvironmentFlags, func(f *ucli.BoolFlag) bool {
		return f.Name == "environment.enable"
	}, func(f *ucli.BoolFlag) *ucli.BoolFlag {
		f.Hidden = false
		f.Value = false

		return f
	})

	NewPlumber(
		func(p *Plumber) *ucli.Command {
			return cli.App(CLI_NAME, DESCRIPTION, VERSION,
				cli.Command(p, "login", "Login to the given NPM registries.",
					setup.Step,
					setup.LoginStep,
				),

				cli.Command(p, "install", "Install node.js dependencies with the given package manager.",
					setup.Step,
					setup.LoginStep,
					install.Step(install.Deps{Node: setup.NodeCtx, Environment: setup.EnvironmentCtx}),
				),

				cli.Command(p, "build", "",
					setup.Step,
					setup.EnvironmentStep,
					build.Step(build.Deps{Node: setup.NodeCtx, Environment: setup.EnvironmentCtx}),
				),

				cli.Command(p, "run", "",
					setup.Step,
					setup.EnvironmentStep,
					run.Step(run.Deps{Node: setup.NodeCtx, Environment: setup.EnvironmentCtx}),
				),
			)
		},
	).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
