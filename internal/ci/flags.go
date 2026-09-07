package ci

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CI = "Gitlab Pipeline"
)

// NewFlags reads the job and commit coordinates the CI runner exports, so
// a report can point back at the pipeline that produced it.
func NewFlags(dst *iac.Metadata) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category:    CATEGORY_CI,
			Name:        "ci.job-name",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_JOB_NAME")),
			Usage:       "GitLab CI job name to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.JobName,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CI,
			Name:        "ci.job-url",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_JOB_URL")),
			Usage:       "GitLab CI job URL to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.JobUrl,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CI,
			Name:        "ci.pipeline-id",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_PIPELINE_ID")),
			Usage:       "GitLab CI pipeline id to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.PipelineId,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CI,
			Name:        "ci.pipeline-url",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_PIPELINE_URL")),
			Usage:       "GitLab CI pipeline URL to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.PipelineUrl,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CI,
			Name:        "ci.commit-sha",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_COMMIT_SHA")),
			Usage:       "Git commit sha to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.CommitSha,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CI,
			Name:        "ci.commit-short-sha",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_COMMIT_SHORT_SHA")),
			Usage:       "Short git commit sha to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.CommitShortSha,
		},
	}
}
