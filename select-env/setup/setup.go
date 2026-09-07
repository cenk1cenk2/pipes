// Package setup holds the environment selection the pipe writes out. The
// instances live here rather than in main so that the flags that fill them and
// the task that reads them back are the same ones.
package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

var (
	Environment    = &environment.Config{}
	EnvironmentCtx = &environment.Ctx{}
)

// New selects the environment out of the source control references and reads
// its variables, which is everything the pipe does before writing them out.
func New(p *Plumber) *TaskList {
	return environment.SetupTaskList(p, Environment, EnvironmentCtx)
}
