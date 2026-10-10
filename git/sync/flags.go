package sync

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

const (
	CategorySync     = "Sync"
	CategoryIdentity = "Identity"
	CategoryProject  = "GitLab Project"
)

//revive:disable:line-length-limit

var Flags = CombineFlags(
	[]cli.Flag{
		&cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.mode",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_MODE"),
			),
			Usage:       "Whether the changes are only reported or committed and opened as a merge request on the target as well. format(enum(publish, report))",
			Required:    false,
			Value:       ModeReport,
			Destination: &P.Mode,
		},

		&cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.target.project",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_TARGET_PROJECT"),
				cli.EnvVar("CI_PROJECT_PATH"),
			),
			Usage:       "Path of the GitLab project the paths are synced into. Defaults to the project of the pipeline.",
			Required:    false,
			Value:       "",
			Destination: &P.Target.Project,
		},

		&cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.target.branch",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_TARGET_BRANCH"),
			),
			Usage:       "Branch of the target the paths are synced onto, and the merge request targets.",
			Required:    false,
			Value:       "next",
			Destination: &P.Target.Branch,
		},

		flags.YAMLFlag(&P.Paths, &cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.paths",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_PATHS"),
			),
			Usage:    "Sources of the job synced onto destinations in the target, which are replaced completely. Include and exclude are globs anchored at the source. format(yaml([]{ source: string, destination: string, include?: []string, exclude?: []string }))",
			Required: true,
		}),

		&cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.branch",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_BRANCH"),
			),
			Usage:       "Branch the changes are pushed to and the merge request is opened from. Defaults to sync/$CI_PROJECT_PATH_SLUG-$CI_JOB_NAME_SLUG.",
			Required:    false,
			Value:       "",
			Destination: &P.Branch,
		},

		&cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.commit-message",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_COMMIT_MESSAGE"),
			),
			Usage:       "Message of the commit, whose subject is the title of the merge request as well.",
			Required:    false,
			Value:       "chore: sync generated files",
			Destination: &P.CommitMessage,
		},

		&cli.BoolFlag{
			Category: CategorySync,
			Name:     "git.sync.allow-empty",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_ALLOW_EMPTY"),
			),
			Usage:       "Sync a missing or empty source as an empty destination instead of failing.",
			Required:    false,
			Value:       false,
			Destination: &P.AllowEmpty,
		},

		&cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.patch",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_PATCH"),
			),
			Usage:       "File the full staged diff is written to, for the job artifacts. A relative path resolves against the working directory.",
			Required:    false,
			Value:       "git-sync.patch",
			Destination: &P.Patch,
		},

		&cli.StringFlag{
			Category: CategorySync,
			Name:     "git.sync.token",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_SYNC_GL_TOKEN"),
			),
			Usage:       "GitLab token that fetches and pushes the target and opens the merge request on it. A publish also reads the head of the default branch of the project of the pipeline with it, so a target in another project needs read access to the project of the pipeline as well.",
			Required:    false,
			Value:       "",
			Destination: &P.Token,
		},

		&cli.StringFlag{
			Category: CategoryIdentity,
			Name:     "git.pipes.author.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_PIPES_AUTHOR_NAME"),
			),
			Usage:       "Name of the author of the commit. Required in publish mode.",
			Required:    false,
			Value:       "",
			Destination: &P.Author.Name,
		},

		&cli.StringFlag{
			Category: CategoryIdentity,
			Name:     "git.pipes.author.email",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_PIPES_AUTHOR_EMAIL"),
			),
			Usage:       "Email of the author of the commit. Required in publish mode.",
			Required:    false,
			Value:       "",
			Destination: &P.Author.Email,
		},

		&cli.StringFlag{
			Category: CategoryIdentity,
			Name:     "git.pipes.committer.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_PIPES_COMMITTER_NAME"),
			),
			Usage:       "Name of the committer of the commit. Required in publish mode.",
			Required:    false,
			Value:       "",
			Destination: &P.Committer.Name,
		},

		&cli.StringFlag{
			Category: CategoryIdentity,
			Name:     "git.pipes.committer.email",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GIT_PIPES_COMMITTER_EMAIL"),
			),
			Usage:       "Email of the committer of the commit. Required in publish mode.",
			Required:    false,
			Value:       "",
			Destination: &P.Committer.Email,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.commit-author",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_COMMIT_AUTHOR"),
			),
			Usage:       "Author of the commit of the pipeline, credited as a co-author of the sync commit.",
			Required:    false,
			Value:       "",
			Destination: &P.Project.CommitAuthor,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.path",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PROJECT_PATH"),
			),
			Usage:       "Path of the GitLab project the pipeline runs for.",
			Required:    false,
			Value:       "",
			Destination: &P.Project.Path,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.path-slug",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PROJECT_PATH_SLUG"),
			),
			Usage:       "Slug of the path of the GitLab project the pipeline runs for.",
			Required:    false,
			Value:       "",
			Destination: &P.Project.PathSlug,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.job-slug",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_JOB_NAME_SLUG"),
			),
			Usage:       "Slug of the name of the job.",
			Required:    false,
			Value:       "",
			Destination: &P.Project.JobSlug,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.dir",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PROJECT_DIR"),
			),
			Usage:       "Checkout of the project the target branch is checked out next to when the project is the target.",
			Required:    false,
			Value:       ".",
			Destination: &P.Project.Dir,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.server-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_SERVER_URL"),
			),
			Usage:       "URL of the GitLab instance the projects are fetched from and pushed to.",
			Required:    false,
			Value:       "",
			Destination: &P.Project.ServerUrl,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.default-branch",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_DEFAULT_BRANCH"),
			),
			Usage:       "Default branch of the project, which the commit of the pipeline has to still be the head of for a publish.",
			Required:    false,
			Value:       "",
			Destination: &P.Project.DefaultBranch,
		},

		&cli.StringFlag{
			Category: CategoryProject,
			Name:     "git.sync.project.open-merge-requests",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_OPEN_MERGE_REQUESTS"),
			),
			Usage:       "Merge requests open from the branch of the pipeline, which the report is posted on outside of a merge request pipeline.",
			Required:    false,
			Value:       "",
			Destination: &P.Project.OpenMergeRequests,
		},
	},
	git.NewFlags(git.Options{Destination: &P.Refs}),
	gitlab.NewFlags(gitlab.Options{Destination: &P.MergeRequestReport}),
	ci.NewFlags(ci.Options{Destination: &P.ReportMetadata}),
)
