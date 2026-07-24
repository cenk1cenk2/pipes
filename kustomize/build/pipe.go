package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
)

type (
	Pipe struct {
		EnableHelm     bool
		HelmCommand    string
		LoadRestrictor string `validate:"omitempty,oneof=rootOnly none"`
		KubeVersion    string
		Output         string `validate:"omitempty,dirpath"`
		FailFast       bool

		MergeRequestReport gitlab.MergeRequestReportConfig
	}

	OverlayResult struct {
		Overlay  string
		Yaml     []byte
		DocCount int
		Err      error
	}

	Ctx struct {
		Results []OverlayResult
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
				RenderOverlays(tl).Job(),
				MergeRequestReport(tl).Job(),
			)
		})
}
