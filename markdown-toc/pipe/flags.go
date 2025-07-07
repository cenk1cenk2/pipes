package pipe

import (
	"github.com/urfave/cli/v3"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "markdown-toc.pattern",
		Usage:   "Pattern for markdown.",
		Value:   cli.NewStringSlice("README.md"),
		EnvVars: []string{"MARKDOWN_TOC_PATTERNS"},
	},
	&cli.IntFlag{
		Name:        "markdown-toc.start_depth",
		Usage:       "Start depth for the elements in the given document.",
		Value:       1,
		EnvVars:     []string{"MARKDOWN_TOC_START_DEPTH"},
		Destination: &TL.Pipe.Markdown.StartDepth,
	},
	&cli.IntFlag{
		Name:        "markdown-toc.end_depth",
		Usage:       "End depth for the elements in the given document.",
		Value:       5,
		EnvVars:     []string{"MARKDOWN_TOC_END_DEPTH"},
		Destination: &TL.Pipe.Markdown.EndDepth,
	},
	&cli.IntFlag{
		Name:        "markdown-toc.indentation",
		Usage:       "Indentation for each element.",
		Value:       2,
		EnvVars:     []string{"MARKDOWN_TOC_INDENTATION"},
		Destination: &TL.Pipe.Markdown.Indentation,
	},
	&cli.StringFlag{
		Name:        "markdown-toc.list_identifier",
		Usage:       "Identifier for each list element.",
		Value:       "-",
		EnvVars:     []string{"MARKDOWN_TOC_LIST_IDENTIFIER"},
		Destination: &TL.Pipe.Markdown.ListIdentifier,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.Markdown.Patterns = tl.CliContext.StringSlice("markdown-toc.pattern")

	return nil
}
