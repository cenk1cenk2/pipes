package plan

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
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
		ShouldRunBefore(func(tl *TaskList) error {
			if !P.MergeRequestReport.Enabled {
				P.MergeRequestReport.MergeRequestIid = 0
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			C.Report = TerraformReportSource()

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				TerraformPlan(tl).Job(),
				terraform.SummaryTask(tl, &C.Report).Job(),
				terraform.MergeRequestReportTask(tl, &C.Report).Job(),
				TerraformPlanCleanup(tl).Job(),
			)
		})
}
