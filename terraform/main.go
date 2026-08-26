package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
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

func main() {
	NewPlumber(
		func(p *Plumber) *ucli.Command {
			tool := setup.Step
			credentials := login.Step(login.Deps{Tool: setup.C})
			backend := state.Step(state.Deps{Tool: setup.C, CI: &setup.P.CiVariables})

			return cli.App(CLI_NAME, DESCRIPTION, VERSION,
				cli.Command(p, "install", "Install terraform project.", tool, credentials, backend, install.Step(install.Deps{Tool: setup.C})),
				cli.Command(p, "lint", "Lint terraform project with terraform.", tool, lint.Step(lint.Deps{Tool: setup.C})),
				cli.Command(p, "plan", "Plan terraform project.", tool, credentials, backend, plan.Step(plan.Deps{Tool: setup.C, State: state.P, Notes: gitlab.NewNotes})),
				cli.Command(p, "apply", "Apply terraform project.", tool, credentials, backend, apply.Step(apply.Deps{Tool: setup.C})),
				cli.Command(p, "publish", "Publish terraform project.", publish.Step(publish.Deps{Registry: gitlab.NewModuleRegistry})),
			)
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
