package publish

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"
	helmv2 "helm.sh/helm/v4/pkg/chart/v2"

	"gitlab.kilic.dev/devops/pipes/helm/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Helm publish tasks", func() {
	seed := func(cwd, name string) {
		*setup.C = setup.Ctx{
			Cwd:   cwd,
			Env:   map[string]string{},
			Chart: &helmv2.Chart{Metadata: &helmv2.Metadata{Name: name}},
		}
	}

	// a package level flag reads its environment only on the first parse, so the pipe
	// and the versions are seeded rather than parsed.
	run := func(
		runner *tests.TestingCommandRunner,
		pipe Pipe,
		versions []string,
		task func(*TaskList) *Task,
	) error {
		GinkgoHelper()

		seed("charts/app", "app")
		*P = pipe
		C.Versions = versions

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-helm",
			CommandName: "publish",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(task(tl).Job())
						})
				},
			},
		}).Run()
	}

	pipe := func() Pipe {
		return Pipe{Chart: Chart{
			Target:      "oci://registry.example.com/charts",
			Destination: "./dist/",
		}}
	}

	Describe("Helm package", func() {
		It("packages the chart in the directory the setup step resolved", func() {
			runner := fixtures.Runner()

			Expect(run(runner, pipe(), []string{"1.0.0"}, HelmPackage)).To(Succeed())

			invocation, ok := runner.LastInvocation()
			Expect(ok).To(BeTrue())
			Expect(invocation.Name).To(Equal("helm"))
			Expect(invocation.Args).To(Equal([]string{"package", "-d", "./dist/", ".", "--version", "1.0.0"}))
			Expect(invocation.Dir).To(Equal("charts/app"))
		})

		// the application version is what the chart reports as the version of the thing
		// it deploys, which most charts leave to whatever is committed in Chart.yaml.
		It("carries the application version only when one was given", func() {
			runner := fixtures.Runner()

			p := pipe()
			p.Chart.AppVersion = "2.3.4"

			Expect(run(runner, p, []string{"1.0.0"}, HelmPackage)).To(Succeed())

			invocation, _ := runner.LastInvocation()
			Expect(invocation.Args).To(ContainElements("--app-version", "2.3.4"))
		})

		It("packages every version it was given", func() {
			runner := fixtures.Runner()

			Expect(run(runner, pipe(), []string{"1.0.0", "1.0"}, HelmPackage)).To(Succeed())

			Expect(runner.InvocationNames()).To(Equal([]string{"helm", "helm"}))
		})

		// nothing to package is a pipeline whose version conditions selected nothing,
		// not a failure, so the task disables itself rather than running helm on it.
		It("runs nothing without a version", func() {
			runner := fixtures.Runner()

			Expect(run(runner, pipe(), []string{}, HelmPackage)).To(Succeed())

			Expect(runner.InvocationNames()).To(BeEmpty())
		})
	})

	Describe("Helm publish", func() {
		It("pushes the archive the package task wrote, named after the chart", func() {
			runner := fixtures.Runner()

			Expect(run(runner, pipe(), []string{"1.0.0"}, HelmPublish)).To(Succeed())

			invocation, ok := runner.LastInvocation()
			Expect(ok).To(BeTrue())
			Expect(invocation.Name).To(Equal("helm"))
			Expect(invocation.Args).
				To(Equal([]string{"push", "dist/app-1.0.0.tgz", "oci://registry.example.com/charts"}))
			Expect(invocation.Dir).To(Equal("charts/app"))
		})

		It("pushes one archive per version", func() {
			runner := fixtures.Runner()

			Expect(run(runner, pipe(), []string{"1.0.0", "1.0"}, HelmPublish)).To(Succeed())

			Expect(runner.InvocationNames()).To(Equal([]string{"helm", "helm"}))
		})

		It("runs nothing without a version", func() {
			runner := fixtures.Runner()

			Expect(run(runner, pipe(), []string{}, HelmPublish)).To(Succeed())

			Expect(runner.InvocationNames()).To(BeEmpty())
		})
	})
})
