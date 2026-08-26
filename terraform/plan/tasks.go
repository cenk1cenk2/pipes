package plan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
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

	if cwd := setup.C.Cwd; cwd != "" && cwd != "." {
		discriminators = append(discriminators, cwd)
	}

	return discriminators
}

func TerraformReportSource() iac.Source {
	metadata := P.ReportMetadata
	metadata.Target = terraformStateName()
	metadata.Cwd = setup.C.Cwd

	// terraform show reads the plan back out of the file terraform plan wrote, so
	// without one there is nothing to summarize.
	summaryOutput := P.Summary.Output
	if P.Plan.Output == "" {
		summaryOutput = ""
	}

	return iac.Source{
		Read: func(t *Task) (iac.Report, error) {
			if P.Plan.Output == "" {
				return iac.Report{}, fmt.Errorf("terraform plan output is required for the plan report")
			}

			show := t.CreateCommand(
				"terraform",
				"show",
				"-json",
				P.Plan.Output,
			).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				SetLogLevel(LOG_LEVEL_TRACE, LOG_LEVEL_WARN, LOG_LEVEL_DEBUG).
				EnableStreamRecording()

			if err := show.Run(); err != nil {
				return iac.Report{}, err
			}

			return parseTerraformShowPlan([]byte(strings.Join(show.GetStdoutStream(), "")), metadata)
		},
		Summary:        iac.Summarize,
		SummaryOutput:  summaryOutput,
		Cwd:            setup.C.Cwd,
		MergeRequest:   P.MergeRequestReport,
		Notes:          gitlab.NewNotes,
		Discriminators: terraformReportDiscriminators,
		Metadata:       metadata,
	}
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
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
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
				output = filepath.Join(setup.C.Cwd, output)
			}

			if err := os.Remove(output); err != nil {
				return fmt.Errorf("remove Terraform plan file %s: %w", output, err)
			}

			t.Log.Infof("Removed Terraform plan file: %s", output)

			return nil
		})
}
