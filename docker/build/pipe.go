package build

import (
	"encoding/json"
	"fmt"

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
var raw = &struct {
	DockerImageTagAsLatest  string
	DockerImageTagsSanitize string
	DockerImageTagsTemplate string
}{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if v := raw.DockerImageTagAsLatest; v != "" {
				if err := json.Unmarshal([]byte(v), &P.DockerImage.TagAsLatest); err != nil {
					return fmt.Errorf("Can not unmarshal Docker image tags for latest: %w", err)
				}
			}

			if v := raw.DockerImageTagsSanitize; v != "" {
				if err := json.Unmarshal([]byte(v), &P.DockerImage.TagsSanitize); err != nil {
					return fmt.Errorf("Can not unmarshal Docker image sanitizing tag conditions: %w", err)
				}
			}

			if v := raw.DockerImageTagsTemplate; v != "" {
				if err := json.Unmarshal([]byte(v), &P.DockerImage.TagsTemplate); err != nil {
					return fmt.Errorf("Can not unmarshal Docker image templating tag conditions: %w", err)
				}
			}

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
