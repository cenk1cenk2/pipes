package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

type (
	Image struct {
		Platforms      []string
		Name           string
		Tags           []string
		TagAsLatest    []string
		TagsFile       string
		TagsFileStrict bool
		TagsSanitize   []versions.Match
		TagsTemplate   []versions.Match
		Pull           bool
		Push           bool
		BuildArgs      map[string]string
		LatestTag      string
		Cache          string `validate:"omitempty,dirpath"`
		Format         string `validate:"oneof=oci docker"`
		StorageDriver  string `validate:"oneof=overlay overlay2 vfs"`
	}

	File struct {
		Context string
		Name    string
	}

	Manifest struct {
		Target string
		File   string
	}

	Pipe struct {
		Git git.Refs
		Image
		File
		Manifest
	}

	Ctx struct {
		Tags []string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				tags(tl).Job(),
				build(tl).Job(),
				push(tl).Job(),
			)
		})
}
