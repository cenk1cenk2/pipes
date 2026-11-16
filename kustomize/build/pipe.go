package build

import (
	"cmp"
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type (
	Pipe struct {
		EnableHelm     bool
		HelmCommand    string
		LoadRestrictor string `validate:"omitempty,oneof=rootOnly none"`
		KubeVersion    string
		Output         string `validate:"omitempty,dirpath"`
		// the commit the overlays are rendered on again, so the merge request report
		// shows what the change does to the rendered manifests.
		DiffBase           string
		TargetBranchSha    string
		DiffBaseSha        string
		MergeRequestReport gitlab.MergeRequestReportConfig
		ReportMetadata     terraform.Metadata
		ReportLog          terraform.LogConfig
	}

	OverlayResult struct {
		Overlay   string
		Yaml      []byte
		Resources []*yaml.RNode
		DocCount  int
		Err       error
	}

	Ctx struct {
		Results []OverlayResult
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

// Dials only once a report is actually written, and is swapped out by the specs, which
// have no GitLab to talk to.
var newNotes gitlab.NotesFactory = gitlab.NewNotes

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			if !P.MergeRequestReport.Enabled {
				P.MergeRequestReport.MergeRequestIid = 0
			}

			P.DiffBase = cmp.Or(P.DiffBase, P.TargetBranchSha, P.DiffBaseSha)

			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				build(tl).Job(),
				report(tl).Job(),
			)
		})
}
