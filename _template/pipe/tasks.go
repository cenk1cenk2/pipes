package pipe

import (
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type Ctx struct {
}

var Context Ctx

func TaskVerifyVariables() utils.Task {
	metadata := utils.TaskMetadata{Context: "verify"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {
		return nil
	}}
}
