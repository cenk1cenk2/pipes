package pipe

import (
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
	"github.com/urfave/cli/v2"
)

type (
	Gitlab struct {
		Token             string
		JobToken          string
		ParentProjectId   string
		ParentPipelineId  string
		DownloadArtifacts cli.StringSlice
	}

	Plugin struct {
		Gitlab Gitlab
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			TaskVerifyVariables(),
			TaskDiscoverArtifacts(),
			TaskDownloadArtifacts(),
			TaskUnarchiveArtifacts(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
