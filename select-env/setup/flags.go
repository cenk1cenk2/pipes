package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

var Flags = CombineFlags(flags.NewGitFlags(
	flags.GitFlagsSetup{
		GitBranchDestination: &P.Git.Branch,
		GitTagDestination:    &P.Git.Tag,
	},
), []cli.Flag{
	&cli.BoolFlag{
		Category: flags.CATEGORY_ENVIRONMENT,
		Name:     "environment.enable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("ENVIRONMENT_ENABLE"),
		),
		Usage:       "Enable environment injection.",
		Required:    false,
		Hidden:      true,
		Value:       true,
		Destination: &P.Environment.Enable,
	},

	// CATEGORY_ENVIRONMENT

	&cli.StringFlag{
		Category: flags.CATEGORY_ENVIRONMENT,
		Name:     "environment.conditions",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("ENVIRONMENT_CONDITIONS"),
		),
		Usage: `Regex pattern to select an environment.
      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.
      json([]struct{ match: RegExp, environment: string })`,
		Required:    false,
		Value:       flags.FLAG_DEFAULT_ENVIRONMENT_CONDITIONS,
		Destination: &raw.EnvironmentConditions,
	},

	&cli.BoolFlag{
		Category: flags.CATEGORY_ENVIRONMENT,
		Name:     "environment.fail-on-no-reference",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("ENVIRONMENT_FAIL_ON_NO_REFERENCE"),
		),
		Usage:       "Fail on missing environment references.",
		Required:    false,
		Value:       true,
		Destination: &P.Environment.FailOnNoReference,
	},

	&cli.BoolFlag{
		Category: flags.CATEGORY_ENVIRONMENT,
		Name:     "environment.strict",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("ENVIRONMENT_STRICT"),
		),
		Usage:       "Fail on no environment selected.",
		Required:    false,
		Value:       true,
		Destination: &P.Environment.Strict,
	},
})
