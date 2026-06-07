package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
)

//revive:disable:line-length-limit

var Flags = CombineFlags(
	[]cli.Flag{
		&cli.StringFlag{
			Name: "pulumi.preview.plan",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_PLAN"),
			),
			Usage:       "Output file for pulumi plan.",
			Required:    false,
			Value:       "plan.json",
			Destination: &P.Plan,
		},

		&cli.StringFlag{
			Name: "pulumi.preview.summary-output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PULUMI_SUMMARY_OUTPUT"),
			),
			Usage:       "Output file for Pulumi preview summary. Leave empty to skip summary generation.",
			Required:    false,
			Value:       "pulumi-summary.json",
			Destination: &P.Summary.Output,
		},
	},
	gitlab.NewMergeRequestReportFlags(&P.MergeRequestReportConfig),
	[]cli.Flag{
		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "gitlab-mr-report.job-name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_JOB_NAME"),
			),
			Usage:       "GitLab CI job name to include in merge request report metadata.",
			Required:    false,
			Value:       "",
			Destination: &P.ReportMetadata.JobName,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "gitlab-mr-report.job-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_JOB_URL"),
			),
			Usage:       "GitLab CI job URL to include in merge request report metadata.",
			Required:    false,
			Value:       "",
			Destination: &P.ReportMetadata.JobUrl,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "gitlab-mr-report.pipeline-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PIPELINE_ID"),
			),
			Usage:       "GitLab CI pipeline id to include in merge request report metadata.",
			Required:    false,
			Value:       "",
			Destination: &P.ReportMetadata.PipelineId,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "gitlab-mr-report.pipeline-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PIPELINE_URL"),
			),
			Usage:       "GitLab CI pipeline URL to include in merge request report metadata.",
			Required:    false,
			Value:       "",
			Destination: &P.ReportMetadata.PipelineUrl,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "gitlab-mr-report.commit-sha",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_COMMIT_SHA"),
			),
			Usage:       "Git commit sha to include in merge request report metadata.",
			Required:    false,
			Value:       "",
			Destination: &P.ReportMetadata.CommitSha,
		},

		&cli.StringFlag{
			Category: flags.CATEGORY_GITLAB_PIPELINE,
			Name:     "gitlab-mr-report.commit-short-sha",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_COMMIT_SHORT_SHA"),
			),
			Usage:       "Short git commit sha to include in merge request report metadata.",
			Required:    false,
			Value:       "",
			Destination: &P.ReportMetadata.CommitShortSha,
		},
	},
)
