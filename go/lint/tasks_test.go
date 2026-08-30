package lint

import (
	"time"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var _ = Describe("Go lint", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, workspace bool) error {
		GinkgoHelper()

		*P = Pipe{Timeout: 5 * time.Minute}
		*C = Ctx{}
		deps := Deps{Tool: &setup.Ctx{
			Ctx:       &tool.Ctx{Cwd: "projects/api", Env: map[string]string{"GOPATH": "/cache"}},
			Workspace: workspace,
		}}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "lint",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(GoLint(tl, deps).Job())
						})
				},
			},
		}).Run()
	}

	modules := func(stdout string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{Name: "go", Stdout: stdout}
	}

	It("lints the working directory when the setup resolved no workspace", func() {
		runner := fixtures.Runner()

		Expect(run(runner, false)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"golangci-lint"}))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"run", "-v", "--timeout", "5m0s"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
		Expect(invocation.Env).To(ContainElement("GOPATH=/cache"))
	})

	It("lints every module of the workspace by its own package pattern", func() {
		runner := fixtures.Runner(modules("/repository/api\n/repository/worker\n"))

		Expect(run(runner, true)).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).To(HaveLen(2))
		Expect(invocations[0].Args).To(Equal([]string{"list", "-m", "-f", "{{.Dir}}"}))
		Expect(invocations[0].Dir).To(Equal("projects/api"))
		Expect(invocations[1].Args).
			To(Equal([]string{"run", "-v", "--timeout", "5m0s", "/repository/api/...", "/repository/worker/..."}))
	})
})
