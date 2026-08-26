package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
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

	Ctx struct {
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if !P.MergeRequestReport.Enabled {
				P.MergeRequestReport.MergeRequestIid = 0
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			source := PulumiReportSource()

			return JobSequence(
				PulumiPlan(tl).Job(),
				iac.SummaryTask(tl, source).Job(),
				iac.MergeRequestReportTask(tl, source).Job(),
			)
		})
}
