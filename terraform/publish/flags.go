package publish

import (
	"fmt"
	"regexp"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "gitlab.kilic.dev/libraries/plumber/v5"
	"golang.org/x/exp/slices"
)

//revive:disable:line-length-limit

const (
	CATEGORY_MODULE          = "Module"
	CATEGORY_REGISTRY        = "Registry"
	CATEGORY_REGISTRY_GITLAB = "Registry - Gitlab"
)

var Flags = TL.Plumber.AppendFlags(flags.NewTagsFileFlags(
	flags.TagsFileFlagsSetup{
		TagsFileDestination: &TL.Pipe.Module.TagsFile,
		TagsFileRequired:    true,
		TagsFileValue: ".tags",
	},
), []cli.Flag{
	// CATEGORY_MODULE

	&cli.StringFlag{
		Category:    CATEGORY_MODULE,
		Name:        "terraform-module.name",
		Usage:       "Name for the module that will be published.",
		Required:    true,
		EnvVars:     []string{"TF_MODULE_NAME", "CI_PROJECT_NAME"},
		Value:       "",
		Destination: &TL.Pipe.Module.Name,
	},

	&cli.StringFlag{
		Category:    CATEGORY_MODULE,
		Name:        "terraform-module.cwd",
		Usage:       "Directory for the module that will be published.",
		Required:    false,
		EnvVars:     []string{"TF_MODULE_CWD", "TF_ROOT"},
		Value:       ".",
		Destination: &TL.Pipe.Module.Cwd,
	},

	&cli.StringFlag{
		Category:    CATEGORY_MODULE,
		Name:        "terraform-module.system",
		Usage:       "Module system for the module that will be published.",
		Required:    false,
		EnvVars:     []string{"TF_MODULE_SYSTEM"},
		Value:       "local",
		Destination: &TL.Pipe.Module.System,
	},

	// CATEGORY_REGISTRY

	&cli.StringFlag{
		Category:    CATEGORY_REGISTRY,
		Name:        "terraform-module.registry",
		Usage:       "Registry of the module that will be published.",
		Required:    false,
		EnvVars:     []string{"TF_MODULE_REGISTRY"},
		Value:       TF_REGISTRY_GITLAB,
		Destination: &TL.Pipe.Registry.Name,
	},

	// CATEGORY_REGISTRY_GITLAB

	&cli.StringFlag{
		Category:    CATEGORY_REGISTRY_GITLAB,
		Name:        "terraform-module.registry.gitlab.api-url",
		Usage:       "Gitlab API URL for publish call.",
		Required:    false,
		EnvVars:     []string{"CI_API_V4_URL"},
		Destination: &TL.Pipe.Registry.Gitlab.ApiUrl,
	},

	&cli.StringFlag{
		Category:    CATEGORY_REGISTRY_GITLAB,
		Name:        "terraform-module.registry.gitlab.project-id",
		Usage:       "Gitlab project id for publish call.",
		Required:    false,
		EnvVars:     []string{"CI_PROJECT_ID"},
		Destination: &TL.Pipe.Registry.Gitlab.ProjectId,
	},

	&cli.StringFlag{
		Category:    CATEGORY_REGISTRY_GITLAB,
		Name:        "terraform-module.registry.gitlab.token",
		Usage:       "Gitlab API token for publish call.",
		Required:    false,
		EnvVars:     []string{"CI_JOB_TOKEN"},
		Destination: &TL.Pipe.Registry.Gitlab.Token,
	},
})

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	registry := tl.CliContext.String("terraform-module.registry")
	if !slices.Contains([]string{TF_REGISTRY_GITLAB}, registry) {
		return fmt.Errorf("Registry type is not supported: %s", registry)
	}

	cleanse := regexp.MustCompile(`[_ ]`)

	tl.Pipe.Module.Name = cleanse.ReplaceAllString(tl.Pipe.Module.Name, "-")

	return nil
}
