// Command pipe-helm runs Helm operations for CI pipelines.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/helm/install"
	"gitlab.kilic.dev/devops/pipes/helm/lint"
	"gitlab.kilic.dev/devops/pipes/helm/login"
	"gitlab.kilic.dev/devops/pipes/helm/publish"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

const name = "pipe-helm"

const description = "Helm charts for CI pipelines."

var VERSION = "latest"

func newCommand(p *plumber.Plumber) *ucli.Command {
	tool := setup.Step
	credentials := login.Step

	return cli.App(name, description, VERSION,
		cli.Command(p, "install", "Install Helm chart dependencies.", tool, credentials, install.Step(install.Deps{Tool: setup.C.Ctx})),
		cli.Command(p, "lint", "Lint Helm chart templates.", tool, lint.Step(lint.Deps{Tool: setup.C.Ctx})),
		cli.Command(p, "publish", "Publish Helm chart templates.", tool, credentials, publish.Step(publish.Deps{Tool: setup.C})),
	)
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
