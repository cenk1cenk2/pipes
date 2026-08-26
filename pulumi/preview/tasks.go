package preview

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
)

// Only the values that actually vary between concurrent preview jobs on one merge
// request belong in the marker, since anything else changes the identifier for every
// consumer without disambiguating anything.
func pulumiReportDiscriminators(deps Deps) []string {
	discriminators := []string{deps.Stack.Stack}

	if cwd := deps.Tool.Cwd; cwd != "" && cwd != "." {
		discriminators = append(discriminators, cwd)
	}

	return discriminators
}

func PulumiReportSource(deps Deps) iac.Source {
	metadata := P.ReportMetadata
	metadata.Target = deps.Stack.Stack
	metadata.Cwd = deps.Tool.Cwd

	return iac.Source{
		Read: func(_ *Task) (iac.Report, error) {
			planPath := P.Plan
			if !filepath.IsAbs(planPath) {
				planPath = filepath.Join(deps.Tool.Cwd, planPath)
			}

			data, err := os.ReadFile(planPath)
			if err != nil {
				return iac.Report{}, fmt.Errorf("read Pulumi plan file %s: %w", planPath, err)
			}

			return parsePulumiPlanReport(data, metadata)
		},
		Summary:        iac.Summarize,
		SummaryOutput:  P.Summary.Output,
		Cwd:            deps.Tool.Cwd,
		MergeRequest:   P.MergeRequestReport,
		Notes:          deps.Notes,
		Discriminators: func() []string { return pulumiReportDiscriminators(deps) },
		Metadata:       metadata,
	}
}

func PulumiPlan(tl *TaskList, deps Deps) *Task {
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
				SetDir(deps.Tool.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
