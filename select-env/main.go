package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	"gitlab.kilic.dev/devops/pipes/select-env/write"
)

func main() {
	NewPlumber(
		func(p *Plumber) *ucli.Command {
			return cli.Root(p, CLI_NAME, DESCRIPTION, VERSION,
				setup.Step,
				write.Step(write.Deps{Environment: setup.EnvironmentCtx}),
			)
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
