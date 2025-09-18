package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
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
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Commands: []*cli.Command{
					{
						Name:        "install",
						Description: "Install terraform project.",
						Flags:       CombineFlags(setup.Flags, login.Flags, state.Flags, install.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									login.New(p),
									state.New(p),
									install.New(p),
								),
							)
						},
					},

					{
						Name:        "lint",
						Description: "Lint terraform project with terraform.",
						Flags:       CombineFlags(setup.Flags, lint.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									lint.New(p),
								),
							)
						},
					},

					{
						Name:        "plan",
						Description: "Plan terraform project.",
						Flags:       CombineFlags(setup.Flags, login.Flags, state.Flags, plan.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									login.New(p),
									state.New(p),
									plan.New(p),
								),
							)
						},
					},

					{
						Name:        "apply",
						Description: "Apply terraform project.",
						Flags:       CombineFlags(setup.Flags, login.Flags, state.Flags, apply.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									setup.New(p),
									login.New(p),
									state.New(p),
									apply.New(p),
								),
							)
						},
					},

					{
						Name:        "publish",
						Description: "Publish terraform project.",
						Flags:       CombineFlags(publish.Flags),
						Action: func(_ context.Context, _ *cli.Command) error {
							return p.RunJobs(
								CombineTaskLists(
									publish.New(p),
								),
							)
						},
					},
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
