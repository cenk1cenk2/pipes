// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/update"
)

const Name = "pipe-update-docker-hub-readme"

const Description = "Updates the readme file on DockerHub or any compatible API."

// Options are the services the pipe reaches outside the machine for. A zero
// value is the production wiring, so only a spec ever fills one in.
type Options struct {
	Hub hub.ClientFactory
}

func (o Options) defaults() Options {
	if o.Hub == nil {
		o.Hub = hub.NewClient
	}

	return o
}

func New(p *plumber.Plumber, version string, opts Options) *ucli.Command {
	opts = opts.defaults()

	return cli.Root(p, Name, Description, version,
		update.Step(update.Deps{Hub: opts.Hub}),
	)
}
