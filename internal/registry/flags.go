package registry

import (
	"fmt"
	"strings"

	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

const DEFAULT_URI = "docker.io"

// Credentials is the registry a pipe logs in to.
type Credentials struct {
	Uri      string
	Username string
	Password string
}

// Spec is how one pipe names its registry flags. The category and the label are
// per pipe because they are what the generated documentation prints as the
// section heading and the description.
type Spec struct {
	Category string
	// Label reads as the subject of the usage sentence, as in "Helm registry url
	// to login to".
	Label   string
	Prefix  string
	Command string
	// LegacyEnv is the prefix of the environment names the flags answered to
	// before they were unified. They are kept forever and listed first.
	LegacyEnv string
}

func NewFlags(spec Spec, dst *Credentials) []ucli.Flag {
	return []ucli.Flag{
		&ucli.StringFlag{
			Category:    spec.Category,
			Name:        spec.name("uri"),
			Sources:     spec.envs("URI"),
			Usage:       fmt.Sprintf("%s url to login to.", spec.Label),
			Required:    false,
			Value:       DEFAULT_URI,
			Destination: &dst.Uri,
		},

		&ucli.StringFlag{
			Category:    spec.Category,
			Name:        spec.name("username"),
			Sources:     spec.envs("USERNAME"),
			Usage:       fmt.Sprintf("%s username for the given registry.", spec.Label),
			Required:    false,
			Destination: &dst.Username,
		},

		&ucli.StringFlag{
			Category:    spec.Category,
			Name:        spec.name("password"),
			Sources:     spec.envs("PASSWORD"),
			Usage:       fmt.Sprintf("%s password for the given registry.", spec.Label),
			Required:    false,
			Destination: cli.MarkSecret(&dst.Password),
		},
	}
}

func (s Spec) name(key string) string {
	return fmt.Sprintf("%s.%s.registry.%s", s.Prefix, s.Command, key)
}

func (s Spec) envs(key string) ucli.ValueSourceChain {
	return cli.EnvVars(
		fmt.Sprintf("%s_REGISTRY_%s", s.LegacyEnv, key),
		strings.ToUpper(fmt.Sprintf("%s_%s_REGISTRY_%s", s.Prefix, s.Command, key)),
	)
}
