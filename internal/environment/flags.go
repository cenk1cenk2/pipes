package environment

import (
	"strings"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/git"
)

//revive:disable:line-length-limit

const CATEGORY_ENVIRONMENT = "Environment"

// Options is what the environment flags are built onto.
type Options struct {
	Destination *Config
}

// Config is the environment selection a pipe was configured with.
type Config struct {
	Enable            bool
	Conditions        []Condition
	FailOnNoReference bool
	Strict            bool
	Git               git.Refs
}

// NewFlags builds the environment flags onto cfg. The git flags come first,
// since the references they carry are what the conditions match against. The
// enable flag ships hidden and on; a pipe that only injects an environment on
// request unhides it and flips the default with OverwriteCliFlag.
func NewFlags(opts Options) []cli.Flag {
	return append(git.NewFlags(git.Options{Destination: &opts.Destination.Git}), []cli.Flag{
		&cli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.enable",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("ENVIRONMENT_ENABLE")),
			Usage:       "Enable environment injection.",
			Required:    false,
			Hidden:      true,
			Value:       true,
			Destination: &opts.Destination.Enable,
		},

		flags.JSONFlag(&opts.Destination.Conditions, &cli.StringFlag{
			Category: CATEGORY_ENVIRONMENT,
			Name:     "environment.conditions",
			Sources:  cli.NewValueSourceChain(cli.EnvVar("ENVIRONMENT_CONDITIONS")),
			Usage: strings.TrimSpace(`
Regex pattern to select an environment.
Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.

format(json([]struct{ match: RegExp, environment: string }))
`),
			Required: false,
			Value:    DEFAULT_CONDITIONS,
		}),

		&cli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.fail-on-no-reference",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("ENVIRONMENT_FAIL_ON_NO_REFERENCE")),
			Usage:       "Fail on missing environment references.",
			Required:    false,
			Value:       true,
			Destination: &opts.Destination.FailOnNoReference,
		},

		&cli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.strict",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("ENVIRONMENT_STRICT")),
			Usage:       "Fail when no environment is selected.",
			Required:    false,
			Value:       true,
			Destination: &opts.Destination.Strict,
		},
	}...)
}
