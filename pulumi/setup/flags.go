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
		Category:    CATEGORY_PULUMI,
		Name:        "pulumi.cwd",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("PULUMI_CWD")),
		Usage:       "Working directory for pulumi commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},
}
