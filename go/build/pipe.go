package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
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

	// Deps is the resolved go tool: the directory the build runs in and the
	// environment the cache setup has written into.
	Deps struct {
		Tool *tool.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GoBuild(tl, deps).Job(),
			)
		})
}
