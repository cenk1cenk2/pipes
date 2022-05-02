package pipe

import (
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	Gitlab struct {
		Url               string
		ApiUrl            string
		Token             string
		JobToken          string
		ParentProjectId   string
		ParentPipelineId  string
		DownloadArtifacts string
	}

	Plugin struct {
		Gitlab Gitlab
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			DiscoverArtifacts(),
			DownloadArtifacts(),
			UnarchiveArtifacts(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
