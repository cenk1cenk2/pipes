package pipe

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
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
