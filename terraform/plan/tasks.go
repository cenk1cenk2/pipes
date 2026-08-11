package plan

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
	"gitlab.kilic.dev/devops/pipes/common/report/iac"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	"gitlab.kilic.dev/devops/pipes/terraform/state"
)

// The state name only means something once it has been set away from its default,
// which most backends never do.
func terraformStateName() string {
	if name := state.P.State.Name; name != "default" {
		return name
	}

	return ""
}

// Only the values that actually vary between concurrent plan jobs on one merge
// request belong in the marker, since anything else changes the identifier for every
// consumer without disambiguating anything.
func terraformReportDiscriminators() []string {
	discriminators := []string{}

	if name := terraformStateName(); name != "" {
		discriminators = append(discriminators, name)
	}

	if cwd := setup.P.Project.Cwd; cwd != "" && cwd != "." {
		discriminators = append(discriminators, cwd)
	}

	return discriminators
}

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

					if P.Plan.PreviewForMergeRequests && P.Plan.PipelineSource == "merge_request_event" {
						c.AppendArgs("-lock=false")
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

func TerraformSummary(tl *TaskList) *Task {
	return tl.CreateTask("summary").
		ShouldDisable(func(t *Task) bool {
			if P.Summary.Output == "" {
				t.Log.Debugln("Skipping Terraform summary because no summary output file is configured.")

				return true
			}

			if P.Plan.Output == "" {
				t.Log.Debugln("Skipping Terraform summary because no plan output file is configured.")

				return true
			}

			return false
		}).
		Set(func(t *Task) error {
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
					summary, err := summarizeTerraformShowPlan([]byte(strings.Join(c.GetStdoutStream(), "")))
					if err != nil {
						return err
					}

					body, err := renderSummary(summary)
					if err != nil {
						return err
					}

					output := P.Summary.Output
					if !filepath.IsAbs(output) {
						output = filepath.Join(setup.P.Project.Cwd, output)
					}

					if err := os.WriteFile(output, body, 0o644); err != nil {
						return fmt.Errorf("write Terraform summary %s: %w", output, err)
					}

					t.Log.Infof("Wrote Terraform summary: %s", output)

					return nil
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
			if !P.MergeRequestReport.Enabled {
				return true
			}

			if P.MergeRequestReport.MergeRequestId == 0 {
				t.Log.Debugln("Skipping GitLab merge request report because this is not a merge request pipeline.")

				return true
			}

			return false
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
					metadata := P.ReportMetadata
					metadata.Target = terraformStateName()
					metadata.Cwd = setup.P.Project.Cwd

					report, err := parseTerraformShowPlan([]byte(strings.Join(c.GetStdoutStream(), "")), metadata)
					if err != nil {
						return err
					}

					body, err := iac.RenderMergeRequestReport(report)
					if err != nil {
						return err
					}

					config := P.MergeRequestReport
					config.Identifier = gitlab.ResolveReportIdentifier(
						config.Identifier,
						metadata.JobName,
						terraformReportDiscriminators()...,
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
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerraformPlanCleanup(tl *TaskList) *Task {
	return tl.CreateTask("cleanup").
		ShouldDisable(func(t *Task) bool {
			if !P.Plan.PreviewForMergeRequests || P.Plan.PipelineSource != "merge_request_event" {
				return true
			}

			if P.Plan.Output == "" {
				return true
			}

			return false
		}).
		Set(func(t *Task) error {
			output := P.Plan.Output
			if !filepath.IsAbs(output) {
				output = filepath.Join(setup.P.Project.Cwd, output)
			}

			if err := os.Remove(output); err != nil {
				return fmt.Errorf("remove Terraform plan file %s: %w", output, err)
			}

			t.Log.Infof("Removed Terraform plan file: %s", output)

			return nil
		})
}
