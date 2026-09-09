package publish

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"
)

//revive:disable:line-length-limit

const (
	CategoryModule         = "Module"
	CategoryRegistry       = "Registry"
	CategoryRegistryGitLab = "Registry - GitLab"
)

var Flags = CombineFlags(
	tagsfile.NewFlags(tagsfile.Options{Destination: &P.Module.TagsFile, Value: ".tags"}),
	[]cli.Flag{
		// CategoryModule

		&cli.StringFlag{
			Category: CategoryModule,
			Name:     "terraform.publish.module.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_PUBLISH_MODULE_NAME"),
				cli.EnvVar("TF_MODULE_NAME"),
				cli.EnvVar("CI_PROJECT_NAME"),
			),
			Usage:       "Name for the module that will be published.",
			Required:    true,
			Value:       "",
			Destination: &P.Module.Name,
		},

		&cli.StringFlag{
			Category: CategoryModule,
			Name:     "terraform.publish.module.cwd",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_PUBLISH_MODULE_CWD"),
				cli.EnvVar("TF_MODULE_CWD"),
				cli.EnvVar("TF_ROOT"),
			),
			Usage:       "Directory for the module that will be published.",
			Required:    false,
			Value:       ".",
			Destination: &P.Module.Cwd,
		},

		&cli.StringFlag{
			Category: CategoryModule,
			Name:     "terraform.publish.module.system",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_PUBLISH_MODULE_SYSTEM"),
				cli.EnvVar("TF_MODULE_SYSTEM"),
			),
			Usage:       "Module system for the module that will be published.",
			Required:    false,
			Value:       "local",
			Destination: &P.Module.System,
		},

		// CategoryRegistry

		&cli.StringFlag{
			Category: CategoryRegistry,
			Name:     "terraform.publish.registry.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_NAME"),
				cli.EnvVar("TF_MODULE_REGISTRY"),
			),
			Usage:       `Registry of the module that will be published. format(enum("gitlab"))`,
			Required:    false,
			Value:       TFRegistryGitLab,
			Destination: &P.Registry.Name,
		},

		// CategoryRegistryGitLab

		&cli.StringFlag{
			Category: CategoryRegistryGitLab,
			Name:     "terraform.publish.registry.gitlab.api-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_GITLAB_API_URL"),
				cli.EnvVar("CI_API_V4_URL"),
			),
			Usage:       "GitLab API URL for the publish call.",
			Required:    false,
			Destination: &P.Registry.Gitlab.ApiUrl,
		},

		&cli.StringFlag{
			Category: CategoryRegistryGitLab,
			Name:     "terraform.publish.registry.gitlab.project-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_GITLAB_PROJECT_ID"),
				cli.EnvVar("CI_PROJECT_ID"),
			),
			Usage:       "GitLab project id for the publish call.",
			Required:    false,
			Destination: &P.Registry.Gitlab.ProjectId,
		},

		&cli.StringFlag{
			Category: CategoryRegistryGitLab,
			Name:     "terraform.publish.registry.gitlab.token",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TERRAFORM_PUBLISH_REGISTRY_GITLAB_TOKEN"),
				cli.EnvVar("CI_JOB_TOKEN"),
			),
			Usage:       "GitLab API token for the publish call.",
			Required:    false,
			Destination: &P.Registry.Gitlab.Token,
		},
	})
