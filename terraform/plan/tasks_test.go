package plan

import (
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	"gitlab.kilic.dev/devops/pipes/terraform/state"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

// The tasks read the setup and the state of the pipe around them off their
// package level instances, so a spec seeds those the same way it seeds its own.
func seed(name, cwd string) {
	*setup.C = setup.Ctx{Cwd: cwd, Env: map[string]string{}}
	*state.P = state.Pipe{State: state.State{Name: name}}
}

var _ = Describe("Terraform report discriminators", func() {
	It("holds nothing back when neither value has moved off its default", func() {
		seed("default", ".")

		Expect(terraformReportDiscriminators()).To(BeEmpty())
	})

	It("carries a named state", func() {
		seed("production", ".")

		Expect(terraformReportDiscriminators()).To(Equal([]string{"production"}))
	})

	It("carries a root module other than the working directory", func() {
		seed("default", ".deploy/sun")

		Expect(terraformReportDiscriminators()).To(Equal([]string{".deploy/sun"}))
	})

	It("carries both, state first", func() {
		seed("production", ".deploy/sun")

		Expect(terraformReportDiscriminators()).To(Equal([]string{"production", ".deploy/sun"}))
	})

	It("reports the state name only once it has been named", func() {
		seed("default", ".")
		Expect(terraformStateName()).To(BeEmpty())

		seed("production", ".")
		Expect(terraformStateName()).To(Equal("production"))
	})
})

var _ = Describe("Terraform plan", func() {
	// Only the plan task runs, since the report tasks that follow it would reach for
	// a plan file the stubbed command never wrote.
	run := func(runner *tests.TestingCommandRunner, environment map[string]string) error {
		GinkgoHelper()

		seed("production", ".deploy/sun")

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-terraform",
			CommandName: "plan",
			Flags:       Flags,
			Environment: environment,
			WithoutEnvironment: []string{
				"CI_PIPELINE_SOURCE",
				"TERRAFORM_PLAN_ARGS",
				"TERRAFORM_PLAN_OUTPUT",
				"TERRAFORM_PLAN_PIPELINE_SOURCE",
				"TERRAFORM_PLAN_PREVIEW_FOR_MERGE_REQUESTS",
			},
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(TerraformPlan(tl).Job())
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

var _ = Describe("Terraform plan cleanup", func() {
	// The pipe is seeded directly rather than through the flags, since a package
	// level flag only reads its environment on the first parse of the process and
	// the specs would otherwise decide each other's pipeline source.
	run := func(cwd string, pipe Pipe) error {
		GinkgoHelper()

		*P = pipe
		seed("production", cwd)

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-terraform",
			CommandName: "plan",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(TerraformPlanCleanup(tl).Job())
						})
				},
			},
		}).Run()
	}

	planned := func() string {
		GinkgoHelper()

		cwd := tests.TempDir()
		Expect(os.WriteFile(filepath.Join(cwd, "plan"), []byte("plan"), 0600)).To(Succeed())

		return cwd
	}

	preview := func(source string) Pipe {
		return Pipe{Plan: Plan{Output: "plan", PipelineSource: source, PreviewForMergeRequests: true}}
	}

	// A preview plan is taken without the state lock, so it is not the plan that
	// gets applied and must not outlive the job that wrote it.
	It("removes the plan a merge request pipeline previewed", func() {
		cwd := planned()

		Expect(run(cwd, preview("merge_request_event"))).To(Succeed())

		Expect(filepath.Join(cwd, "plan")).NotTo(BeAnExistingFile())
	})

	It("keeps the plan every other pipeline wrote to be applied", func() {
		cwd := planned()

		Expect(run(cwd, preview("push"))).To(Succeed())

		Expect(filepath.Join(cwd, "plan")).To(BeAnExistingFile())
	})

	It("keeps the plan when previews are turned off", func() {
		cwd := planned()

		pipe := preview("merge_request_event")
		pipe.Plan.PreviewForMergeRequests = false

		Expect(run(cwd, pipe)).To(Succeed())

		Expect(filepath.Join(cwd, "plan")).To(BeAnExistingFile())
	})
})
