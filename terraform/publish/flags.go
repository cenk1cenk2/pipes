package publish

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"
)

//revive:disable:line-length-limit

const (
	CATEGORY_MODULE          = "Module"
	CATEGORY_REGISTRY        = "Registry"
	CATEGORY_REGISTRY_GITLAB = "Registry - Gitlab"
)

var Flags = CombineFlags(
	tagsfile.NewFlags(&P.Module.TagsFile, ".tags", nil, false),
	[]cli.Flag{
		// CATEGORY_MODULE

		&cli.StringFlag{
			Category: CATEGORY_MODULE,
			Name:     "terraform.publish.module.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TF_MODULE_NAME"),
				cli.EnvVar("CI_PROJECT_NAME"),
				cli.EnvVar("TERRAFORM_PUBLISH_MODULE_NAME"),
			),
			Usage:       "Name for the module that will be published.",
			Required:    true,
			Value:       "",
			Destination: &P.Module.Name,
		},

		&cli.StringFlag{
			Category: CATEGORY_MODULE,
			Name:     "terraform.publish.module.cwd",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TF_MODULE_CWD"),
				cli.EnvVar("TF_ROOT"),
				cli.EnvVar("TERRAFORM_PUBLISH_MODULE_CWD"),
			),
			Usage:       "Directory for the module that will be published.",
			Required:    false,
			Value:       ".",
			Destination: &P.Module.Cwd,
		},

		&cli.StringFlag{
			Category: CATEGORY_MODULE,
			Name:     "terraform.publish.module.system",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TF_MODULE_SYSTEM"),
				cli.EnvVar("TERRAFORM_PUBLISH_MODULE_SYSTEM"),
			),
			Usage:       "Module system for the module that will be published.",
			Required:    false,
			Value:       "local",
			Destination: &P.Module.System,
		},

		// CATEGORY_REGISTRY

		&cli.StringFlag{
			Category: CATEGORY_REGISTRY,
			Name:     "terraform.publish.registry.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TF_MODULE_REGISTRY"),
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_NAME"),
			),
			Usage:       "Registry of the module that will be published.",
			Required:    false,
			Value:       TF_REGISTRY_GITLAB,
			Destination: &P.Registry.Name,
		},

		// CATEGORY_REGISTRY_GITLAB

		&cli.StringFlag{
			Category: CATEGORY_REGISTRY_GITLAB,
			Name:     "terraform.publish.registry.gitlab.api-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_API_V4_URL"),
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_GITLAB_API_URL"),
			),
			Usage:       "Gitlab API URL for publish call.",
			Required:    false,
			Destination: &P.Registry.Gitlab.ApiUrl,
		},

		&cli.StringFlag{
			Category: CATEGORY_REGISTRY_GITLAB,
			Name:     "terraform.publish.registry.gitlab.project-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PROJECT_ID"),
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_GITLAB_PROJECT_ID"),
			),
			Usage:       "Gitlab project id for publish call.",
			Required:    false,
			Destination: &P.Registry.Gitlab.ProjectId,
		},

		&cli.StringFlag{
			Category: CATEGORY_REGISTRY_GITLAB,
			Name:     "terraform.publish.registry.gitlab.token",
			Usage:    "Gitlab API token for publish call.",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_JOB_TOKEN"),
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_GITLAB_TOKEN"),
			),
			Required:    false,
			Destination: &P.Registry.Gitlab.Token,
		},
	})
