package install

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var _ = Describe("Go install", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, workspace bool) error {
		GinkgoHelper()

		*P = pipe
		deps := Deps{Tool: &setup.Ctx{
			Ctx:       &tool.Ctx{Cwd: "projects/api", Env: map[string]string{"GOPATH": "/cache"}},
			Workspace: workspace,
		}}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(
								GoModVendor(tl, deps).Job(),
								GoModVerify(tl, deps).Job(),
							)
						})
				},
			},
		}).Run()
	}

	It("vendors the single module in the working directory", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{}, false)).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Name).To(Equal("go"))
		Expect(invocations[0].Args).To(Equal([]string{"mod", "vendor"}))
		Expect(invocations[0].Dir).To(Equal("projects/api"))
		Expect(invocations[0].Env).To(ContainElement("GOPATH=/cache"))
	})

	It("vendors the whole workspace when the setup resolved one", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{}, true)).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Args).To(Equal([]string{"work", "vendor"}))
	})

	It("appends the arguments to whichever vendor command it picked", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Args: "-e"}, true)).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Args).To(Equal([]string{"work", "vendor", "-e"}))
	})

	It("verifies the modules unless the pipeline turned it off", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Verify: true}, false)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"mod", "verify"}))
		Expect(invocation.Dir).To(Equal("projects/api"))

		disabled := fixtures.Runner()
		Expect(run(disabled, Pipe{}, false)).To(Succeed())

		Expect(disabled.Invocations()).To(HaveLen(1))
	})
})
