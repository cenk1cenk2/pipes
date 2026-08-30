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

// answer lets the runner reply to one invocation of the given command. Seeding
// any response makes the runner strict, so a spec that stubs one command has to
// account for every other command the run makes.
func answer(name string, args ...string) tests.TestingCommandResponse {
	return tests.TestingCommandResponse{Name: name, Args: args}
}

func stdout(out, name string, args ...string) tests.TestingCommandResponse {
	response := answer(name, args...)
	response.Stdout = out

	return response
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

	It("vendors the whole workspace when the pipeline asked for one", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "install", "--go.workspace")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("work vendor")))
	})

	// The explicit flag is what the pipelines set, so the toolchain is never asked
	// where the workspace is when it is on. Seeding a response makes the runner
	// strict, so every command of the run has to be answered for.
	It("vendors the workspace the toolchain reports when the flag is off", func() {
		runner := fixtures.Runner(
			answer("go", "version"),
			stdout("/repository/go.work\n", "go", "env", "GOWORK"),
			answer("go", "work", "vendor"),
			answer("go", "mod", "verify"),
		)

		Expect(run(runner, "install")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("env GOWORK")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("work vendor")))
	})

	It("does not ask the toolchain where the workspace is when the flag is on", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "install", "--go.workspace")).To(Succeed())

		Expect(formatted(runner)).NotTo(ContainElement(ContainSubstring("env GOWORK")))
	})

	It("lints every module of the workspace when the pipeline asked for one", func() {
		runner := fixtures.Runner(
			answer("go", "version"),
			stdout("/repository/api\n/repository/worker\n", "go", "list", "-m", "-f", "{{.Dir}}"),
			answer("golangci-lint"),
		)

		Expect(run(runner, "lint", "--go.workspace")).To(Succeed())

		Expect(formatted(runner)).
			To(ContainElement(ContainSubstring("/repository/api/... /repository/worker/...")))
	})

	// The tool and its arguments are positional, which is the only command in the
	// tree that takes any.
	It("runs the named tool on tool", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "tool", "mockery", "generate")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("tool mockery generate")))
	})
})
