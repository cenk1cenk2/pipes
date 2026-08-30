package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var _ = Describe("Go workspace", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		// Only the resolved tool is rebuilt, so what the previous run decided about
		// the workspace is still there for the run under test to overwrite.
		*P = pipe
		C.Ctx = &tool.Ctx{Cwd: "projects/api", Env: map[string]string{}}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(GoWorkspace(tl).Job())
						})
				},
			},
		}).Run()
	}

	gowork := func(stdout string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{Name: "go", Args: []string{"env", "GOWORK"}, Stdout: stdout}
	}

	It("takes the workspace from the flag without asking the toolchain", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Workspace: true})).To(Succeed())

		Expect(C.Workspace).To(BeTrue())
		Expect(runner.Invocations()).To(BeEmpty())
	})

	It("falls back to the workspace the toolchain reports", func() {
		runner := fixtures.Runner(gowork("/repository/go.work\n"))

		Expect(run(runner, Pipe{})).To(Succeed())

		Expect(C.Workspace).To(BeTrue())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("go"))
		Expect(invocation.Args).To(Equal([]string{"env", "GOWORK"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
	})

	// go env reports "off" instead of an empty value when workspace mode is
	// explicitly disabled, which is not a workspace to vendor or to enumerate.
	It("stays on the single module when the toolchain reports workspace mode off", func() {
		runner := fixtures.Runner(gowork("off\n"))

		Expect(run(runner, Pipe{})).To(Succeed())

		Expect(C.Workspace).To(BeFalse())
	})

	It("stays on the single module when the toolchain reports nothing", func() {
		runner := fixtures.Runner(gowork(""))

		Expect(run(runner, Pipe{})).To(Succeed())

		Expect(C.Workspace).To(BeFalse())
	})

	It("forgets a workspace that the current run did not detect", func() {
		Expect(run(fixtures.Runner(gowork("/repository/go.work\n")), Pipe{})).To(Succeed())
		Expect(C.Workspace).To(BeTrue())

		Expect(run(fixtures.Runner(gowork("off\n")), Pipe{})).To(Succeed())
		Expect(C.Workspace).To(BeFalse())
	})
})
