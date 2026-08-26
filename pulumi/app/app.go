// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

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

const Name = "pipe-pulumi"

const Description = "Pulumi actions for CI pipelines."

// Options are the services the pipe reaches outside the machine for. A zero
// value is the production wiring, so only a spec ever fills one in.
type Options struct {
	Notes gitlab.NotesFactory
}

func (o Options) defaults() Options {
	if o.Notes == nil {
		o.Notes = gitlab.NewNotes
	}

	return o
}

func New(p *plumber.Plumber, version string, opts Options) *ucli.Command {
	opts = opts.defaults()

	tool := setup.Step
	selected := stack.Step(stack.Deps{Tool: setup.C})

	return cli.App(Name, Description, version,
		cli.Command(p, "preview", "Preview the Pulumi changes.", tool, selected, preview.Step(preview.Deps{Tool: setup.C, Stack: stack.P, Notes: opts.Notes})),
		cli.Command(p, "up", "Apply the Pulumi changes.", tool, selected, up.Step(up.Deps{Tool: setup.C})),
	)
}
