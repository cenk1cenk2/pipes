package pipe

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type Ctx struct {
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify"},
		Task: func(t *utils.Task) error {
			err := utils.ValidateAndSetDefaults(t.Metadata, &Pipe)

			if err != nil {
				return err
			}

			return nil
		},
	}
}

func InstallPackages() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "install"},
		Task: func(t *utils.Task) error {
			apks := Pipe.Apk.Value()

			if len(apks) > 0 {
				t.Log.Debugf(
					"Will install packages from APK repository: %s",
					strings.Join(apks, ", "),
				)

				cmd := exec.Command("apk", "--no-cache")

				cmd.Args = append(cmd.Args, apks...)

				t.Commands = append(t.Commands, cmd)
			}

			packages := Pipe.Node.Value()

			if len(packages) > 0 {
				t.Log.Debugf(
					"Will install packages from NPM repository: %s",
					strings.Join(packages, ", "),
				)

				cmd := exec.Command("yarn", "global", "add")

				cmd.Args = append(cmd.Args, packages...)

				t.Commands = append(t.Commands, cmd)
			}

			return nil
		},
	}
}

func RunSemanticRelease() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "semantic-release"},
		Task: func(t *utils.Task) error {
			var exe string
			if Pipe.UseMulti {
				exe = MULTI_SEMANTIC_RELEASE_EXE
			} else {
				exe = SEMANTIC_RELEASE_EXE
			}

			cmd := exec.Command(exe)

			if Pipe.SemanticRelease.IsDryRun {
				cmd.Args = append(cmd.Args, "--dry-run")
			}

			if t.Log.Logger.Level == logrus.DebugLevel {
				cmd.Args = append(cmd.Args, "--debug")
			}

			t.Commands = append(t.Commands, cmd)

			return nil
		},
	}
}
