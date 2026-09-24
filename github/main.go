// Command pipe-github authenticates as a GitHub App for CI pipelines.
package main

import (
	"context"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/github/status"
	"gitlab.kilic.dev/devops/pipes/github/token"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        "pipe-github",
			Version:     version,
			Description: "GitHub App tokens and commit statuses in the pipeline.",
			Commands: []*cli.Command{
				{
					Name: "token",
					Description: strings.TrimSpace(`
Mint a GitHub App installation token into a dotenv file. Every GitHub App flag is required here.

The token lives for an hour. Expose the file as an artifacts:reports:dotenv report, and every job that needs this one receives the token as a variable, which overrides the job-level variable of the same name.
`),
					Flags: CombineFlags(token.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							token.New(p),
						))
					},
				},
				{
					Name: "status",
					Description: strings.TrimSpace(`
Post a commit status to GitHub.

Given the GitHub App flags, the status mints its own token, narrowed to the repository and to writing statuses. That needs the private key in the status job as well, so prefer the token from the dotenv file and mint in place only where the pipeline can outlive the token.

A commit GitHub does not have, as on a branch that only exists on GitLab, is skipped with a warning; every other failure fails the job.
`),
					Flags: CombineFlags(status.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							status.New(p),
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
