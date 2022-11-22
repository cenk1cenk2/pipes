package pipe

import (
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Git flags.GitFlags

	Docker struct {
		UseBuildKit     bool
		UseBuildx       bool
		BuildxPlatforms string
		BuildxInstance  string
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

	DockerRegistry struct {
		Registry string
		Username string
		Password string
	}

	Pipe struct {
		Ctx

		Git
		Docker
		DockerImage
		DockerFile
		DockerRegistry
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				Setup(tl).Job(),

				DockerTagsParent(tl).Job(),

				tl.JobParallel(
					DockerVersion(tl).Job(),
					DockerBuildXVersion(tl).Job(),
				),

				DockerLoginParent(tl).Job(),

				tl.JobParallel(
					DockerBuildParent(tl).Job(),
					DockerBuildXParent(tl).Job(),
				),

				DockerInspect(tl).Job(),
			)
		})
}
