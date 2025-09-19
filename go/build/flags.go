package build

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v2"
)

//revive:disable:line-length-limit

const (
	CATEGORY_BUILD = "Build"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_ARGS"),
		),
		Usage:       "Arguments to append to build command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},

	&cli.StringFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.output",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_OUTPUT"),
		),
		Usage:       "Output location for the build artifacts.",
		Required:    false,
		Value:       "./dist/",
		Destination: &P.Output,
	},

	&cli.StringFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.binary-name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_BINARY_NAME"),
		),
		Usage:       "Name of the binary to output during build.",
		Required:    false,
		Value:       "bin",
		Destination: &P.BinaryName,
	},

	&cli.StringFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.binary-template",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_BINARY_TEMPLATE"),
		),
		Usage:       "Binary naming for the build artifact. format(Template(map[string]))",
		Required:    false,
		Value:       "{{ .name }}{{ if .os }}-{{ .os }}{{ end }}{{ if .arch }}-{{ .arch }}{{ end }}",
		Destination: &P.BinaryTemplate,
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.ld-flags",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_LD_FLAGS"),
		),
		Usage:       "Arguments for the linker during the build process.",
		Required:    false,
		Value:       []string{},
		Destination: &P.LdFlags,
	},

	&cli.BoolFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.enable-cgo",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_ENABLE_CGO"),
			cli.EnvVar("CGO_ENABLED"),
		),
		Usage:       "Enable CGO during the build process.",
		Required:    false,
		Value:       false,
		Destination: &P.EnableCGO,
	},

	&cli.StringFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.targets",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_TARGETS"),
		),
		Usage:            "Build targets for the build process. format(yaml([]struct{ os: string?, arch: string? }))",
		Required:         false,
		Value:            `[]`,
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.BuildTargets); err != nil {
				return fmt.Errorf("Cannot unmarshal build targets: %w", err)
			}

			return nil
		},
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.tags",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_TAGS"),
		),
		Usage:       "Build tags for the build process.",
		Required:    false,
		Value:       []string{},
		Destination: &P.BuildTags,
	},

	&cli.StringFlag{
		Category: CATEGORY_BUILD,
		Name:     "go.build.variables",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_BUILD_VARIABLES"),
		),
		Usage:            "Build variables for the build process. format(yaml(map[string]string))",
		Required:         false,
		Value:            `{}`,
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.BuildVariables); err != nil {
				return fmt.Errorf("Cannot unmarshal build variables: %w", err)
			}

			return nil
		},
	},
}
