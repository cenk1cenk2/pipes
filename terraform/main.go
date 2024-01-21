package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/terraform/apply"
	"gitlab.kilic.dev/devops/pipes/terraform/install"
	"gitlab.kilic.dev/devops/pipes/terraform/lint"
	"gitlab.kilic.dev/devops/pipes/terraform/login"
	"gitlab.kilic.dev/devops/pipes/terraform/plan"
	"gitlab.kilic.dev/devops/pipes/terraform/publish"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	"gitlab.kilic.dev/devops/pipes/terraform/state"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Commands: cli.Commands{
					{
						Name:        "install",
						Description: "Install terraform project.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, state.Flags, install.Flags),
						Action: func(c *cli.Context) error {
							tl := &install.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
									state.New(p).SetCliContext(c).Job(),
									install.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:        "lint",
						Description: "Lint terraform project with terraform.",
						Flags:       p.AppendFlags(setup.Flags, lint.Flags),
						Action: func(c *cli.Context) error {
							tl := &lint.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									lint.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:        "plan",
						Description: "Plan terraform project.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, state.Flags, plan.Flags),
						Action: func(c *cli.Context) error {
							tl := &plan.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
									state.New(p).SetCliContext(c).Job(),
									plan.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:        "apply",
						Description: "Apply terraform project.",
						Flags:       p.AppendFlags(setup.Flags, login.Flags, state.Flags, apply.Flags),
						Action: func(c *cli.Context) error {
							tl := &apply.TL

							return tl.RunJobs(
								tl.JobSequence(
									setup.New(p).SetCliContext(c).Job(),
									login.New(p).SetCliContext(c).Job(),
									state.New(p).SetCliContext(c).Job(),
									apply.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},

					{
						Name:        "publish",
						Description: "Publish terraform project.",
						Flags:       p.AppendFlags(publish.Flags),
						Action: func(c *cli.Context) error {
							tl := &publish.TL

							return tl.RunJobs(
								tl.JobSequence(
									publish.New(p).SetCliContext(c).Job(),
								),
							)
						},
					},
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags:       true,
			ExcludeHelpCommand: true,
		}).
		Run()
}
