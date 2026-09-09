package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/build"
	"gitlab.kilic.dev/devops/pipes/go/install"
	"gitlab.kilic.dev/devops/pipes/go/lint"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	gotool "gitlab.kilic.dev/devops/pipes/go/tool"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        "pipe-go",
			Version:     version,
			Description: "Pipe for Go builds.",
			Commands: []*cli.Command{
				{
					Name:        "install",
					Description: "Vendor go modules.",
					Flags:       CombineFlags(setup.Flags, install.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							install.New(p),
						))
					},
				},
				{
					Name:        "build",
					Description: "Build an application.",
					Flags:       CombineFlags(setup.Flags, build.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							build.New(p),
						))
					},
				},
				{
					Name:        "lint",
					Description: "Run golangci-lint on the project.",
					Flags:       CombineFlags(setup.Flags, lint.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							lint.New(p),
						))
					},
				},
				{
					Name:        "tool",
					Description: "Run a specified go tool.",
					Flags:       CombineFlags(setup.Flags, gotool.Flags),
					Arguments:   gotool.Arguments,
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							gotool.New(p),
						))
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
