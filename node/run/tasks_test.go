package run

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

// The tasks read the package manager and the environment of the pipe around
// them off their package level instances, so a spec seeds those the same way it
// seeds its own.
func seed(packageManager string) {
	*setup.NodeCtx = node.Ctx{PackageManager: node.PackageManager{
		Exe:      packageManager,
		Commands: node.PackageManagers[packageManager],
	}}
	*setup.EnvironmentCtx = environment.Ctx{
		Environment: "production",
		EnvVars:     map[string]string{"API_URL": "https://api.example.com"},
	}
}

var _ = Describe("Node run", func() {
	// The whole task list runs rather than the task alone, since what the script
	// resolves to is decided before the tasks are built. The arguments are
	// registered because the command line is what fills them.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, args ...string) error {
		GinkgoHelper()

		*P = pipe
		*C = Ctx{}
		seed("pnpm")

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "run",
			Args:        append([]string{"pipe-node", "run"}, args...),
			Arguments:   Arguments,
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	It("runs the arguments it was given as the script and its arguments", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Run: Run{Cwd: "."}}, "tsc", "src")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pnpm"))
		Expect(invocation.Args).To(Equal([]string{"run", "tsc", "src"}))
		Expect(invocation.Dir).To(Equal("."))
	})

	// The script flag is one string holding both halves, so the pipe has to cut it
	// where the arguments would otherwise have arrived already separated.
	It("cuts the script flag into the script and its arguments", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Run: Run{Script: "lint --fix", Cwd: "."}})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"run", "lint", "--fix"}))
	})

	// The cut happens before the templating, so the script half of the flag has to
	// be written without a space in it for the two to survive as one action each.
	It("templates the script against the selected environment", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{
			Run: Run{Script: "deploy:{{.Environment}} {{ index .EnvVars \"API_URL\" }}", Cwd: "."},
		})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"run", "deploy:production", "https://api.example.com"}))
	})
})
