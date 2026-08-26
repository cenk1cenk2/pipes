package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func TerraformLint(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask().
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				TerraformFmtCheck(t.TL, deps).Job(),
				TerraformValidate(t.TL, deps).Job(),
			)
		})
}

func TerraformFmtCheck(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("fmt", "check").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"fmt",
				"-check",
				"-diff",
				"-recursive",
			).
				Set(func(c *Command) error {
					if P.Lint.FormatCheckArgs != "" {
						c.AppendArgs(P.Lint.FormatCheckArgs)
					}

					return nil
				}).
				SetDir(deps.Tool.Cwd).
				AppendEnvironment(deps.Tool.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerraformValidate(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("validate").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"validate",
			).
				Set(func(c *Command) error {
					if P.Lint.ValidateArgs != "" {
						c.AppendArgs(P.Lint.ValidateArgs)
					}

					return nil
				}).
				SetDir(deps.Tool.Cwd).
				AppendEnvironment(deps.Tool.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
