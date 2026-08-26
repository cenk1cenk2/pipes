package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/pulumi/preview"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
	"gitlab.kilic.dev/devops/pipes/pulumi/up"
)

func main() {
	NewPlumber(
		func(p *Plumber) *ucli.Command {
			tool := setup.Step
			selected := stack.Step(stack.Deps{Tool: setup.C})

			return cli.App(CLI_NAME, DESCRIPTION, VERSION,
				cli.Command(p, "preview", "Preview the Pulumi changes.", tool, selected, preview.Step(preview.Deps{Tool: setup.C, Stack: stack.P, Notes: gitlab.NewNotes})),
				cli.Command(p, "up", "Apply the Pulumi changes.", tool, selected, up.Step(up.Deps{Tool: setup.C})),
			)
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
