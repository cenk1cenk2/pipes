// Command pipe-kustomize runs Kustomize operations for CI pipelines.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/kustomize/build"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

func newCommand(p *plumber.Plumber) *ucli.Command {
	tool := setup.Step

	return cli.App(CLI_NAME, DESCRIPTION, VERSION,
		cli.Command(p, "build", "Build and validate Kustomize overlays.", tool, build.Step(build.Deps{Tool: setup.C})),
	)
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
