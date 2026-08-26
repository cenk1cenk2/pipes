package setup

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

//revive:disable:line-length-limit

const (
	CATEGORY_BUILDAH = "Buildah"
)

var Spec = tool.Spec{
	Name:        "buildah",
	Category:    CATEGORY_BUILDAH,
	FlagPrefix:  "buildah",
	EnvPrefix:   "BUILDAH",
	VersionArgs: []string{"--version"},
}

var Flags = []cli.Flag(tool.Flags(Spec, P))
