package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Gitlab struct {
		ApiUrl            string
		Token             string
		JobToken          string
		ParentProjectId   string
		ParentPipelineId  string
		DownloadArtifacts string
	}

	Pipe struct {
		Ctx

		Gitlab
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
				DiscoverArtifacts(tl).Job(),
				DownloadArtifacts(tl).Job(),
				UnarchiveArtifacts(tl).Job(),
			)
		})
}
