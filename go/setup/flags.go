package setup

import (
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
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
	tool.Flags(Spec, &P.Config),
	[]ucli.Flag{
		&ucli.StringFlag{
			Category:    CATEGORY_SETUP,
			Name:        "go.cache",
			Sources:     cli.EnvVars("GO_CACHE"),
			Usage:       "Enable go cache.",
			Required:    false,
			Value:       "./.go/",
			Destination: &P.Cache,
		},
	})
