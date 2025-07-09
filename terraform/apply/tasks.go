package apply

import (
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "github.com/cenk1cenk2/plumber/v6"
)

func TerraformApply(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("apply").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"terraform",
				"apply",
				"-input=false",
			).
				Set(func(c *Command[Pipe]) error {
					if t.Pipe.Apply.Output != "" {
						c.AppendArgs(t.Pipe.Apply.Output)
					}

					if t.Pipe.Apply.Args != "" {
						c.AppendArgs(t.Pipe.Apply.Args)
					}

					return nil
				}).
				SetDir(setup.TL.Pipe.Project.Cwd).
				AppendEnvironment(setup.TL.Pipe.Ctx.EnvVars).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
