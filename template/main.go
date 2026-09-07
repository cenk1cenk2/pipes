// Command pipe-template is the scaffold a new pipe is copied from.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/template/pipe"
)

func main() {
	plumber.NewPlumber(func(p *plumber.Plumber) *cli.Command {
		return &cli.Command{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Flags:       plumber.CombineFlags(pipe.Flags),
			// The task lists are built in here rather than alongside the flags, since a
			// stage reads the parsed flag values as it constructs.
			Action: func(_ context.Context, _ *cli.Command) error {
				return p.RunJobs(plumber.CombineTaskLists(
					pipe.New(p),
				))
			},
		}
	}).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
