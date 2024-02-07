package install

import (
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func TerraformInstall(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("install").
		Set(func(t *Task[Pipe]) error {
			for _, ws := range setup.TL.Pipe.Project.Workspaces {
				t.CreateSubtask(ws).
					Set(func(t *Task[Pipe]) error {
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
							SetDir(ws).
							AppendEnvironment(setup.TL.Pipe.Ctx.EnvVars).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task[Pipe]) error {
						return t.RunCommandJobAsJobSequence()
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}
