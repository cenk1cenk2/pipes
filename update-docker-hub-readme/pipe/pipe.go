package pipe

import (
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	DockerHub struct {
		Username string
		Password string
		Address  string
	}

	Readme struct {
		Repository  string
		File        string
		Description string
	}

	Plugin struct {
		DockerHub DockerHub
		Readme    Readme
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			LoginToDockerHubRegistry(),
			UpdateDockerReadme(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
