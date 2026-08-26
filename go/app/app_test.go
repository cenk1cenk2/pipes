package app_test

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/app"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
)

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(runner *tests.TestingCommandRunner, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return app.New(p, "test", app.Options{})
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(append([]string{"pipe-go"}, args...)...)
}

func formatted(runner *tests.TestingCommandRunner) []string {
	GinkgoHelper()

	invocations := runner.Invocations()
	commands := make([]string, len(invocations))

	for i, invocation := range invocations {
		commands[i] = invocation.Formatted
	}

	return commands
}

var _ = Describe("New", func() {
	It("vendors the modules on install", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "install")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("mod vendor")))
	})

	It("builds the application on build", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "build")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("build -mod=vendor -v")))
	})

	It("runs golangci-lint on lint", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "lint")).To(Succeed())

		Expect(runner.InvocationNames()).To(ContainElement("golangci-lint"))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("run -v --timeout")))
	})

	// The tool and its arguments are positional, which is the only command in the
	// tree that takes any.
	It("runs the named tool on tool", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "tool", "mockery", "generate")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("tool mockery generate")))
	})
})
