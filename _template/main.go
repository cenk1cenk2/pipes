// Command pipe-template is the scaffold a new pipe is copied from.
package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/template/pipe"
)

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        CLIName,
			Version:     version,
			Usage:       Description,
			Description: Description,
			Flags:       CombineFlags(pipe.Flags),
			Action: func(_ context.Context, _ *cli.Command) error {
				return p.RunJobs(CombineTaskLists(
					pipe.New(p),
				))
			},
		}
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
