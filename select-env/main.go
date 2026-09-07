// Command select-env selects a set of environment variable prefixes from the
// conditions.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	"gitlab.kilic.dev/devops/pipes/select-env/write"
)

func newCommand(p *plumber.Plumber) *cli.Command {
	return &cli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Flags:       plumber.CombineFlags(setup.Flags, write.Flags),
		Action: func(_ context.Context, _ *cli.Command) error {
			return p.RunJobs(plumber.CombineTaskLists(
				setup.New(p),
				write.New(p, write.Deps{Environment: setup.EnvironmentCtx}),
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
