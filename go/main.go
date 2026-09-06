// Command pipe-go builds Go applications from CI/CD.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/build"
	"gitlab.kilic.dev/devops/pipes/go/install"
	"gitlab.kilic.dev/devops/pipes/go/lint"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	gotool "gitlab.kilic.dev/devops/pipes/go/tool"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

const name = "pipe-go"

const description = "Build Go applications with the CI pipe."

var VERSION = "latest"

// options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type options struct{}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, _ options) *ucli.Command {
	// The setup step is not aliased the way the other pipes alias it, since this
	// pipe already has a subcommand named after the tool.
	return cli.App(name, description, version,
		cli.Command(p, "install", "Vendor go modules.", setup.Step, install.Step(install.Deps{Tool: setup.C})),
		cli.Command(p, "build", "Build an application.", setup.Step, build.Step(build.Deps{Tool: setup.C.Ctx})),
		cli.Command(p, "lint", "Run golangci-lint on the project.", setup.Step, lint.Step(lint.Deps{Tool: setup.C})),
		cli.Command(p, "tool", "Run a specified go tool.", setup.Step, gotool.Step(gotool.Deps{Tool: setup.C.Ctx})),
	)
}

func main() {
	plumber.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, VERSION, options{})
	}).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
