package plan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

const terraformSummaryOutput = "terraform-summary-report.json"

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

func TerraformReport(tl *TaskList) *Task {
	return tl.CreateTask("report").
		ShouldDisable(func(t *Task) bool {
			return !P.Report.Enabled
		}).
		Set(func(t *Task) error {
			if P.Plan.Output == "" {
				return fmt.Errorf("terraform plan output is required for report")
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
					plan, err := parseTerraformShowPlan([]byte(strings.Join(c.GetStdoutStream(), "")))
					if err != nil {
						return err
					}

					report := buildTerraformReport(plan)
					for _, item := range []reportMetadata{
						{Name: "Working directory", Value: setup.P.Project.Cwd},
						{Name: "Plan file", Value: P.Plan.Output},
					} {
						if item.Value != "" {
							report.Metadata = append(report.Metadata, item)
						}
					}

					body, err := renderReport(report)
					if err != nil {
						return err
					}

					reportPath := P.Report.Output
					if !filepath.IsAbs(reportPath) {
						reportPath = filepath.Join(setup.P.Project.Cwd, reportPath)
					}
					if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
						return fmt.Errorf("create Terraform report directory: %w", err)
					}
					if err := os.WriteFile(reportPath, body, 0o644); err != nil {
						return fmt.Errorf("write Terraform report %s: %w", reportPath, err)
					}
					t.Log.Infof("Wrote Terraform report artifact: %s", reportPath)

					summaryReport, err := renderSummary(summarizeTerraformPlan(plan))
					if err != nil {
						return err
					}

					summaryReportPath := terraformSummaryOutput
					if !filepath.IsAbs(summaryReportPath) {
						summaryReportPath = filepath.Join(setup.P.Project.Cwd, summaryReportPath)
					}
					if err := os.MkdirAll(filepath.Dir(summaryReportPath), 0o755); err != nil {
						return fmt.Errorf("create Terraform summary report directory: %w", err)
					}
					if err := os.WriteFile(summaryReportPath, summaryReport, 0o644); err != nil {
						return fmt.Errorf("write Terraform summary report %s: %w", summaryReportPath, err)
					}
					t.Log.Infof("Wrote Terraform summary report artifact: %s", summaryReportPath)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
