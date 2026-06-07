package plan

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func TerraformPlan(tl *TaskList) *Task {
	return tl.CreateTask("plan").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"plan",
				"-input=false",
			).
				Set(func(c *Command) error {
					if P.Plan.Output != "" {
						c.AppendArgs(fmt.Sprintf("-out=%s", P.Plan.Output))
					}

					if P.Plan.Args != "" {
						c.AppendArgs(P.Plan.Args)
					}

					return nil
				}).
				SetDir(setup.P.Project.Cwd).
				AppendEnvironment(setup.C.EnvVars).
				SetRetries(&CommandRetry{
					Tries: P.Plan.RetryTries,
					Delay: P.Plan.RetryDelay,
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerraformMergeRequestReport(tl *TaskList) *Task {
	return tl.CreateTask("merge-request-report").
		ShouldDisable(func(t *Task) bool {
			return !P.MergeRequestReport.Enabled || P.MergeRequestReport.MergeRequestId == 0
		}).
		Set(func(t *Task) error {
			if P.Plan.Output == "" {
				return fmt.Errorf("terraform plan output is required for GitLab merge request report")
			}

			t.CreateCommand(
				"terraform",
				"show",
				"-json",
				P.Plan.Output,
			).
				SetDir(setup.P.Project.Cwd).
				AppendEnvironment(setup.C.EnvVars).
				SetLogLevel(LOG_LEVEL_TRACE, LOG_LEVEL_WARN, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					report, err := parseTerraformShowPlan([]byte(strings.Join(c.GetStdoutStream(), "")))
					if err != nil {
						return err
					}

					for _, item := range []mergeRequestReportMetadata{
						{Name: "Terraform working directory", Value: setup.P.Project.Cwd},
						{Name: "Terraform plan file", Value: P.Plan.Output},
						{Name: "GitLab project id", Value: P.MergeRequestReport.ProjectId},
						{Name: "GitLab merge request", Value: fmt.Sprintf("!%d", P.MergeRequestReport.MergeRequestId)},
						{Name: "Report identifier", Value: P.MergeRequestReport.Identifier},
					} {
						if item.Value != "" {
							report.Metadata = append(report.Metadata, item)
						}
					}

					body, err := renderMergeRequestReport(report)
					if err != nil {
						return err
					}

					result, err := gitlab.UpsertMergeRequestReport(
						context.Background(),
						P.MergeRequestReport,
						body,
					)
					if err != nil {
						return err
					}

					t.Log.Infof("Upserted GitLab merge request report note: %d", result.NoteId)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
