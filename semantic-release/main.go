package main

import (
	"github.com/urfave/cli/v2"

	pipe "github.com/cenk1cenk2/ci-cd-pipes/semantic-release/pipe"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

func main() {
	utils.CliCreate(
		&cli.App{
			Name:    pipe.CLI_NAME,
			Version: pipe.VERSION,
			Action:  run,
			Flags:   pipe.Flags,
		},
	)
}

func run(c *cli.Context) error {
	utils.CliGreet(c)

	return pipe.Pipe.Exec()
}
