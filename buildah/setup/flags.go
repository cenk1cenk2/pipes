package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_BUILDAH = "Buildah"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    CATEGORY_BUILDAH,
		Name:        "buildah.cwd",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("BUILDAH_CWD")),
		Usage:       "Working directory for buildah commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},
}
