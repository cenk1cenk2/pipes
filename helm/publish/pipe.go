package publish

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

type (
	HelmChart struct {
		Name              string
		Versions          []string
		VersionFile       string
		VersionFileStrict bool
		VersionsSanitize  []HelmChartMatch
		VersionsTemplate  []HelmChartMatch
		Destination       string `validate:"dirpath"`
	}

	HelmChartMatch struct {
		Match    string `json:"match"    yaml:"match"    validate:"required"`
		Template string `json:"template" yaml:"template" validate:"required"`
	}

	Pipe struct {
		Git flags.GitFlags
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
				HelmChartVersionsParent(tl).Job(),
				HelmPackage(tl).Job(),
				HelmPublish(tl).Job(),
			)
		})
}
