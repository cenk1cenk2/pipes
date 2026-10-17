package build

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

type (
	Image struct {
		Platforms           []string
		Name                string
		Tags                []string
		TagsFromEnvironment []string
		TagAsLatest         []string
		TagsSanitize        []versions.Match
		TagsTemplate        []versions.Match
		Pull                bool
		Push                bool
		BuildArgs           map[string]string
		LatestTag           string
		Cache               string `validate:"omitempty,dirpath"`
		Format              string `validate:"oneof=oci docker"`
		StorageDriver       string `validate:"oneof=overlay overlay2 vfs"`
	}

	File struct {
		Context string
		Name    string
	}

	Manifest struct {
		Target          string
		Job             string
		EnvFile         string
		PublishedPrefix string
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

// PushTarget is the image the build pushes, empty when it does not push.
func (i Image) PushTarget() string {
	if !i.Push {
		return ""
	}

	return i.Name
}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				image(tl).Job(),
				tags(tl).Job(),
				build(tl).Job(),
				push(tl).Job(),
				manifest(tl).Job(),
			)
		})
}
