package cli

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"
)

// Step is a pipe stage as the command that runs it sees it: the flags it needs
// parsed and the task list it contributes. A command is the ordered composition
// of its steps, so a stage shared between commands is declared once.
type Step struct {
	Flags     []ucli.Flag
	Arguments []ucli.Argument
	New       func(p *plumber.Plumber) *plumber.TaskList
}

// Command builds a subcommand that runs every step in the order it was given.
func Command(p *plumber.Plumber, name, description string, steps ...Step) *ucli.Command {
	return &ucli.Command{
		Name:        name,
		Description: description,
		Flags:       stepFlags(steps),
		Arguments:   stepArguments(steps),
		Action:      stepAction(p, steps),
	}
}

// App builds the root command of a pipe that dispatches to subcommands.
func App(name, description, version string, commands ...*ucli.Command) *ucli.Command {
	return &ucli.Command{
		Name:        name,
		Version:     version,
		Usage:       description,
		Description: description,
		Commands:    commands,
	}
}

// Root is App for the single-action pipes, where the steps run from the root
// command because there is no second thing the pipe could be asked to do.
func Root(p *plumber.Plumber, name, description, version string, steps ...Step) *ucli.Command {
	command := Command(p, name, description, steps...)
	command.Version = version
	command.Usage = description

	return command
}

func stepFlags(steps []Step) []ucli.Flag {
	flags := []ucli.Flag{}

	for _, step := range steps {
		flags = append(flags, step.Flags...)
	}

	return flags
}

func stepArguments(steps []Step) []ucli.Argument {
	arguments := []ucli.Argument{}

	for _, step := range steps {
		arguments = append(arguments, step.Arguments...)
	}

	return arguments
}

func stepAction(p *plumber.Plumber, steps []Step) ucli.ActionFunc {
	return func(_ context.Context, _ *ucli.Command) error {
		// The task lists are built inside the action rather than alongside the
		// flags, since a step reads the parsed flag values as it constructs.
		lists := make([]*plumber.TaskList, 0, len(steps))

		for _, step := range steps {
			lists = append(lists, step.New(p))
		}

		return p.RunJobs(plumber.CombineTaskLists(lists...))
	}
}
