package install

import (
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func InitTerraform(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("install").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"terraform",
				"init",
				"-input=false",
			).
				Set(func(c *Command[Pipe]) error {
					if t.Pipe.Install.Reconfigure {
						c.AppendArgs("-reconfigure")
					}

					if t.Pipe.Install.Args != "" {
						c.AppendArgs(t.Pipe.Install.Args)
					}

					return nil
				}).
				SetDir(setup.TL.Pipe.Project.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
