// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

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

const Name = "pipe-go"

const Description = "Build Go applications with the CI pipe."

// Options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type Options struct{}

func New(p *plumber.Plumber, version string, _ Options) *ucli.Command {
	// The setup step is not aliased the way the other pipes alias it, since this
	// pipe already has a subcommand named after the tool.
	return cli.App(Name, Description, version,
		cli.Command(p, "install", "Vendor go modules.", setup.Step, install.Step(install.Deps{Tool: setup.C})),
		cli.Command(p, "build", "Build an application.", setup.Step, build.Step(build.Deps{Tool: setup.C})),
		cli.Command(p, "lint", "Run golangci-lint on the project.", setup.Step, lint.Step(lint.Deps{Tool: setup.C})),
		cli.Command(p, "tool", "Run a specified go tool.", setup.Step, gotool.Step(gotool.Deps{Tool: setup.C})),
	)
}
