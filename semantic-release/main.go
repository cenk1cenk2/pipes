// Command pipe-semantic-release releases applications through the
// semantic-release library.
package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/semantic-release/release"
	"gitlab.kilic.dev/devops/pipes/semantic-release/setup"
)

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		// the environment is opt-in here, and the flags are shared package level values,
		// so this runs before the command tree reads them.
		OverwriteCliFlag(setup.EnvironmentFlags, func(f *cli.BoolFlag) bool {
			return f.Name == "environment.enable"
		}, func(f *cli.BoolFlag) *cli.BoolFlag {
			f.Hidden = false
			f.Value = false

			return f
		})

		return &cli.Command{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Flags:       CombineFlags(setup.EnvironmentFlags, setup.NodeFlags, setup.LoginFlags, release.Flags),
			Action: func(_ context.Context, _ *cli.Command) error {
				return p.RunJobs(CombineTaskLists(
					setup.NewEnvironment(p),
					setup.New(p),
					setup.NewLogin(p),
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
