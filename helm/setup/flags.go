package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryHelm = "Helm"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CategoryHelm,
		Name:     "helm.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_CWD"),
			cli.EnvVar("HELM_ROOT"),
		),
		Usage:       "Working directory for helm commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},
}
