package setup

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

//revive:disable:line-length-limit

const (
	CATEGORY_PULUMI = "pulumi"
)

var Spec = tool.Spec{
	Name:        "pulumi",
	Category:    CATEGORY_PULUMI,
	FlagPrefix:  "pulumi",
	EnvPrefix:   "PULUMI",
	VersionArgs: []string{"version"},
}

var Flags = []cli.Flag(tool.Flags(Spec, P))
