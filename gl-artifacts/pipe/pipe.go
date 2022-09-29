package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
	return TL.New(p).SetTasks(
		TL.JobSequence(
			Setup(&TL).Job(),
			DiscoverArtifacts(&TL).Job(),
			DownloadArtifacts(&TL).Job(),
			UnarchiveArtifacts(&TL).Job(),
		),
	)
}
