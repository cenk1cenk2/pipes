package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

type (
	Git flags.GitFlags

	ContainerImage struct {
		Platforms      []string
		Name           string
		Tags           []string
		TagAsLatest    []string
		TagsFile       string
		TagsFileStrict bool
		TagsSanitize   []ContainerImageMatch
		TagsTemplate   []ContainerImageMatch
		Pull           bool
		Push           bool
		BuildArgs      map[string]string
		LatestTag      string
		Cache          string `validate:"omitempty,dirpath"`
		Format         string `validate:"oneof=oci docker"`
		StorageDriver  string `validate:"oneof=overlay overlay2 vfs"`
	}

	ContainerFile struct {
		Context string
		Name    string
	}

	ContainerManifest struct {
		Target string
		File   string
	}

	ContainerImageMatch struct {
		Match    string `json:"match"    yaml:"match"    validate:"required"`
		Template string `json:"template" yaml:"template" validate:"required"`
	}

	Pipe struct {
		Git
		ContainerImage
		ContainerFile
		ContainerManifest
	}

	Ctx struct {
		Tags       []string
		References []string
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
				ParseReferences(tl).Job(),
				ContainerImageTagsParent(tl).Job(),
				ContainerBuild(tl).Job(),
				ContainerPush(tl).Job(),
			)
		})
}
