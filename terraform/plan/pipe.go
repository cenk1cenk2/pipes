package plan

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
)

type (
	Plan struct {
		Args       string
		Output     string
		RetryDelay time.Duration
		RetryTries uint32
	}

	Pipe struct {
		Plan
		MergeRequestReport gitlab.MergeRequestReportConfig
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				TerraformPlan(tl).Job(),
				TerraformMergeRequestReport(tl).Job(),
			)
		})
}
