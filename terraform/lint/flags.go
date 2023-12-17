package lint

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "terraform-lint.fmt-check.enable",
		Usage:       "Enable terraform fmt.",
		Required:    false,
		EnvVars:     []string{"TF_LINT_FMT_CHECK_ENABLE"},
		Value:       true,
		Destination: &TL.Pipe.Lint.FormatCheckEnable,
	},
	&cli.StringFlag{
		Name:        "terraform-lint.fmt-check.args",
		Usage:       "Additional arguments for terraform fmt.",
		Required:    false,
		EnvVars:     []string{"TF_LINT_FMT_CHECK_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.Lint.FormatCheckArgs,
	},
	&cli.BoolFlag{
		Name:        "terraform-lint.validate.enable",
		Usage:       "Enable terraform validate.",
		Required:    false,
		EnvVars:     []string{"TF_LINT_VALIDATE_ENABLE"},
		Value:       true,
		Destination: &TL.Pipe.Lint.ValidateEnable,
	},
	&cli.StringFlag{
		Name:        "terraform-lint.validate.args",
		Usage:       "Additional arguments for terraform validate.",
		Required:    false,
		EnvVars:     []string{"TF_LINT_VALIDATE_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.Lint.ValidateArgs,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
