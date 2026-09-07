package setup

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/tool"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CATEGORY_SETUP = "Setup"
)

var Spec = tool.Spec{
	Name:        "go",
	Category:    CATEGORY_SETUP,
	FlagPrefix:  "go",
	EnvPrefix:   "GO",
	VersionArgs: []string{"version"},
}

var Flags = CombineFlags(
	tool.NewFlags(Spec, &P.Config),
	[]cli.Flag{
		&cli.StringFlag{
			Category:    CATEGORY_SETUP,
			Name:        "go.cache",
			Sources:     flags.EnvVars("GO_CACHE"),
			Usage:       "Enable go cache.",
			Required:    false,
			Value:       "./.go/",
			Destination: &P.Cache,
		},

		&cli.BoolFlag{
			Category:    CATEGORY_SETUP,
			Name:        "go.workspace",
			Sources:     flags.EnvVars("GO_WORKSPACE"),
			Usage:       "Drive the modules as a Go workspace instead of the single module in the working directory.",
			Required:    false,
			Value:       false,
			Destination: &P.Workspace,
		},
	})
