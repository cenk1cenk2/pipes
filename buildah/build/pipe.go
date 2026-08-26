package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
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

	// Deps is the registry the login step authenticated against, whose uri every
	// tag is prefixed with so the image is built under the name it is pushed to.
	Deps struct {
		Registry *registry.Credentials
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			collector := ContainerImageTags(deps)

			return JobSequence(
				collector.Tasks(tl, &C.Tags).Job(),
				ContainerManifestFileWrite(tl, collector).Job(),
				ContainerBuild(tl).Job(),
				ContainerPush(tl).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
