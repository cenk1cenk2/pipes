package pipe

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Name: "markdown-toc.pattern",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("MARKDOWN_TOC_PATTERNS"),
		),
		Usage:       "Pattern for markdown.",
		Value:       []string{"README.md"},
		Destination: &P.Markdown.Patterns,
	},
	&cli.IntFlag{
		Name: "markdown-toc.start_depth",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("MARKDOWN_TOC_START_DEPTH"),
		),
		Usage:       "Start depth for the elements in the given document.",
		Value:       1,
		Destination: &P.Markdown.StartDepth,
	},
	&cli.IntFlag{
		Name: "markdown-toc.end_depth",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("MARKDOWN_TOC_END_DEPTH"),
		),
		Usage:       "End depth for the elements in the given document.",
		Value:       5,
		Destination: &P.Markdown.EndDepth,
	},
	&cli.IntFlag{
		Name: "markdown-toc.indentation",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("MARKDOWN_TOC_INDENTATION"),
		),
		Usage:       "Indentation for each element.",
		Value:       2,
		Destination: &P.Markdown.Indentation,
	},
	&cli.StringFlag{
		Name: "markdown-toc.list_identifier",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("MARKDOWN_TOC_LIST_IDENTIFIER"),
		),
		Usage:       "Identifier for each list element.",
		Value:       "-",
		Destination: &P.Markdown.ListIdentifier,
	},
}
