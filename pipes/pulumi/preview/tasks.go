package preview

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
	"gitlab.kilic.dev/devops/pipes/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pipes/pulumi/stack"
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

func PulumiReportSource() iac.Source {
	metadata := P.ReportMetadata
	metadata.Target = stack.P.Stack
	metadata.Cwd = setup.C.Cwd

	return iac.Source{
		Read: func(_ *Task) (iac.Report, error) {
			planPath := P.Plan
			if !filepath.IsAbs(planPath) {
				planPath = filepath.Join(setup.C.Cwd, planPath)
			}

			data, err := os.ReadFile(planPath)
			if err != nil {
				return iac.Report{}, fmt.Errorf("read Pulumi plan file %s: %w", planPath, err)
			}

			return parsePulumiPlanReport(data, metadata)
		},
		Summary:        iac.Summarize,
		SummaryOutput:  P.Summary.Output,
		Cwd:            setup.C.Cwd,
		MergeRequest:   P.MergeRequestReport,
		Notes:          gitlab.NewNotes,
		Discriminators: pulumiReportDiscriminators,
		Metadata:       metadata,
	}
}

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
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
