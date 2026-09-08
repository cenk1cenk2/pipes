package gitlab

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryGitLabMergeRequestReport = "GitLab Merge Request Report"
)

// Options is what the merge request report flags are built onto.
type Options struct {
	Destination *MergeRequestReportConfig
}

func NewFlags(opts Options) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Category: CategoryGitLabMergeRequestReport,
			Name:     "gitlab-mr-report.enabled",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITLAB_MR_REPORT_ENABLED"),
			),
			Usage:       "Enable GitLab merge request report note on the given merge request.",
			Required:    false,
			Value:       false,
			Destination: &opts.Destination.Enabled,
		},

		&cli.StringFlag{
			Category: CategoryGitLabMergeRequestReport,
			Name:     "gitlab-mr-report.token",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GL_PIPES_TOKEN"),
			),
			Usage:       "GitLab API token for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.Token,
		},

		&cli.StringFlag{
			Category: CategoryGitLabMergeRequestReport,
			Name:     "gitlab-mr-report.api-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_API_V4_URL"),
			),
			Usage:       "GitLab API URL for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.ApiUrl,
		},

		&cli.StringFlag{
			Category: CategoryGitLabMergeRequestReport,
			Name:     "gitlab-mr-report.project-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PROJECT_ID"),
			),
			Usage:       "GitLab project id for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.ProjectId,
		},

		&cli.Int64Flag{
			Category: CategoryGitLabMergeRequestReport,
			Name:     "gitlab-mr-report.merge-request-iid",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_MERGE_REQUEST_IID"),
			),
			Usage:       "GitLab merge request iid for merge request report notes.",
			Required:    false,
			Value:       0,
			Destination: &opts.Destination.MergeRequestIid,
		},

		&cli.StringFlag{
			Category: CategoryGitLabMergeRequestReport,
			Name:     "gitlab-mr-report.identifier",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITLAB_MR_REPORT_IDENTIFIER"),
			),
			Usage:       "Hidden marker identifier for merge request report notes. Defaults to the job name combined with the stack or state under report.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.Identifier,
		},
	}
}
