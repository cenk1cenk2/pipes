package apply

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
	"gitlab.kilic.dev/devops/pipes/terraform/plan"
)

type (
	Apply struct {
		Args   string
		Output string
	}

	Pipe struct {
		Apply
		MergeRequestReport gitlab.MergeRequestReportConfig
		ReportMetadata     terraform.Metadata
	}

	Ctx struct {
		Report terraform.Source
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			if !P.MergeRequestReport.Enabled {
				P.MergeRequestReport.MergeRequestIid = 0
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			C.Report = plan.ReportSource(plan.Pipe{
				Plan:               plan.Plan{Output: P.Apply.Output},
				MergeRequestReport: P.MergeRequestReport,
				ReportMetadata:     P.ReportMetadata,
			})

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				apply(tl).Job(),
			)
		})
}
