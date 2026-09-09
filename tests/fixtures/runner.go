// Package fixtures is the test scaffolding the pipes share, whichever module a spec lives in.
package fixtures

import (
	"github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
)

// Runner builds a command runner that answers the given responses instead of
// executing anything. An invocation nothing matches still records.
func Runner(responses ...tests.TestingCommandResponse) *tests.TestingCommandRunner {
	return tests.NewTestingCommandRunner().AddResponses(responses...)
}

// Cli runs a task list command against the given runner. Every field of the spec
// stays available, since the pipes only agree on wanting their commands stubbed.
func Cli(runner *tests.TestingCommandRunner, spec tests.TaskListCli) *tests.TaskListCliFixture {
	spec.Runtime = plumber.Runtime{CommandRunner: runner.Runner()}

	return tests.NewTaskListCli(spec)
}

// Task is a task on a plumber of its own that logs into the spec output, for the
// helpers that take the task they run under instead of a bare logger.
func Task(name ...string) *plumber.Task {
	return plumber.NewTaskList(tests.NewPlumber().Plumber).CreateTask(name...)
}
