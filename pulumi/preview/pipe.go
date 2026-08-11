package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
	"gitlab.kilic.dev/devops/pipes/common/report/iac"
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
				P.MergeRequestReport.MergeRequestId = 0
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				PulumiPlan(tl).Job(),
				PulumiSummary(tl).Job(),
				PulumiMergeRequestReport(tl).Job(),
			)
		})
}
