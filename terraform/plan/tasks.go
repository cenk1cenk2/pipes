package plan

import (
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func TerraformPlan(tl *TaskList) *Task {
	return tl.CreateTask("plan").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"plan",
				"-input=false",
			).
				Set(func(c *Command) error {
					if P.Plan.Output != "" {
						c.AppendArgs(fmt.Sprintf("-out=%s", P.Plan.Output))
					}

					if P.Plan.Args != "" {
						c.AppendArgs(P.Plan.Args)
					}

					return nil
				}).
				SetDir(setup.P.Project.Cwd).
				AppendEnvironment(setup.C.EnvVars).
				SetRetries(&CommandRetry{
					Tries: P.Plan.RetryTries,
					Delay: P.Plan.RetryDelay,
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
