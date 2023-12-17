package plan

import (
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func TerraformPlan(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("plan").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"terraform",
				"plan",
				"-input=false",
			).
				Set(func(c *Command[Pipe]) error {
					if t.Pipe.Plan.Output != "" {
						c.AppendArgs(t.Pipe.Plan.Output)
					}

					if t.Pipe.Plan.Args != "" {
						c.AppendArgs(t.Pipe.Plan.Args)
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
