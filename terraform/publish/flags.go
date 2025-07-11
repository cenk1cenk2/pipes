package publish

import (
	"context"
	"regexp"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

const (
	CATEGORY_MODULE          = "Module"
	CATEGORY_REGISTRY        = "Registry"
	CATEGORY_REGISTRY_GITLAB = "Registry - Gitlab"
)

var Flags = CombineFlags(flags.NewTagsFileFlags(
	flags.TagsFileFlagsSetup{
		TagsFileDestination: &P.Module.TagsFile,
		TagsFileRequired:    false,
		TagsFileValue:       ".tags",
	},
), []cli.Flag{
	// CATEGORY_MODULE

	&cli.StringFlag{
		Category: CATEGORY_MODULE,
		Name:     "terraform-module.name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_MODULE_NAME"),
			cli.EnvVar("CI_PROJECT_NAME"),
		),
		Usage:       "Name for the module that will be published.",
		Required:    true,
		Value:       "",
		Destination: &P.Module.Name,
		Action: func(ctx context.Context, c *cli.Command, s string) error {
			regexp.MustCompile(`[_ ]`).ReplaceAllString(s, "-")

			P.Module.Name = s

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_MODULE,
		Name:     "terraform-module.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_MODULE_CWD"),
			cli.EnvVar("TF_ROOT"),
		),
		Usage:       "Directory for the module that will be published.",
		Required:    false,
		Value:       ".",
		Destination: &P.Module.Cwd,
	},

	&cli.StringFlag{
		Category: CATEGORY_MODULE,
		Name:     "terraform-module.system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_MODULE_SYSTEM"),
		),
		Usage:       "Module system for the module that will be published.",
		Required:    false,
		Value:       "local",
		Destination: &P.Module.System,
	},

	// CATEGORY_REGISTRY

	&cli.StringFlag{
		Category: CATEGORY_REGISTRY,
		Name:     "terraform-module.registry",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_MODULE_REGISTRY"),
		),
		Usage:       "Registry of the module that will be published.",
		Required:    false,
		Value:       TF_REGISTRY_GITLAB,
		Destination: &P.Registry.Name,
	},

	// CATEGORY_REGISTRY_GITLAB

	&cli.StringFlag{
		Category: CATEGORY_REGISTRY_GITLAB,
		Name:     "terraform-module.registry.gitlab.api-url",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CI_API_V4_URL"),
		),
		Usage:       "Gitlab API URL for publish call.",
		Required:    false,
		Destination: &P.Registry.Gitlab.ApiUrl,
	},

	&cli.StringFlag{
		Category: CATEGORY_REGISTRY_GITLAB,
		Name:     "terraform-module.registry.gitlab.project-id",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CI_PROJECT_ID"),
		),
		Usage:       "Gitlab project id for publish call.",
		Required:    false,
		Destination: &P.Registry.Gitlab.ProjectId,
	},

	&cli.StringFlag{
		Category: CATEGORY_REGISTRY_GITLAB,
		Name:     "terraform-module.registry.gitlab.token",
		Usage:    "Gitlab API token for publish call.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CI_JOB_TOKEN"),
		),
		Required:    false,
		Destination: &P.Registry.Gitlab.Token,
	},
})
