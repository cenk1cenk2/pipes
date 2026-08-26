package plan

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
	"gitlab.kilic.dev/devops/pipes/terraform/state"
)

type (
	Plan struct {
		Args                    string
		Output                  string
		PipelineSource          string
		PreviewForMergeRequests bool
		RetryDelay              time.Duration
		RetryTries              uint32
	}

	Summary struct {
		Output string
	}

	Pipe struct {
		Plan
		Summary
		MergeRequestReport gitlab.MergeRequestReportConfig
		ReportMetadata     iac.Metadata
	}

	// Deps is everything the plan reads from the pipe around it: the resolved
	// terraform tool it runs with, the state whose name tells concurrent plan jobs
	// on one merge request apart, and the notes the report is written through.
	Deps struct {
		Tool  *tool.Ctx
		State *state.Pipe
		Notes gitlab.NotesFactory
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if !P.MergeRequestReport.Enabled {
				P.MergeRequestReport.MergeRequestIid = 0
			}

			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			source := TerraformReportSource(deps)

			return JobSequence(
				TerraformPlan(tl, deps).Job(),
				iac.SummaryTask(tl, source).Job(),
				iac.MergeRequestReportTask(tl, source).Job(),
				TerraformPlanCleanup(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
