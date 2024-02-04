package lint

import (
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func TerraformLint(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask().
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return t.TL.JobParallel(
				TerraformFmtCheck(t.TL).Job(),
				TerraformValidate(t.TL).Job(),
			)
		})
}

func TerraformFmtCheck(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("fmt", "check").
		Set(func(t *Task[Pipe]) error {
			for _, ws := range setup.TL.Pipe.Project.Workspaces {
				func(ws string) {
					t.CreateSubtask(ws).
						Set(func(t *Task[Pipe]) error {
							t.CreateCommand(
								"terraform",
								"fmt",
								"-check",
								"-diff",
								"-recursive",
							).
								Set(func(c *Command[Pipe]) error {
									if t.Pipe.Lint.FormatCheckArgs != "" {
										c.AppendArgs(t.Pipe.Lint.FormatCheckArgs)
									}

									return nil
								}).
								SetDir(ws).
								AppendEnvironment(setup.TL.Pipe.Ctx.EnvVars).
								AddSelfToTheTask()

							return nil
						}).
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunCommandJobAsJobParallel()
						}).
						AddSelfToTheParentAsParallel()
				}(ws)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func TerraformValidate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("validate").
		Set(func(t *Task[Pipe]) error {
			for _, ws := range setup.TL.Pipe.Project.Workspaces {
				func(ws string) {
					t.CreateSubtask(ws).
						Set(func(t *Task[Pipe]) error {
							t.CreateCommand(
								"terraform",
								"validate",
							).
								Set(func(c *Command[Pipe]) error {
									if t.Pipe.Lint.ValidateArgs != "" {
										c.AppendArgs(t.Pipe.Lint.ValidateArgs)
									}

									return nil
								}).
								SetDir(ws).
								AppendEnvironment(setup.TL.Pipe.Ctx.EnvVars).
								AddSelfToTheTask()

							return nil
						}).
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunCommandJobAsJobParallel()
						}).
						AddSelfToTheParentAsParallel()
				}(ws)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}
