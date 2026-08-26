package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/kustomize/build"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

func main() {
	NewPlumber(
		func(p *Plumber) *ucli.Command {
			tool := setup.Step

			return cli.App(CLI_NAME, DESCRIPTION, VERSION,
				cli.Command(p, "build", "Build and validate Kustomize overlays.", tool, build.Step(build.Deps{Tool: setup.C})),
			)
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
