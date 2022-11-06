package setup

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	CATEGORY_ENVIRONMENT = "Environment"
)

var Flags = TL.Plumber.AppendFlags(flags.NewGitFlags(
	flags.GitFlagsSetup{
		GitBranchDestination: &TL.Pipe.Git.Branch,
		GitTagDestination:    &TL.Pipe.Git.Tag,
	},
), []cli.Flag{

	// CATEGORY_ENVIRONMENT

	&cli.StringFlag{
		Category: CATEGORY_ENVIRONMENT,
		Name:     "environment.conditions",
		Usage:    `Regex pattern to select an environment. Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. json([]struct{ match: RegExp, environment: string })`,
		Required: false,
		EnvVars:  []string{"ENVIRONMENT_CONDITIONS"},
		Value:    flags.FLAG_DEFAULT_ENVIRONMENT_CONDITIONS,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_ENVIRONMENT,
		Name:        "environment.fail-on-no-reference",
		Usage:       "Whether to fail on missing environment references.",
		Required:    false,
		EnvVars:     []string{"ENVIRONMENT_FAIL_ON_NO_REFERENCE"},
		Value:       true,
		Destination: &TL.Pipe.Environment.FailOnNoReference,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_ENVIRONMENT,
		Name:        "environment.strict",
		Usage:       "Whether to fail on missing environment selection.",
		Required:    false,
		EnvVars:     []string{"ENVIRONMENT_STRICT"},
		Value:       false,
		Destination: &TL.Pipe.Environment.Strict,
	},
})

func ProcessFlags(tl *TaskList[Pipe]) Job {
	return tl.CreateBasicJob(func() error {
		// setup selection of environment conditions
		if err := json.Unmarshal([]byte(tl.CliContext.String("environment.conditions")), &TL.Pipe.Conditions); err != nil {
			return fmt.Errorf("Can not unmarshal environment conditions: %w", err)
		}

		return tl.Validate(&tl.Pipe)
	})
}
