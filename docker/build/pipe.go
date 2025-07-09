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
		Ctx

		Git
		Docker
		DockerImage
		DockerFile
		DockerManifest
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				Setup(tl).Job(),
				DockerTagsParent(tl).Job(),

				tl.JobParallel(
					DockerBuildParent(tl).Job(),
					DockerBuildXParent(tl).Job(),
				),

				DockerInspect(tl).Job(),
			)
		})
}
