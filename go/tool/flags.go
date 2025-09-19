package tool

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_TOOL = "Tool"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_TOOL,
		Name:     "go.tool",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_TOOL"),
		),
		Usage:       "Binary that provides the tooling.",
		Required:    false,
		Value:       "",
		Destination: &P.Tool,
	},

	&cli.StringFlag{
		Category: CATEGORY_TOOL,
		Name:     "go.tool.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_TOOL_ARGS"),
		),
		Usage:       "Arguments to append to tool command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},
}
