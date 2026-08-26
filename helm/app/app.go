// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

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

const Name = "pipe-helm"

const Description = "Helm charts for CI pipelines."

// Options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type Options struct{}

func New(p *plumber.Plumber, version string, _ Options) *ucli.Command {
	tool := setup.Step
	credentials := login.Step

	return cli.App(Name, Description, version,
		cli.Command(p, "install", "Install Helm chart dependencies.", tool, credentials, install.Step(install.Deps{Tool: setup.C.Ctx})),
		cli.Command(p, "lint", "Lint Helm chart templates.", tool, lint.Step(lint.Deps{Tool: setup.C.Ctx})),
		cli.Command(p, "publish", "Publish Helm chart templates.", tool, credentials, publish.Step(publish.Deps{Tool: setup.C})),
	)
}
