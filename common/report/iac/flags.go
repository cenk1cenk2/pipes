package iac

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

func NewMetadataFlags(metadata *Metadata) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "report.job-name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_JOB_NAME"),
			),
			Usage:       "GitLab CI job name to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &metadata.JobName,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "report.job-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_JOB_URL"),
			),
			Usage:       "GitLab CI job URL to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &metadata.JobUrl,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "report.pipeline-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PIPELINE_ID"),
			),
			Usage:       "GitLab CI pipeline id to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &metadata.PipelineId,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "report.pipeline-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PIPELINE_URL"),
			),
			Usage:       "GitLab CI pipeline URL to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &metadata.PipelineUrl,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "report.commit-sha",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_COMMIT_SHA"),
			),
			Usage:       "Git commit sha to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &metadata.CommitSha,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "report.commit-short-sha",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_COMMIT_SHORT_SHA"),
			),
			Usage:       "Short git commit sha to include in the plan report metadata.",
			Required:    false,
			Value:       "",
			Destination: &metadata.CommitShortSha,
		},
	}
}
