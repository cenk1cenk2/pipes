package setup

import (
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

// The flags are built once, so the command registers the same slice the task list reads back.
var Flags = environment.NewFlags(environment.Options{Destination: Environment})
