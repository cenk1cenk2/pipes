// Command pipe-template is the scaffold a new pipe is copied from.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/template/pipe"
)

const name = "pipe-template"

const description = "template-cli"

var VERSION = "latest"

// options is where a service the pipe reaches outside the machine for would be
// injected. This scaffold reaches for nothing, so there is nothing to swap.
type options struct{}

// newCommand builds the command tree. The version is a parameter so a spec can
// build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, _ options) *ucli.Command {
	return &ucli.Command{
		Name:        name,
		Version:     version,
		Usage:       description,
		Description: description,
		Flags:       plumber.CombineFlags(pipe.Flags),
		Action: func(_ context.Context, _ *ucli.Command) error {
			return p.RunJobs(
				plumber.CombineTaskLists(
					pipe.New(p),
				),
			)
		},
	}
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
