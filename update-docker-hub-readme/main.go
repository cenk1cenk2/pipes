// Command pipe-update-docker-hub-readme updates the readme file on DockerHub.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/update"
)

func main() {
	plumber.NewPlumber(func(p *plumber.Plumber) *cli.Command {
		return &cli.Command{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Flags:       plumber.CombineFlags(update.Flags),
			Action: func(_ context.Context, _ *cli.Command) error {
				return p.RunJobs(plumber.CombineTaskLists(
					update.New(p),
				))
			},
		}
	}).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
