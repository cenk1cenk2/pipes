package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
)

type (
	Summary struct {
		Output string
	}

	Pipe struct {
		Plan string
		Summary
		MergeRequestReport gitlab.MergeRequestReportConfig
		ReportMetadata     iac.Metadata
	}

	// Deps is everything the preview reads from the pipe around it: the resolved
	// pulumi tool it runs with, the stack whose name tells concurrent preview jobs
	// on one merge request apart, and the notes the report is written through.
	Deps struct {
		Tool  *tool.Ctx
		Stack *stack.Pipe
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
			source := PulumiReportSource(deps)

			return JobSequence(
				PulumiPlan(tl, deps).Job(),
				iac.SummaryTask(tl, source).Job(),
				iac.MergeRequestReportTask(tl, source).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
