// Command pipe-template is the scaffold a new pipe is copied from.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/template/pipe"
)

func newCommand(p *plumber.Plumber) *ucli.Command {
	return &ucli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Flags:       plumber.CombineFlags(pipe.Flags),
		Action: func(_ context.Context, _ *ucli.Command) error {
			return p.RunJobs(plumber.CombineTaskLists(
				pipe.New(p),
			))
		},
	}
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
