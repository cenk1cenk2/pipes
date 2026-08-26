package setup

import (
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CATEGORY_KUSTOMIZE = "Kustomize"
)

var Spec = tool.Spec{
	Name:          "kustomize",
	Category:      CATEGORY_KUSTOMIZE,
	FlagPrefix:    "kustomize",
	EnvPrefix:     "KUSTOMIZE",
	CwdEnvAliases: []string{"KUSTOMIZE_ROOT"},
	VersionArgs:   []string{"version"},
}

var Flags = CombineFlags(
	tool.Flags(Spec, &P.Config),
	[]ucli.Flag{
		&ucli.StringSliceFlag{
			Category:    CATEGORY_KUSTOMIZE,
			Name:        "kustomize.paths",
			Sources:     cli.EnvVars("KUSTOMIZE_PATHS"),
			Usage:       "Explicit overlay paths to build relative to the working directory.",
			Required:    false,
			Destination: &P.Paths,
		},
	})
