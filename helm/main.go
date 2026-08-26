package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/helm/install"
	"gitlab.kilic.dev/devops/pipes/helm/lint"
	"gitlab.kilic.dev/devops/pipes/helm/login"
	"gitlab.kilic.dev/devops/pipes/helm/publish"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

func main() {
	NewPlumber(
		func(p *Plumber) *ucli.Command {
			tool := setup.Step
			credentials := login.Step

			return cli.App(CLI_NAME, DESCRIPTION, VERSION,
				cli.Command(p, "install", "Install Helm chart dependencies.", tool, credentials, install.Step(install.Deps{Tool: setup.C.Ctx})),
				cli.Command(p, "lint", "Lint Helm chart templates.", tool, lint.Step(lint.Deps{Tool: setup.C.Ctx})),
				cli.Command(p, "publish", "Publish Helm chart templates.", tool, credentials, publish.Step(publish.Deps{Tool: setup.C})),
			)
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
