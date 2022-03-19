package main

import (
	"github.com/urfave/cli/v2"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	pipe "gitlab.kilic.dev/bdsm/gitlab-pipes/_template/pipe"
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
