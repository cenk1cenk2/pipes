// Command pipe-semantic-release releases applications through the
// semantic-release library.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/semantic-release/release"
	"gitlab.kilic.dev/devops/pipes/semantic-release/setup"
)

const name = "pipe-semantic-release"

const description = "Releases applications through the semantic-release library."

var VERSION = "latest"

// options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type options struct{}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, _ options) *ucli.Command {
	// The environment feature is opt-in for this pipe, unlike the pipes that own
	// their environment. The flags are shared package level values, so this runs
	// before the command tree reads them.
	plumber.OverwriteCliFlag(setup.EnvironmentFlags, func(f *ucli.BoolFlag) bool {
		return f.Name == "environment.enable"
	}, func(f *ucli.BoolFlag) *ucli.BoolFlag {
		f.Hidden = false
		f.Value = false

		return f
	})

	return cli.Root(p, name, description, version,
		setup.EnvironmentStep,
		setup.Step,
		setup.LoginStep,
		release.Step,
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
