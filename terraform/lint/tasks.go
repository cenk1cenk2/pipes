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
	return tl.CreateTask().
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Lint.FormatCheckEnable
		}).
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return tl.JobParallel(
				TerraformFmtCheckCwd(tl).Job(),
				TerraformFmtCheckWorkspace(tl).Job(),
			)
		})
}

func TerraformFmtCheckCwd(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("fmt", "check", setup.TL.Pipe.Cwd).
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(setup.TL.Pipe.Project.Workspaces) > 0
		}).
		Set(func(t *Task[Pipe]) error {
			return createFmtCheckCommand(t, setup.TL.Pipe.Project.Cwd)
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerraformFmtCheckWorkspace(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("fmt", "check").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(setup.TL.Pipe.Project.Workspaces) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			for _, ws := range setup.TL.Pipe.Project.Workspaces {
				func(ws string) {
					t.CreateSubtask(ws).
						Set(func(t *Task[Pipe]) error {
							return createFmtCheckCommand(t, ws)
						}).
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunCommandJobAsJobSequence()
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
	return tl.CreateTask().
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Lint.ValidateEnable
		}).
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return tl.JobParallel(
				TerraformValidateCwd(tl).Job(),
				TerraformValidateWorkspace(tl).Job(),
			)
		})
}

func TerraformValidateCwd(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("validate", setup.TL.Pipe.Cwd).
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(setup.TL.Pipe.Project.Workspaces) > 0
		}).
		Set(func(t *Task[Pipe]) error {
			return createValidateCommand(t, setup.TL.Pipe.Project.Cwd)
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerraformValidateWorkspace(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("validate").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(setup.TL.Pipe.Project.Workspaces) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			for _, ws := range setup.TL.Pipe.Project.Workspaces {
				func(ws string) {
					t.CreateSubtask(ws).
						Set(func(t *Task[Pipe]) error {
							return createValidateCommand(t, ws)
						}).
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunCommandJobAsJobSequence()
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

func createFmtCheckCommand(t *Task[Pipe], cwd string) error {
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
		SetDir(cwd).
		AppendEnvironment(setup.TL.Pipe.Ctx.EnvVars).
		AddSelfToTheTask()

	return nil
}

func createValidateCommand(t *Task[Pipe], cwd string) error {
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
		SetDir(cwd).
		AppendEnvironment(setup.TL.Pipe.Ctx.EnvVars).
		AddSelfToTheTask()

	return nil
}
