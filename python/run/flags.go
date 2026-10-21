package run

import (
	"github.com/urfave/cli/v3"
)

var Arguments = []cli.Argument{
	&cli.StringArgs{
		Name:        "arg",
		Min:         1,
		Max:         -1,
		UsageText:   "Command to run inside the installed environment.",
		Destination: &P.Command,
	},
}
