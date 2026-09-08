package lint

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Go lint", func() {
	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, modules ...string) error {
		GinkgoHelper()

		*P = Pipe{Timeout: 5 * time.Minute}
		*setup.C = setup.Ctx{
			Cwd:     "projects/api",
			Env:     map[string]string{"GOPATH": "/cache"},
			Modules: modules,
		}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "lint",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(lint(tl).Job())
						})
				},
			},
		}).Run()
	}

	It("lints the working directory when the setup resolved no workspace", func() {
		runner := fixtures.Runner()

		Expect(run(runner)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"golangci-lint"}))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"run", "-v", "--timeout", "5m0s"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
		Expect(invocation.Env).To(ContainElement("GOPATH=/cache"))
	})

	// every module is linted from inside its own directory rather than through a
	// "<module>/..." pattern, since the go tool resolves such a pattern to nothing
	// for a module whose directory name starts with an underscore.
	It("lints every module of the workspace from inside it", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "/repository/_template", "/repository/api")).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).To(HaveLen(2))

		Expect(invocations[0].Args).To(Equal([]string{"run", "-v", "--timeout", "5m0s", "./..."}))
		Expect(invocations[0].Dir).To(Equal("/repository/_template"))

		Expect(invocations[1].Args).To(Equal([]string{"run", "-v", "--timeout", "5m0s", "./..."}))
		Expect(invocations[1].Dir).To(Equal("/repository/api"))
	})

	It("carries the resolved environment into every module", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "/repository/api", "/repository/worker")).To(Succeed())

		for _, invocation := range runner.Invocations() {
			Expect(invocation.Env).To(ContainElement("GOPATH=/cache"))
		}
	})
})
