// Command pipe-update-docker-hub-readme updates the readme file on DockerHub.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/update"
)

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

func newCommand(p *plumber.Plumber, opts options) *ucli.Command {
	opts = opts.defaults()

	return &ucli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Flags:       plumber.CombineFlags(update.Flags),
		Action: func(_ context.Context, _ *ucli.Command) error {
			return p.RunJobs(plumber.CombineTaskLists(
				update.New(p, update.Deps{Hub: opts.Hub}),
			))
		},
	}
}

func main() {
	plumber.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, options{})
	}).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
