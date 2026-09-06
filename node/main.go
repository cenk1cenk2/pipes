// Command pipe-node installs dependencies and builds Node.js applications from
// CI/CD.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/node/build"
	"gitlab.kilic.dev/devops/pipes/node/install"
	"gitlab.kilic.dev/devops/pipes/node/run"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

const name = "pipe-node"

const description = "Pipe for installing node.js dependencies and building node.js applications on CI/CD."

var VERSION = "latest"

func newCommand(p *plumber.Plumber) *ucli.Command {
	// The environment feature is opt-in for this pipe, unlike the pipes that own
	// their environment. The flags are shared package level values, so this runs
	// before the command tree reads them.
	plumber.OverwriteCliFlag(setup.EnvironmentFlags, func(f *ucli.BoolFlag) bool {
		return f.Name == "environment.enable"
	}, func(f *ucli.BoolFlag) *ucli.BoolFlag {
		f.Hidden = false
		f.Value = false

		return f
	})

	return cli.App(name, description, VERSION,
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
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
