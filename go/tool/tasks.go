package tool

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func GoTool(tl *TaskList) *Task {
	return CreateGoToolTask(tl, CreateGoToolTaskOptions{
		Cwd:  setup.P.Cwd,
		Tool: P.Tool,
		Args: P.Args,
	})
}

type CreateGoToolTaskOptions struct {
	Cwd  string
	Tool string
	Args string
}

func CreateGoToolTask(tl *TaskList, options CreateGoToolTaskOptions) *Task {
	return tl.CreateTask("tool", options.Tool).
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"tool",
			).
				SetDir(options.Cwd).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.AppendArgs(options.Tool)

					c.AppendArgs(strings.Split(options.Args, " ")...)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
