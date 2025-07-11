package apply

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func TerraformApply(tl *TaskList) *Task {
	return tl.CreateTask("apply").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"apply",
				"-input=false",
			).
				Set(func(c *Command) error {
					if P.Apply.Output != "" {
						c.AppendArgs(P.Apply.Output)
					}

					if P.Apply.Args != "" {
						c.AppendArgs(P.Apply.Args)
					}

					return nil
				}).
				SetDir(setup.P.Project.Cwd).
				AppendEnvironment(setup.C.EnvVars).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
