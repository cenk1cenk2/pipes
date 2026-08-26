package run

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
)

var _ = Describe("Node run", func() {
	deps := Deps{
		Node: &node.Ctx{PackageManager: node.PackageManager{
			Exe:      "pnpm",
			Commands: node.PackageManagers["pnpm"],
		}},
		Environment: &environment.Ctx{
			Environment: "production",
			EnvVars:     map[string]string{"API_URL": "https://api.example.com"},
		},
	}

	// The whole task list runs rather than the task alone, since what the script
	// resolves to is decided before the tasks are built. The arguments are
	// registered because the command line is what fills them.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, args ...string) error {
		GinkgoHelper()

		*P = pipe
		*C = Ctx{}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "run",
			Args:        append([]string{"pipe-node", "run"}, args...),
			Arguments:   Arguments,
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					return New(p, deps)
				},
			},
		}).Run()
	}

	It("exposes the arguments through its step so the command registers them", func() {
		Expect(Step(deps).Arguments).To(Equal(Arguments))
		Expect(Step(deps).Flags).To(Equal(Flags))
	})

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
