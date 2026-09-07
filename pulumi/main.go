// Command pipe-pulumi runs Pulumi actions for CI pipelines.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/pulumi/preview"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
	"gitlab.kilic.dev/devops/pipes/pulumi/up"
)

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

func newCommand(p *plumber.Plumber, opts options) *ucli.Command {
	opts = opts.defaults()

	return &ucli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Commands: []*ucli.Command{
			{
				Name:        "preview",
				Description: "Preview the Pulumi changes.",
				Flags:       plumber.CombineFlags(setup.Flags, stack.Flags, preview.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						stack.New(p, stack.Deps{Tool: setup.C}),
						preview.New(p, preview.Deps{Tool: setup.C, Stack: stack.P, Notes: opts.Notes}),
					))
				},
			},
			{
				Name:        "up",
				Description: "Apply the Pulumi changes.",
				Flags:       plumber.CombineFlags(setup.Flags, stack.Flags, up.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						stack.New(p, stack.Deps{Tool: setup.C}),
						up.New(p, up.Deps{Tool: setup.C}),
					))
				},
			},
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
