package environment

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/git"
)

//revive:disable:line-length-limit

const CATEGORY_ENVIRONMENT = "Environment"

// Config is the environment selection a pipe was configured with.
type Config struct {
	Enable            bool
	Conditions        []Condition
	FailOnNoReference bool
	Strict            bool
	Git               git.Refs
}

// NewFlags builds the environment flags onto cfg. The git flags come first
// because the references they carry are what the conditions match against.
//
// The pipes that only inject an environment on request unhide the enable flag
// and flip its default with OverwriteCliFlag, so the flag stays hidden and on
// here for the pipe whose whole job this is.
func NewFlags(cfg *Config) []cli.Flag {
	return append(git.NewFlags(&cfg.Git), []cli.Flag{
		&cli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.enable",
			Sources:     flags.EnvVars("ENVIRONMENT_ENABLE"),
			Usage:       "Enable environment injection.",
			Required:    false,
			Hidden:      true,
			Value:       true,
			Destination: &cfg.Enable,
		},

		flags.JSONFlag(&cli.StringFlag{
			Category: CATEGORY_ENVIRONMENT,
			Name:     "environment.conditions",
			Sources:  flags.EnvVars("ENVIRONMENT_CONDITIONS"),
			Usage: `Regex pattern to select an environment.
      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.
      json([]struct{ match: RegExp, environment: string })`,
			Required: false,
			Value:    DEFAULT_CONDITIONS,
		}, &cfg.Conditions),

		&cli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.fail-on-no-reference",
			Sources:     flags.EnvVars("ENVIRONMENT_FAIL_ON_NO_REFERENCE"),
			Usage:       "Fail on missing environment references.",
			Required:    false,
			Value:       true,
			Destination: &cfg.FailOnNoReference,
		},

		&cli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.strict",
			Sources:     flags.EnvVars("ENVIRONMENT_STRICT"),
			Usage:       "Fail on no environment selected.",
			Required:    false,
			Value:       true,
			Destination: &cfg.Strict,
		},
	}...)
}
