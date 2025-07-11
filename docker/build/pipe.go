package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

type (
	Git flags.GitFlags

	Docker struct {
		BuildxPlatforms string
	}

	DockerImage struct {
		Name           string
		Tags           []string
		TagAsLatest    []string
		TagsFile       string
		TagsFileStrict bool
		TagsSanitize   []TagsSanitizeJson
		TagsTemplate   []TagsTemplateJson
		Pull           bool
		Inspect        bool
		BuildArgs      []string
	}

	DockerFile struct {
		Context string
		Name    string
	}

	DockerManifest struct {
		Target     string
		OutputFile string
	}

	Pipe struct {
		Git
		Docker
		DockerImage
		DockerFile
		DockerManifest
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
				Setup(tl).Job(),
				DockerTagsParent(tl).Job(),

				JobParallel(
					DockerBuildParent(tl).Job(),
					DockerBuildXParent(tl).Job(),
				),

				DockerInspect(tl).Job(),
			)
		})
}
