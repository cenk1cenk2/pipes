package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	CiVariables struct {
		ProjectId string
		ApiUrl    string
	}

	Pipe struct {
		tool.Config
		CiVariables
		LogLevel string `validate:"omitempty,oneof=trace debug info warn error"`
	}
)

var P = &Pipe{}
var C = tool.NewCtx()

// Step is what every terraform command starts with, since the tool has to be
// resolved before anything can be asked of it.
var Step = cli.Step{Flags: Flags, New: New}

func New(p *plumber.Plumber) *plumber.TaskList {
	return tool.Setup(p, Spec, &P.Config, C, GenerateTerraformEnvVars).
		ShouldRunBefore(func(_ *plumber.TaskList) error {
			return cli.Validated(p, P)
		})
}
