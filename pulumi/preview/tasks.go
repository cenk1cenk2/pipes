package preview

import (
	"context"
	"fmt"
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
)

func PulumiPlan(tl *TaskList) *Task {
	return tl.CreateTask("plan").
		Set(func(t *Task) error {
			t.CreateCommand(
				"pulumi",
				"preview",
				"--non-interactive",
				"--diff",
				"--save-plan",
				P.Plan,
			).
				SetDir(setup.P.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func MergeRequestReport(tl *TaskList) *Task {
	return tl.CreateTask("gitlab", "merge-request-report").
		ShouldDisable(func(t *Task) bool {
			return !P.MergeRequestReportConfig.Enabled || P.MergeRequestReportConfig.MergeRequestId == 0
		}).
		Set(func(t *Task) error {
			data, err := os.ReadFile(P.Plan)
			if err != nil {
				return fmt.Errorf("read Pulumi plan file %s: %w", P.Plan, err)
			}

			metadata := P.ReportMetadata
			metadata.Stack = stack.P.Stack

			report, err := parsePulumiPlanReport(data, metadata)
			if err != nil {
				return err
			}

			body, err := renderMergeRequestReport(report)
			if err != nil {
				return err
			}

			result, err := gitlab.UpsertMergeRequestReport(
				context.Background(),
				P.MergeRequestReportConfig,
				body,
			)
			if err != nil {
				return err
			}

			t.Log.Infof("Upserted GitLab merge request report note: %d", result.NoteId)

			return nil
		})
}
