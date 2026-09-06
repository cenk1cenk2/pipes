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

// options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type options struct{}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, _ options) *ucli.Command {
	tool := setup.Step
	credentials := login.Step

	return cli.App(name, description, version,
		cli.Command(p, "install", "Install Helm chart dependencies.", tool, credentials, install.Step(install.Deps{Tool: setup.C.Ctx})),
		cli.Command(p, "lint", "Lint Helm chart templates.", tool, lint.Step(lint.Deps{Tool: setup.C.Ctx})),
		cli.Command(p, "publish", "Publish Helm chart templates.", tool, credentials, publish.Step(publish.Deps{Tool: setup.C})),
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
