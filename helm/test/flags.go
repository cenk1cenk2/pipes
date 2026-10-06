package test

import (
	"github.com/urfave/cli/v3"
)

const CategoryHelmTest = "Helm Test"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Category: CategoryHelmTest,
		Name:     "helm.test.files",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_TEST_FILES"),
		),
		Usage:       "Test suite files relative to the chart, helm-unittest reads tests/*_test.yaml when none is given. format(glob)",
		Required:    false,
		Destination: &P.Files,
	},

	&cli.StringSliceFlag{
		Category: CategoryHelmTest,
		Name:     "helm.test.values",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_TEST_VALUES"),
		),
		Usage:       "Values files overriding the chart values for every test suite. format(glob)",
		Required:    false,
		Destination: &P.Values,
	},

	&cli.BoolFlag{
		Category: CategoryHelmTest,
		Name:     "helm.test.strict",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_TEST_STRICT"),
		),
		Usage:       "Parse the test suites strictly, failing on unknown fields.",
		Required:    false,
		Value:       false,
		Destination: &P.Strict,
	},

	&cli.StringFlag{
		Category: CategoryHelmTest,
		Name:     "helm.test.output.file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_TEST_OUTPUT_FILE"),
		),
		Usage:       "File to write the test results to. Leave empty to write no report.",
		Required:    false,
		Value:       "",
		Destination: &P.Output.File,
	},

	&cli.StringFlag{
		Category: CategoryHelmTest,
		Name:     "helm.test.output.type",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_TEST_OUTPUT_TYPE"),
		),
		Usage:       `Format of the test results file. format(enum("JUnit", "NUnit", "XUnit", "Sonar"))`,
		Required:    false,
		Value:       "JUnit",
		Destination: &P.Output.Type,
	},
}
