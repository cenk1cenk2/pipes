package build

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryBuild = "Build"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CategoryBuild,
		Name:     "python.build.output",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PYTHON_BUILD_OUTPUT"),
		),
		Usage:       "Directory the distributions are written to, relative to the working directory.",
		Required:    false,
		Value:       "dist",
		Destination: &P.Output,
	},

	&cli.StringFlag{
		Category: CategoryBuild,
		Name:     "python.build.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PYTHON_BUILD_ARGS"),
		),
		Usage:       "Arguments to append to the build command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},
}
