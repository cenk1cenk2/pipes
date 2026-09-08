package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryKustomize = "Kustomize"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CategoryKustomize,
		Name:     "kustomize.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("KUSTOMIZE_CWD"),
			cli.EnvVar("KUSTOMIZE_ROOT"),
		),
		Usage:       "Working directory for kustomize commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	&cli.StringSliceFlag{
		Category:    CategoryKustomize,
		Name:        "kustomize.paths",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("KUSTOMIZE_PATHS")),
		Usage:       "Explicit overlay paths to build relative to the working directory.",
		Required:    false,
		Destination: &P.Paths,
	},
}
