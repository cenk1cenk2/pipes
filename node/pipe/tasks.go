package pipe

import (
	"fmt"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
)

type Ctx struct {
	PackageManager
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

			Context.PackageManager = PackageManager{
				Exe:      Pipe.Node.PackageManager,
				Commands: PackageManagers[Pipe.Node.PackageManager],
			}

			t.Log.Debugln(fmt.Sprintf("Using package manager: %s", Pipe.Node.PackageManager))

			return nil
		},
	}
}
