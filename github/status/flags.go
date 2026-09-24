package status

import (
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/github/client"
)

//revive:disable:line-length-limit

const (
	CategoryStatus = "Status"
)

var Flags = append(client.NewFlags(client.Options{Destination: &P.App}), []cli.Flag{
	&cli.StringFlag{
		Category: CategoryStatus,
		Name:     "status.token",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_STATUS_TOKEN"),
			cli.EnvVar("GH_TOKEN"),
		),
		Usage:       "GitHub token to post the status with. Not needed when the GitHub App flags are given, since the status then mints its own token.",
		Required:    false,
		Value:       "",
		Destination: &P.Status.Token,
	},

	&cli.StringFlag{
		Category: CategoryStatus,
		Name:     "status.project",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_STATUS_PROJECT"),
		),
		Usage:       "GitHub repository to post the status to, as owner/repository.",
		Required:    true,
		Value:       "",
		Destination: &P.Status.Project,
	},

	&cli.StringFlag{
		Category: CategoryStatus,
		Name:     "status.state",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_STATUS_STATE"),
			cli.EnvVar("GITHUB_STATUS_REPORT"),
		),
		Usage:       "State of the status. format(enum(pending, success, failure, error))",
		Required:    true,
		Value:       "",
		Destination: &P.Status.State,
	},

	&cli.StringFlag{
		Category: CategoryStatus,
		Name:     "status.sha",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_STATUS_SHA"),
			cli.EnvVar("CI_COMMIT_SHA"),
		),
		Usage:       "Commit sha to post the status for.",
		Required:    true,
		Value:       "",
		Destination: &P.Status.Sha,
	},

	&cli.StringFlag{
		Category: CategoryStatus,
		Name:     "status.target-url",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_STATUS_TARGET_URL"),
			cli.EnvVar("CI_PIPELINE_URL"),
		),
		Usage:       "URL the status links to.",
		Required:    false,
		Value:       "",
		Destination: &P.Status.TargetUrl,
	},

	&cli.StringFlag{
		Category: CategoryStatus,
		Name:     "status.context",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_STATUS_CONTEXT"),
		),
		Usage:       "Context that tells this status apart from the others on the commit.",
		Required:    false,
		Value:       "Gitlab CI",
		Destination: &P.Status.Context,
	},

	&cli.StringFlag{
		Category: CategoryStatus,
		Name:     "status.description",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_STATUS_DESCRIPTION"),
		),
		Usage:       "Short description of the status.",
		Required:    false,
		Value:       "",
		Destination: &P.Status.Description,
	},
}...)
