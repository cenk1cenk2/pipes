// Command pipe-go builds Go applications from CI/CD.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/build"
	"gitlab.kilic.dev/devops/pipes/go/install"
	"gitlab.kilic.dev/devops/pipes/go/lint"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	gotool "gitlab.kilic.dev/devops/pipes/go/tool"
)

func main() {
	plumber.NewPlumber(func(p *plumber.Plumber) *cli.Command {
		return &cli.Command{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Commands: []*cli.Command{
				{
					Name:        "install",
					Description: "Vendor go modules.",
					Flags:       plumber.CombineFlags(setup.Flags, install.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(plumber.CombineTaskLists(
							setup.New(p),
							install.New(p),
						))
					},
				},
				{
					Name:        "build",
					Description: "Build an application.",
					Flags:       plumber.CombineFlags(setup.Flags, build.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(plumber.CombineTaskLists(
							setup.New(p),
							build.New(p),
						))
					},
				},
				{
					Name:        "lint",
					Description: "Run golangci-lint on the project.",
					Flags:       plumber.CombineFlags(setup.Flags, lint.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(plumber.CombineTaskLists(
							setup.New(p),
							lint.New(p),
						))
					},
				},
				{
					Name:        "tool",
					Description: "Run a specified go tool.",
					Flags:       plumber.CombineFlags(setup.Flags, gotool.Flags),
					Arguments:   gotool.Arguments,
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(plumber.CombineTaskLists(
							setup.New(p),
							gotool.New(p),
						))
					},
				},
			},
		}
	}).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
