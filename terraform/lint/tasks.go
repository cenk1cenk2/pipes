package lint

import (
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func LintTerraform(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("parent").
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return t.TL.JobParallel(
				TerraformFmtCheck(t.TL).Job(),
				TerraformValidate(t.TL).Job(),
			)
		})
}

func TerraformFmtCheck(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("fmt", "check").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Lint.FormatCheckEnable
		}).
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
				SetDir(setup.TL.Pipe.Project.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerraformValidate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("validate").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Lint.ValidateEnable
		}).
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
				SetDir(setup.TL.Pipe.Project.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
