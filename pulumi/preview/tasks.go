package preview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
)

// Only the values that actually vary between concurrent preview jobs on one merge
// request belong in the marker, since anything else changes the identifier for every
// consumer without disambiguating anything.
func pulumiReportDiscriminators() []string {
	discriminators := []string{stack.P.Stack}

	if cwd := setup.C.Cwd; cwd != "" && cwd != "." {
		discriminators = append(discriminators, cwd)
	}

	return discriminators
}

func reportSource() terraform.Source {
	metadata := P.ReportMetadata
	metadata.Target = stack.P.Stack
	metadata.Cwd = setup.C.Cwd

	return terraform.Source{
		Read: func(ctx context.Context, t *Task) (terraform.Report, error) {
			planPath := P.Plan
			if !filepath.IsAbs(planPath) {
				planPath = filepath.Join(setup.C.Cwd, planPath)
			}

			data, err := os.ReadFile(planPath)
			if err != nil {
				return terraform.Report{}, fmt.Errorf("read Pulumi plan file %s: %w", planPath, err)
			}

			return parsePulumiPlanReport(data, readStackState(ctx, t), metadata)
		},
		Summary:        terraform.Summarize,
		SummaryOutput:  P.Summary.Output,
		Cwd:            setup.C.Cwd,
		MergeRequest:   P.MergeRequestReport,
		Notes:          gitlab.NewNotes,
		Discriminators: pulumiReportDiscriminators,
		Metadata:       metadata,
		Log:            P.ReportLog,
	}
}

// The state is what turns a change into a comparison, but a report that only shows
// what the plan carries beats no report, so failing to read it does not fail the job.
func readStackState(ctx context.Context, t *Task) map[string]map[string]any {
	data, err := exportStackState(ctx, setup.C.Cwd, stack.P.Stack)
	if err == nil {
		var state map[string]map[string]any
		if state, err = parseStackState(data); err == nil {
			return state
		}
	}

	t.Log.Warn(fmt.Sprintf("Reporting only the values the plan carries, since the stack state is not readable: %s", err))

	return nil
}

func plan(tl *TaskList) *Task {
	return tl.CreateTask("plan").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"pulumi",
				"preview",
				"--non-interactive",
				"--diff",
				"--save-plan",
				P.Plan,
			).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
