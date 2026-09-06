// Command pipe-pulumi runs Pulumi actions for CI pipelines.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/pulumi/preview"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
	"gitlab.kilic.dev/devops/pipes/pulumi/up"
)

const name = "pipe-pulumi"

const description = "Pulumi actions for CI pipelines."

var VERSION = "latest"

// options are the services the pipe reaches outside the machine for. A zero
// value is the production wiring, so only a spec ever fills one in.
type options struct {
	Notes gitlab.NotesFactory
}

func (o options) defaults() options {
	if o.Notes == nil {
		o.Notes = gitlab.NewNotes
	}

	return o
}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, opts options) *ucli.Command {
	opts = opts.defaults()

	tool := setup.Step
	selected := stack.Step(stack.Deps{Tool: setup.C})

	return cli.App(name, description, version,
		cli.Command(p, "preview", "Preview the Pulumi changes.", tool, selected, preview.Step(preview.Deps{Tool: setup.C, Stack: stack.P, Notes: opts.Notes})),
		cli.Command(p, "up", "Apply the Pulumi changes.", tool, selected, up.Step(up.Deps{Tool: setup.C})),
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
