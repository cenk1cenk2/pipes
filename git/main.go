// Command pipe-git writes what a job generated back into a git repository for CI pipelines.
package main

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/git/sync"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:        "pipe-git",
			Version:     version,
			Description: "Git related tasks in the pipeline.",
			Commands: []*cli.Command{
				{
					Name:        "sync",
					Description: "Sync the paths a job generated into a branch of a GitLab project. Replaces every destination on the target branch with its source, writes the staged diff to a patch, logs it and reports it on the merge requests of the pipeline, or in the description of the merge request of the sync branch when the pipeline runs for it. In publish mode, the changes are committed onto the tip of the target branch, force pushed to the sync branch and opened as a merge request against the target branch, whose title and description every later publish replaces. Once there is nothing to commit, the open merge request of the sync branch is commented on and closed, and the sync branch is deleted. Every merge request merged into the commit of a publish is then commented on with the merge request of the sync, or that there are no changes. A pipeline running on the target branch itself syncs like any other, and stops at nothing to commit once the merged output regenerates unchanged. A publish fails once the commit of the pipeline is no longer the head of the default branch.",
					Flags:       sync.Flags,
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							sync.New(p),
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
