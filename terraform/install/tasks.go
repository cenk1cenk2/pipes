package install

import (
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func TerraformInstall(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("install").
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return tl.JobParallel(
				TerraformInstallCwd(tl).Job(),
				TerraformInstallWorkspace(tl).Job(),
			)
		})
}

func TerraformInstallCwd(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask(setup.TL.Pipe.Cwd).
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(setup.TL.Pipe.Project.Workspaces) > 0
		}).
		Set(func(t *Task[Pipe]) error {
			return createInstallCmd(t, setup.TL.Pipe.Project.Cwd)
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerraformInstallWorkspace(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask().
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(setup.TL.Pipe.Project.Workspaces) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			for _, ws := range setup.TL.Pipe.Project.Workspaces {
				func(ws string) {
					t.CreateSubtask(ws).
						Set(func(t *Task[Pipe]) error {
							return createInstallCmd(t, ws)
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

func createInstallCmd(t *Task[Pipe], cwd string) error {
	t.CreateCommand(
		"terraform",
		"init",
		"-input=false",
	).
		Set(func(c *Command[Pipe]) error {
			if t.Pipe.Install.Reconfigure {
				t.Log.Infoln("Will reconfigure state.")

				c.AppendArgs("-reconfigure")
			}

			if t.Pipe.Install.UseLockfile {
				t.Log.Infoln("Using lockfile.")

				c.AppendArgs("-lockfile=readonly")
			}

			if t.Pipe.Install.Args != "" {
				c.AppendArgs(t.Pipe.Install.Args)
			}

			return nil
		}).
		SetDir(cwd).
		AppendEnvironment(setup.TL.Pipe.Ctx.EnvVars).
		AddSelfToTheTask()

	return nil
}
