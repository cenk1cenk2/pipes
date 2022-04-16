package pipe

import (
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
