package main

import (
	"github.com/urfave/cli/v2"

	login "gitlab.kilic.dev/devops/pipes/node/login"
	pipe "gitlab.kilic.dev/devops/pipes/semantic-release/pipe"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

func main() {
	utils.CliCreate(
		&cli.App{
			Name:    pipe.CLI_NAME,
			Version: pipe.VERSION,
			Action:  run,
			Flags:   append(login.Flags, pipe.Flags...),
		},
	)
}

func run(c *cli.Context) error {
	utils.CliGreet(c)

	return pipe.Pipe.Exec()
}
