package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Go install", func() {
	seed := func(pipe Pipe, workspace bool) {
		GinkgoHelper()

		*P = pipe
		*setup.C = setup.Ctx{
			Cwd:       "projects/api",
			Env:       map[string]string{"GOPATH": "/cache"},
			Workspace: workspace,
		}
	}

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, workspace bool) error {
		GinkgoHelper()

		seed(pipe, workspace)

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(
								vendor(tl).Job(),
								verify(tl).Job(),
							)
						})
				},
			},
		}).Run()
	}

	disabled := func(workspace bool, task func(*TaskList) *Task) bool {
		GinkgoHelper()

		seed(Pipe{}, workspace)

		tl := &TaskList{}
		tl.New(NewPlumber(func(_ *Plumber) *cli.Command {
			return &cli.Command{Name: "test"}
		}))

		return task(tl).IsDisabled()
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

		single := fixtures.Runner()

		Expect(run(single, Pipe{Args: "-e"}, false)).To(Succeed())

		invocations = single.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Args).To(Equal([]string{"mod", "vendor", "-e"}))
	})

	// the two vendor tasks are the halves of the workspace condition, so exactly
	// one of them runs under the parent whichever way the setup resolved.
	DescribeTable(
		"vendors through exactly one of the two tasks",
		func(workspace bool, module, workspaceDisabled bool) {
			Expect(disabled(workspace, vendorModule)).To(Equal(module))
			Expect(disabled(workspace, vendorWorkspace)).To(Equal(workspaceDisabled))
		},
		Entry("no workspace", false, false, true),
		Entry("workspace", true, true, false),
	)

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
