// Command select-env selects a set of environment variable prefixes from the
// conditions.
package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	"gitlab.kilic.dev/devops/pipes/select-env/write"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        "select-env",
			Version:     version,
			Description: "Selects an set of environment variable prefix depending on the condition.",
			Flags:       CombineFlags(setup.Flags, write.Flags),
			Action: func(_ context.Context, _ *cli.Command) error {
				return p.RunJobs(CombineTaskLists(
					setup.New(p),
					write.New(p),
				))
			},
		}
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
