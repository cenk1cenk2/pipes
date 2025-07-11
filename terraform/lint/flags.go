package lint

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name: "terraform-lint.fmt-check.enable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_LINT_FMT_CHECK_ENABLE"),
		),
		Usage:       "Enable terraform fmt.",
		Required:    false,
		Value:       true,
		Destination: &P.Lint.FormatCheckEnable,
	},
	&cli.StringFlag{
		Name: "terraform-lint.fmt-check.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_LINT_FMT_CHECK_ARGS"),
		),
		Usage:       "Additional arguments for terraform fmt.",
		Required:    false,
		Value:       "",
		Destination: &P.Lint.FormatCheckArgs,
	},
	&cli.BoolFlag{
		Name: "terraform-lint.validate.enable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_LINT_VALIDATE_ENABLE"),
		),
		Usage:       "Enable terraform validate.",
		Required:    false,
		Value:       true,
		Destination: &P.Lint.ValidateEnable,
	},
	&cli.StringFlag{
		Name: "terraform-lint.validate.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_LINT_VALIDATE_ARGS"),
		),
		Usage:       "Additional arguments for terraform validate.",
		Required:    false,
		Value:       "",
		Destination: &P.Lint.ValidateArgs,
	},
}
