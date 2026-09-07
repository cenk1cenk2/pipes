package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Args           string
		Output         string `validate:"dirpath"`
		BinaryName     string
		BinaryTemplate string
		LinkerFlags    string
		EnableCGO      bool
		BuildTargets   []GoBuildTarget
		BuildVariables map[string]string
		BuildTags      []string
	}

	GoBuildTarget struct {
		Os   string `json:"os,omitempty"   yaml:"os,omitempty"`
		Arch string `json:"arch,omitempty" yaml:"arch,omitempty"`
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GoBuild(tl).Job(),
			)
		})
}
