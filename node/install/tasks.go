package install

import (
	"os/exec"

	pipe "gitlab.kilic.dev/devops/gitlab-pipes/node/pipe"
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
)

type Ctx struct {
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify install"},
		Task: func(t *utils.Task) error {
			err := utils.ValidateAndSetDefaults(t.Metadata, &Pipe)

			if err != nil {
				return err
			}

			return nil
		}}
}

func InstallNodeDependencies() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "node install"},
		Task: func(t *utils.Task) error {
			var args []string

			if Pipe.NodeInstall.UseLockFile {
				args = pipe.Context.PackageManager.Commands.InstallWithLock

				t.Log.Debugln("Using lockfile for installation.")
			} else {
				args = pipe.Context.PackageManager.Commands.Install

				t.Log.Debugln("Installing dependencies without a lockfile.")
			}

			cmd := exec.Command(pipe.Context.PackageManager.Exe)
			cmd.Args = append(cmd.Args, args...)

			cmd.Dir = Pipe.NodeInstall.Cwd

			t.Command = cmd

			return nil
		},
	}
}
