package preview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
	"gitlab.kilic.dev/devops/pipes/common/report/iac"
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

func PulumiSummary(tl *TaskList) *Task {
	return tl.CreateTask("summary").
		ShouldDisable(func(t *Task) bool {
			if P.Summary.Output == "" {
				t.Log.Debugln("Skipping Pulumi summary because no summary output file is configured.")

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

			report, err := parsePulumiPlanReport(data, iac.Metadata{})
			if err != nil {
				return err
			}

			body, err := renderSummary(summarizePulumiReport(report))
			if err != nil {
				return err
			}

			output := P.Summary.Output
			if !filepath.IsAbs(output) {
				output = filepath.Join(setup.P.Cwd, output)
			}

			if err := os.WriteFile(output, body, 0o644); err != nil {
				return fmt.Errorf("write Pulumi summary %s: %w", output, err)
			}

			t.Log.Infof("Wrote Pulumi summary: %s", output)

			return nil
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
			metadata.Target = stack.P.Stack
			metadata.Cwd = setup.P.Cwd

			report, err := parsePulumiPlanReport(data, metadata)
			if err != nil {
				return err
			}

			body, err := iac.RenderMergeRequestReport(report)
			if err != nil {
				return err
			}

			config := P.MergeRequestReportConfig
			config.Identifier = gitlab.ResolveReportIdentifier(
				config.Identifier,
				metadata.JobName,
				stack.P.Stack,
			)
			config.LegacyIdentifiers = []string{metadata.JobName}

			result, err := gitlab.UpsertMergeRequestReport(
				context.Background(),
				config,
				body,
			)
			if err != nil {
				return err
			}

			t.Log.Infof("Upserted GitLab merge request report note: %d", result.NoteId)

			return nil
		})
}
