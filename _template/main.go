package main

import (
	"github.com/urfave/cli/v2"

	pipe "gitlab.kilic.dev/devops/gitlab-pipes/_template/pipe"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

func main() {
	utils.CliCreate(
		&cli.App{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Flags:       pipe.Flags,
			Action:      run,
		},
	)
}

func run(c *cli.Context) error {
	utils.CliGreet(c)

	return pipe.Pipe.Exec()
}
