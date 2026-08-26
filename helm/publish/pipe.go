package publish

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

type (
	HelmChart struct {
		Target            string
		Versions          []string
		VersionFile       string
		VersionFileStrict bool
		VersionsSanitize  []versions.Match
		VersionsTemplate  []versions.Match
		Destination       string `validate:"dirpath"`
		AppVersion        string
	}

	Pipe struct {
		Git git.Refs
		HelmChart
	}

	Ctx struct {
		Versions []string
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
				HelmChartVersions().Tasks(tl, &C.Versions).Job(),
				HelmPackage(tl).Job(),
				HelmPublish(tl).Job(),
			)
		})
}
