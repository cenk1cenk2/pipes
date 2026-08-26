// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	"gitlab.kilic.dev/devops/pipes/select-env/write"
)

const Name = "select-env"

const Description = "Selects an set of environment variable prefix depending on the condition."

// Options is where a service the pipe reaches outside the machine for would be
// injected. This pipe only rewrites its own environment, so there is nothing to
// swap.
type Options struct{}

func New(p *plumber.Plumber, version string, _ Options) *ucli.Command {
	return cli.Root(p, Name, Description, version,
		setup.Step,
		write.Step(write.Deps{Environment: setup.EnvironmentCtx}),
	)
}
