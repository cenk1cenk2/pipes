package plan

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
	"gitlab.kilic.dev/devops/pipes/terraform/state"
)

func deps(name, cwd string) Deps {
	return Deps{
		Tool:  &tool.Ctx{Cwd: cwd, Env: map[string]string{}},
		State: &state.Pipe{State: state.State{Name: name}},
	}
}

var _ = Describe("Terraform report discriminators", func() {
	It("holds nothing back when neither value has moved off its default", func() {
		Expect(terraformReportDiscriminators(deps("default", "."))).To(BeEmpty())
	})

	It("carries a named state", func() {
		Expect(terraformReportDiscriminators(deps("production", "."))).To(Equal([]string{"production"}))
	})

	It("carries a root module other than the working directory", func() {
		Expect(terraformReportDiscriminators(deps("default", ".deploy/sun"))).To(Equal([]string{".deploy/sun"}))
	})

	It("carries both, state first", func() {
		Expect(terraformReportDiscriminators(deps("production", ".deploy/sun"))).
			To(Equal([]string{"production", ".deploy/sun"}))
	})

	It("reports the state name only once it has been named", func() {
		Expect(terraformStateName(deps("default", "."))).To(BeEmpty())
		Expect(terraformStateName(deps("production", "."))).To(Equal("production"))
	})
})

var _ = Describe("Terraform plan", func() {
	// Only the plan task runs, since the report tasks that follow it would reach for
	// a plan file the stubbed command never wrote.
	run := func(runner *tests.TestingCommandRunner, environment map[string]string) error {
		GinkgoHelper()

		deps := deps("production", ".deploy/sun")

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-terraform",
			CommandName: "plan",
			Flags:       Flags,
			Environment: environment,
			WithoutEnvironment: []string{
				"CI_PIPELINE_SOURCE",
				"TF_PLAN_ARGS",
				"TF_PLAN_CACHE",
				"TF_PLAN_OUTPUT",
				"TF_PLAN_PREVIEW_FOR_MRS",
			},
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(TerraformPlan(tl, deps).Job())
						})
				},
			},
		}).Run()
	}

	It("drops the state lock on a merge request pipeline", func() {
		runner := fixtures.Runner()

		Expect(run(runner, map[string]string{"CI_PIPELINE_SOURCE": "merge_request_event"})).To(Succeed())

		Expect(runner.InvocationNames()).To(ContainElement("terraform"))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("terraform"))
		Expect(invocation.Args).To(Equal([]string{"plan", "-input=false", "-out=plan", "-lock=false"}))
		Expect(invocation.Dir).To(Equal(".deploy/sun"))
	})

	It("keeps the state lock on every other pipeline", func() {
		runner := fixtures.Runner()

		Expect(run(runner, nil)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"plan", "-input=false", "-out=plan"}))
	})
})
