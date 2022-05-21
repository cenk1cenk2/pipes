package pipe

import (
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	DockerImage struct {
		TagsFile string
	}

	Github struct {
		Token      string
		Repository string
	}

	Plugin struct {
		Github
		DockerImage
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			GithubLogin(),
			FetchLatestTag(),
			WriteTagsFile(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
