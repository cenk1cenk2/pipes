package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/python/build"
	"gitlab.kilic.dev/devops/pipes/python/install"
	"gitlab.kilic.dev/devops/pipes/python/lint"
	"gitlab.kilic.dev/devops/pipes/python/run"
	"gitlab.kilic.dev/devops/pipes/python/setup"
	"gitlab.kilic.dev/devops/pipes/python/test"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        "pipe-python",
			Version:     version,
			Description: "Pipe for installing, checking and building Python projects with uv on CI/CD.",
			Commands: []*cli.Command{
				{
					Name:        "install",
					Description: "Install the locked dependencies of the project.",
					Flags:       CombineFlags(setup.Flags, install.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							install.New(p),
						))
					},
				},
				{
					Name:        "lint",
					Description: "Run the lint commands inside the installed environment.",
					Flags:       CombineFlags(setup.Flags, lint.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							lint.New(p),
						))
					},
				},
				{
					Name:        "test",
					Description: "Run the test commands inside the installed environment.",
					Flags:       CombineFlags(setup.Flags, test.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							test.New(p),
						))
					},
				},
				{
					Name:        "build",
					Description: "Build the source distribution and the wheel of the project.",
					Flags:       CombineFlags(setup.Flags, build.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							build.New(p),
						))
					},
				},
				{
					Name:        "run",
					Description: "Run a command inside the installed environment.",
					Flags:       setup.Flags,
					Arguments:   run.Arguments,
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							run.New(p),
						))
					},
				},

				DocsCommand(p),
			},
		}
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
