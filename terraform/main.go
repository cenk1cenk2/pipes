// Command pipe-terraform runs terraform inside CI pipelines.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/terraform/apply"
	"gitlab.kilic.dev/devops/pipes/terraform/install"
	"gitlab.kilic.dev/devops/pipes/terraform/lint"
	"gitlab.kilic.dev/devops/pipes/terraform/login"
	"gitlab.kilic.dev/devops/pipes/terraform/plan"
	"gitlab.kilic.dev/devops/pipes/terraform/publish"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	"gitlab.kilic.dev/devops/pipes/terraform/state"
)

const name = "pipe-terraform"

const description = "Running terraform inside the pipelines."

var VERSION = "latest"

// options are the services the pipe reaches outside the machine for. A zero
// value is the production wiring, so only a spec ever fills one in.
type options struct {
	Notes    gitlab.NotesFactory
	Registry func(apiUrl, projectId, token string) gitlab.ModuleRegistry
}

func (o options) defaults() options {
	if o.Notes == nil {
		o.Notes = gitlab.NewNotes
	}

	if o.Registry == nil {
		o.Registry = gitlab.NewModuleRegistry
	}

	return o
}

// newCommand builds the command tree. The version is a parameter so a spec
// can build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, opts options) *ucli.Command {
	opts = opts.defaults()

	tool := setup.Step
	credentials := login.Step(login.Deps{Tool: setup.C})
	backend := state.Step(state.Deps{Tool: setup.C, CI: &setup.P.CiVariables})

	return cli.App(name, description, version,
		cli.Command(p, "install", "Install terraform project.", tool, credentials, backend, install.Step(install.Deps{Tool: setup.C})),
		cli.Command(p, "lint", "Lint terraform project with terraform.", tool, lint.Step(lint.Deps{Tool: setup.C})),
		cli.Command(p, "plan", "Plan terraform project.", tool, credentials, backend, plan.Step(plan.Deps{Tool: setup.C, State: state.P, Notes: opts.Notes})),
		cli.Command(p, "apply", "Apply terraform project.", tool, credentials, backend, apply.Step(apply.Deps{Tool: setup.C})),
		cli.Command(p, "publish", "Publish terraform project.", publish.Step(publish.Deps{Registry: opts.Registry})),
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
