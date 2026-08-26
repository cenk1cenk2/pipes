package publish

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
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

	// Deps is what the setup step resolved: the directory the chart is packaged in
	// and the chart itself, whose name every archive is written and pushed under.
	Deps struct {
		Tool *setup.Ctx
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
			return JobSequence(
				HelmChartVersions(deps).Tasks(tl, &C.Versions).Job(),
				HelmPackage(tl, deps).Job(),
				HelmPublish(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
