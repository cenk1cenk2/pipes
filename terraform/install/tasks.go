package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func TerraformInstall(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"init",
				"-input=false",
			).
				Set(func(c *Command) error {
					if P.Install.Reconfigure {
						t.Log.Infoln("Will reconfigure state.")

						c.AppendArgs("-reconfigure")
					}

					if P.Install.UseLockfile {
						t.Log.Infoln("Using lockfile.")

						c.AppendArgs("-lockfile=readonly")
					}

					if P.Install.Args != "" {
						c.AppendArgs(P.Install.Args)
					}

					return nil
				}).
				SetDir(setup.P.Cwd).
				AppendEnvironment(setup.C.EnvVars).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
