package pipe

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:        "typedoc.pattern",
		Usage:       "Pattern for directories. Supports multiple patterns with comma-separated values.",
		Value:       cli.NewStringSlice("packages/*/"),
		EnvVars:     []string{"TYPEDOC_PATTERNS", "PLUGIN_TYPEDOC_PATTERNS"},
		Destination: &Pipe.TypeDoc.Patterns,
	},
	&cli.StringFlag{
		Name:        "typedoc.arguments",
		Usage:       "Pass in the arguments for TypeDoc.",
		Value:       "--options .typedoc.json --hideInPageTOC --hideBreadcrumbs",
		EnvVars:     []string{"TYPEDOC_ARGUMENTS", "PLUGIN_TYPEDOC_ARGUMENTS"},
		Destination: &Pipe.TypeDoc.Arguments,
	},
}
