package pipe

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
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
