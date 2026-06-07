package preview

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
)

const pulumiSummaryOutput = "pulumi-summary-report.json"

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

func ReportTask(tl *TaskList) *Task {
	return tl.CreateTask("report").
		ShouldDisable(func(t *Task) bool {
			return !P.Report.Enabled
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

			report, err := parsePulumiReport(data, reportMetadata{
				Stack: stack.P.Stack,
			})
			if err != nil {
				return err
			}

			body, err := renderReport(report)
			if err != nil {
				return err
			}

			reportPath := P.Report.Output
			if !filepath.IsAbs(reportPath) {
				reportPath = filepath.Join(setup.P.Cwd, reportPath)
			}
			if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
				return fmt.Errorf("create Pulumi report directory: %w", err)
			}
			if err := os.WriteFile(reportPath, body, 0o644); err != nil {
				return fmt.Errorf("write Pulumi report %s: %w", reportPath, err)
			}

			t.Log.Infof("Wrote Pulumi report artifact: %s", reportPath)

			summaryReport, err := renderSummary(summarizePulumiReport(report))
			if err != nil {
				return err
			}

			summaryReportPath := pulumiSummaryOutput
			if !filepath.IsAbs(summaryReportPath) {
				summaryReportPath = filepath.Join(setup.P.Cwd, summaryReportPath)
			}
			if err := os.MkdirAll(filepath.Dir(summaryReportPath), 0o755); err != nil {
				return fmt.Errorf("create Pulumi summary report directory: %w", err)
			}
			if err := os.WriteFile(summaryReportPath, summaryReport, 0o644); err != nil {
				return fmt.Errorf("write Pulumi summary report %s: %w", summaryReportPath, err)
			}

			t.Log.Infof("Wrote Pulumi summary report artifact: %s", summaryReportPath)

			return nil
		})
}
