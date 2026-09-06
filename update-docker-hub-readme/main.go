// Command pipe-update-docker-hub-readme updates the readme file on DockerHub.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/update"
)

const name = "pipe-update-docker-hub-readme"

const description = "Updates the readme file on DockerHub or any compatible API."

var VERSION = "latest"

// options are the services the pipe reaches outside the machine for. A zero
// value is the production wiring, so only a spec ever fills one in.
type options struct {
	Hub hub.ClientFactory
}

func (o options) defaults() options {
	if o.Hub == nil {
		o.Hub = hub.NewClient
	}

	return o
}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, opts options) *ucli.Command {
	opts = opts.defaults()

	return cli.Root(p, name, description, version,
		update.Step(update.Deps{Hub: opts.Hub}),
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
