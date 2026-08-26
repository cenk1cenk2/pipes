// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/semantic-release/release"
	"gitlab.kilic.dev/devops/pipes/semantic-release/setup"
)

const Name = "pipe-semantic-release"

const Description = "Releases applications through the semantic-release library."

// Options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type Options struct{}

func New(p *plumber.Plumber, version string, _ Options) *ucli.Command {
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

	return cli.Root(p, Name, Description, version,
		setup.EnvironmentStep,
		setup.Step,
		setup.LoginStep,
		release.Step,
	)
}
