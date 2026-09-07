// Command select-env selects a set of environment variable prefixes from the
// conditions.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	"gitlab.kilic.dev/devops/pipes/select-env/write"
)

func newCommand(p *plumber.Plumber) *ucli.Command {
	return cli.Root(p, CLI_NAME, DESCRIPTION, VERSION,
		setup.Step,
		write.Step(write.Deps{Environment: setup.EnvironmentCtx}),
	)
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
