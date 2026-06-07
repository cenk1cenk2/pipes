package preview

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	ReportConfig struct {
		Enabled bool
		Output  string `validate:"required_if=Enabled true"`
	}

	Pipe struct {
		Plan   string
		Report ReportConfig
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
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				PulumiPlan(tl).Job(),
				ReportTask(tl).Job(),
			)
		})
}
