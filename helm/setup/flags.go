package setup

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

//revive:disable:line-length-limit

const (
	CATEGORY_HELM = "Helm"
)

var Spec = tool.Spec{
	Name:          "helm",
	Category:      CATEGORY_HELM,
	FlagPrefix:    "helm",
	EnvPrefix:     "HELM",
	CwdEnvAliases: []string{"HELM_ROOT"},
	VersionArgs:   []string{"version"},
}

var Flags = []cli.Flag(tool.Flags(Spec, P))
