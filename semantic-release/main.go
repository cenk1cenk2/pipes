// Command pipe-semantic-release releases applications through the
// semantic-release library.
package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/semantic-release/release"
	"gitlab.kilic.dev/devops/pipes/semantic-release/setup"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		OverwriteCliFlag(setup.Flags, func(f *cli.BoolFlag) bool {
			return f.Name == "environment.enable"
		}, func(f *cli.BoolFlag) *cli.BoolFlag {
			f.Hidden = false
			f.Value = false

			return f
		})

		return &cli.Command{
			Name:        "pipe-semantic-release",
			Version:     version,
			Description: "Releases applications through the semantic-release library.",
			Flags:       CombineFlags(setup.Flags, release.Flags),
			Action: func(_ context.Context, _ *cli.Command) error {
				return p.RunJobs(CombineTaskLists(
					setup.New(p),
					release.New(p),
				))
			},
		}
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
