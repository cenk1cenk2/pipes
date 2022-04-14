package pipe

import (
	"github.com/urfave/cli/v2"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	Markdown struct {
		Patterns  cli.StringSlice
		Arguments string
	}

	Plugin struct {
		Markdown Markdown
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			FindMarkdownFiles(),
			RunMarkdownToc(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
