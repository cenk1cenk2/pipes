package plan

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
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
)

var TL = TaskList{}

var P = &Pipe{}

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
			source := TerraformReportSource()

			return JobSequence(
				TerraformPlan(tl).Job(),
				iac.SummaryTask(tl, source).Job(),
				iac.MergeRequestReportTask(tl, source).Job(),
				TerraformPlanCleanup(tl).Job(),
			)
		})
}
