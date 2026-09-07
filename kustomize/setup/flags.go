package setup

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/tool"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CATEGORY_KUSTOMIZE = "Kustomize"
)

var Spec = tool.Spec{
	Name:        "kustomize",
	Category:    CATEGORY_KUSTOMIZE,
	FlagPrefix:  "kustomize",
	EnvPrefix:   "KUSTOMIZE",
	VersionArgs: []string{"version"},
}

var Flags = CombineFlags(
	tool.NewFlags(Spec, &P.Config),
	[]cli.Flag{
		&cli.StringSliceFlag{
			Category:    CATEGORY_KUSTOMIZE,
			Name:        "kustomize.paths",
			Sources:     flags.EnvVars("KUSTOMIZE_PATHS"),
			Usage:       "Explicit overlay paths to build relative to the working directory.",
			Required:    false,
			Destination: &P.Paths,
		},
	})
