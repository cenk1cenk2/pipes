package pipe

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_environment = "Environment"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    category_environment,
		Name:        "environment.conditions",
		Usage:       "Regular expression patterns to match for selecting the environment. format(json({ tags: { [name: string]: RegExp }, branches: { [name: string]: RegExp } }))",
		Required:    false,
		EnvVars:     []string{"ENVIRONMENT_CONDITIONS"},
		Value:       `{ "production": "^v\\d*\\.\\d*\\.\\d*$", "stage": "^v\\d*\\.\\d*\\.\\d*-.*$" }`,
		Destination: &TL.Pipe.Environment.Conditions,
	},
}
