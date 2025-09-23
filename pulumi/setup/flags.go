package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_PULUMI = "pulumi"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_PULUMI,
		Name:     "pulumi.cwd",
		Usage:    "Path to the Pulumi working directory.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PULUMI_CWD"),
		),
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},
}
