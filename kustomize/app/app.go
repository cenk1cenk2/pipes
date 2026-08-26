// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/kustomize/build"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

const Name = "pipe-kustomize"

const Description = "Kustomize operations for CI pipelines."

// Options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type Options struct{}

func New(p *plumber.Plumber, version string, _ Options) *ucli.Command {
	tool := setup.Step

	return cli.App(Name, Description, version,
		cli.Command(p, "build", "Build and validate Kustomize overlays.", tool, build.Step(build.Deps{Tool: setup.C})),
	)
}
