package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryBuildah = "Buildah"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    CategoryBuildah,
		Name:        "buildah.cwd",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("BUILDAH_CWD")),
		Usage:       "Working directory for buildah commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},
}
