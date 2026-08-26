package ci

import (
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
)

//revive:disable:line-length-limit

// NewMetadataFlags reads the job and commit coordinates the CI runner exports, so
// a report can point back at the pipeline that produced it.
func NewMetadataFlags(dst *iac.Metadata) []ucli.Flag {
	return []ucli.Flag{
		&ucli.StringFlag{
			Category:    cli.CATEGORY_CI,
			Name:        "report.job-name",
			Sources:     cli.EnvVars("CI_JOB_NAME"),
			Usage:       "GitLab CI job name to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.JobName,
		},

		&ucli.StringFlag{
			Category:    cli.CATEGORY_CI,
			Name:        "report.job-url",
			Sources:     cli.EnvVars("CI_JOB_URL"),
			Usage:       "GitLab CI job URL to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.JobUrl,
		},

		&ucli.StringFlag{
			Category:    cli.CATEGORY_CI,
			Name:        "report.pipeline-id",
			Sources:     cli.EnvVars("CI_PIPELINE_ID"),
			Usage:       "GitLab CI pipeline id to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.PipelineId,
		},

		&ucli.StringFlag{
			Category:    cli.CATEGORY_CI,
			Name:        "report.pipeline-url",
			Sources:     cli.EnvVars("CI_PIPELINE_URL"),
			Usage:       "GitLab CI pipeline URL to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.PipelineUrl,
		},

		&ucli.StringFlag{
			Category:    cli.CATEGORY_CI,
			Name:        "report.commit-sha",
			Sources:     cli.EnvVars("CI_COMMIT_SHA"),
			Usage:       "Git commit sha to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.CommitSha,
		},

		&ucli.StringFlag{
			Category:    cli.CATEGORY_CI,
			Name:        "report.commit-short-sha",
			Sources:     cli.EnvVars("CI_COMMIT_SHORT_SHA"),
			Usage:       "Short git commit sha to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &dst.CommitShortSha,
		},
	}
}
