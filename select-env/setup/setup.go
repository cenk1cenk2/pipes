// Package setup holds the environment selection the pipe writes out, so that the
// flags that fill it and the task that reads it back are the same instances.
package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

var (
	Environment    = &environment.Config{}
	EnvironmentCtx = &environment.Ctx{}
)

// New selects the environment out of the source control references and reads its variables.
func New(p *Plumber) *TaskList {
	return environment.SetupTaskList(p, Environment, EnvironmentCtx)
}
