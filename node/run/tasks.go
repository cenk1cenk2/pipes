package run

import (
	"os"
	"os/exec"
	"strings"

	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
	"gitlab.kilic.dev/devops/pipes/node/pipe"
)

type Ctx struct {
	EnvironmentVariables []string
	SelectedEnvironment  string
	FallbackEnvironment  string
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify-script"},
		Task: func(t *utils.Task) error {
			err := utils.ValidateAndSetDefaults(t.Metadata, &Pipe)

			if err != nil {
				return err
			}

			return nil
		},
	}
}

func RunNodeScript() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "script"},
		Task: func(t *utils.Task) error {
			args := []string{}

			cmd := exec.Command(pipe.Context.PackageManager.Exe)

			args = append(args, pipe.Context.PackageManager.Commands.Run...)
			args = append(args, Pipe.NodeCommand.Script)
			args = append(args, pipe.Context.PackageManager.Commands.RunDelimitter...)
			args = append(args, strings.Split(Pipe.NodeCommand.ScriptArgs, " ")...)

			cmd.Args = append(cmd.Args, args...)

			cmd.Dir = Pipe.NodeCommand.Cwd

			cmd.Env = append(cmd.Env, os.Environ()...)
			cmd.Env = append(cmd.Env, Context.EnvironmentVariables...)

			t.Command = cmd

			return nil
		},
	}
}
