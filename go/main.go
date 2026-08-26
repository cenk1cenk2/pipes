package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/build"
	"gitlab.kilic.dev/devops/pipes/go/install"
	"gitlab.kilic.dev/devops/pipes/go/lint"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/go/tool"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

func main() {
	NewPlumber(
		func(p *Plumber) *ucli.Command {
			// The setup step is not aliased the way the other pipes alias it, since
			// this pipe already has a subcommand named after the tool.
			return cli.App(CLI_NAME, DESCRIPTION, VERSION,
				cli.Command(p, "install", "Vendor go modules.", setup.Step, install.Step(install.Deps{Tool: setup.C})),
				cli.Command(p, "build", "Build an application.", setup.Step, build.Step(build.Deps{Tool: setup.C})),
				cli.Command(p, "lint", "Run golangci-lint on the project.", setup.Step, lint.Step(lint.Deps{Tool: setup.C})),
				cli.Command(p, "tool", "Run a specified go tool.", setup.Step, tool.Step(tool.Deps{Tool: setup.C})),
			)
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
