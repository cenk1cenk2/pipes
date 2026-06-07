package preview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
			if !P.MergeRequestReportConfig.Enabled {
				return true
			}

			if P.MergeRequestReportConfig.MergeRequestId == 0 {
				t.Log.Debugln("Skipping GitLab merge request report because this is not a merge request pipeline.")

				return true
			}

			return false
		}).
		Set(func(t *Task) error {
			planPath := P.Plan
			if !filepath.IsAbs(planPath) {
				planPath = filepath.Join(setup.P.Cwd, planPath)
			}

			data, err := os.ReadFile(planPath)
			if err != nil {
				return fmt.Errorf("read Pulumi plan file %s: %w", planPath, err)
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
