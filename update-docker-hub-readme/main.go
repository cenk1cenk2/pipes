package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/update"
)

func main() {
	NewPlumber(
		func(p *Plumber) *ucli.Command {
			return cli.Root(p, CLI_NAME, DESCRIPTION, VERSION,
				update.Step(update.Deps{Hub: hub.NewClient}),
			)
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
