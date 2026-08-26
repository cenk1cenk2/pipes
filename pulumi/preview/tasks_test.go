package preview

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
)

func deps(name, cwd string) Deps {
	return Deps{
		Tool:  &tool.Ctx{Cwd: cwd, Env: map[string]string{}},
		Stack: &stack.Pipe{Stack: name},
	}
}

var _ = Describe("Pulumi report discriminators", func() {
	It("carries the stack on its own from the working directory", func() {
		Expect(pulumiReportDiscriminators(deps("main", "."))).To(Equal([]string{"main"}))
	})

	It("carries a project directory alongside the stack", func() {
		Expect(pulumiReportDiscriminators(deps("main", "projects/warehouse"))).
			To(Equal([]string{"main", "projects/warehouse"}))
	})
})

var _ = Describe("Pulumi preview", func() {
	// Only the preview task runs, since the report tasks that follow it would reach
	// for a plan file the stubbed command never wrote.
	run := func(runner *tests.TestingCommandRunner, environment map[string]string) error {
		GinkgoHelper()

		deps := deps("main", "projects/warehouse")

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-pulumi",
			CommandName: "preview",
			Flags:       Flags,
			Environment: environment,
			WithoutEnvironment: []string{
				"PULUMI_PLAN",
				"PULUMI_SUMMARY_OUTPUT",
			},
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(PulumiPlan(tl, deps).Job())
						})
				},
			},
		}).Run()
	}

	It("saves the plan the report is read back out of, in the project directory", func() {
		runner := fixtures.Runner()

		Expect(run(runner, nil)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pulumi"))
		Expect(invocation.Args).
			To(Equal([]string{"preview", "--non-interactive", "--diff", "--save-plan", "plan.json"}))
		Expect(invocation.Dir).To(Equal("projects/warehouse"))
	})

	It("saves the plan where the flag pointed it", func() {
		runner := fixtures.Runner()

		Expect(run(runner, map[string]string{"PULUMI_PLAN": "preview.json"})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).
			To(Equal([]string{"preview", "--non-interactive", "--diff", "--save-plan", "preview.json"}))
	})
})
