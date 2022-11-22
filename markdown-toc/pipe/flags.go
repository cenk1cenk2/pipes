package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "markdown-toc.pattern",
		Usage:   "Pattern for markdown.",
		Value:   cli.NewStringSlice("README.md"),
		EnvVars: []string{"MARKDOWN_TOC_PATTERNS", "PLUGIN_MARKDOWN_TOC_PATTERNS"},
	},
	&cli.StringFlag{
		Name:        "markdown-toc.arguments",
		Usage:       "Pass in the arguments for markdown-toc.",
		Value:       "--bullets='-'",
		EnvVars:     []string{"MARKDOWN_TOC_ARGUMENTS", "PLUGIN_MARKDOWN_TOC_ARGUMENTS"},
		Destination: &TL.Pipe.Markdown.Arguments,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.Markdown.Patterns = tl.CliContext.StringSlice("markdown-toc.pattern")

	return nil
}
