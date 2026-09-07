// Package fixtures is the test scaffolding the pipes share: the stubs and the
// loggers a spec needs whichever module it lives in.
package fixtures

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	"github.com/sirupsen/logrus"
)

// Runner builds a command runner that answers the given responses instead of
// executing anything. An invocation nothing matches still records, so a spec can
// assert on what a pipe would have run without seeding a response for it.
func Runner(responses ...tests.TestingCommandResponse) *tests.TestingCommandRunner {
	return tests.NewTestingCommandRunner().AddResponses(responses...)
}

// Cli runs a task list command against the given runner. Every field of the spec
// stays available, since the pipes differ in command name, flags and arguments
// and only agree on wanting their commands stubbed.
func Cli(runner *tests.TestingCommandRunner, spec tests.TaskListCli) *tests.TaskListCliFixture {
	spec.Runtime = plumber.Runtime{CommandRunner: runner.Runner()}

	return tests.NewTaskListCli(spec)
}

// Log is a logger that writes into the spec output at the loudest level, for the
// functions that take one rather than reaching for the task they run under.
func Log() *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(GinkgoWriter)
	logger.SetLevel(logrus.TraceLevel)

	return logrus.NewEntry(logger)
}
