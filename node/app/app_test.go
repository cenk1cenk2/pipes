package app_test

import (
	"os"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/node/app"
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

	return fixture.RunCli(append([]string{"pipe-node"}, args...)...)
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

func commandNamed(command *ucli.Command, name string) *ucli.Command {
	GinkgoHelper()

	for _, sub := range command.Commands {
		if sub.Name == name {
			return sub
		}
	}

	return nil
}

func environmentEnable(flags []ucli.Flag) *ucli.BoolFlag {
	GinkgoHelper()

	for _, flag := range flags {
		if converted, ok := flag.(*ucli.BoolFlag); ok && converted.Name == "environment.enable" {
			return converted
		}
	}

	return nil
}

// The flags are package level values shared across the whole tree, so a value a
// spec parsed stays parsed for the ones after it. The specs assert on the
// commands they are about rather than on the whole invocation list, so what an
// earlier one left behind does not decide whether a later one passes.
var _ = Describe("New", func() {
	It("writes the npmrc of the registries it was given on login", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(
			runner,
			"login",
			"--npm.login", `[{ "username": "user", "token": "npm-token", "registry": "registry.example.test" }]`,
		)).To(Succeed())

		npmrc, err := os.ReadFile(".npmrc")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(npmrc)).To(ContainSubstring("//registry.example.test/:_authToken=npm-token"))
	})

	It("installs the dependencies through the package manager on install", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(runner, "install")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("i --frozen-lockfile")))
	})

	It("runs the build script through the package manager on build", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(runner, "build", "--node.build_script", "compile")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("run compile")))
	})

	// The environment feature is hidden and on by default everywhere it is owned
	// by the pipe. Here it is opt-in, and app.New is what turns it around.
	It("offers the environment selection as an opt-in", func() {
		command := app.New(nil, "test", app.Options{})

		build := commandNamed(command, "build")
		Expect(build).NotTo(BeNil())

		flag := environmentEnable(build.Flags)
		Expect(flag).NotTo(BeNil())
		Expect(flag.Hidden).To(BeFalse())
		Expect(flag.Value).To(BeFalse())
	})

	// The script and its arguments are positional, which is what makes this the
	// only command in the tree that takes any.
	It("runs the script it was given on run", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(runner, "run", "migrate", "up")).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("run migrate")))
	})
})
