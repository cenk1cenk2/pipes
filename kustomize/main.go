// Command pipe-kustomize runs Kustomize operations for CI pipelines.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/kustomize/build"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

const name = "pipe-kustomize"

const description = "Kustomize operations for CI pipelines."

var VERSION = "latest"

// options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type options struct{}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, _ options) *ucli.Command {
	tool := setup.Step

	return cli.App(name, description, version,
		cli.Command(p, "build", "Build and validate Kustomize overlays.", tool, build.Step(build.Deps{Tool: setup.C})),
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
