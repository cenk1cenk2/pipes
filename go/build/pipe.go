package build

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
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

	// Ctx holds the directories a binary is built out of, which leaves out the
	// library modules a workspace carries next to the commands.
	Ctx struct {
		Packages []string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				packages(tl).Job(),
				build(tl).Job(),
			)
		})
}
