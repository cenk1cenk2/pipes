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

const name = "select-env"

const description = "Selects an set of environment variable prefix depending on the condition."

var VERSION = "latest"

// options is where a service the pipe reaches outside the machine for would be
// injected. This pipe only rewrites its own environment, so there is nothing to
// swap.
type options struct{}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, _ options) *ucli.Command {
	return cli.Root(p, name, description, version,
		setup.Step,
		write.Step(write.Deps{Environment: setup.EnvironmentCtx}),
	)
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
