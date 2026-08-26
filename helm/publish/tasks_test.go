package publish

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"
	helmv2 "helm.sh/helm/v4/pkg/chart/v2"

	"gitlab.kilic.dev/devops/pipes/helm/setup"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

func deps(cwd, name string) Deps {
	return Deps{
		Tool: &setup.Ctx{
			Ctx:   &tool.Ctx{Cwd: cwd, Env: map[string]string{}},
			Chart: &helmv2.Chart{Metadata: &helmv2.Metadata{Name: name}},
		},
	}
}

// The pipe is seeded rather than parsed out of the flags: a package level flag
// only reads its environment sources on the first parse of a process, so a suite
// that drove them would depend on the order the specs happen to run in. The
// versions are seeded for the same reason the pipe is, since collecting them
// reads the git references of whatever checkout the specs run in.
func run(
	runner *tests.TestingCommandRunner,
	pipe Pipe,
	versions []string,
	task func(*plumber.TaskList, Deps) *plumber.Task,
) error {
	GinkgoHelper()

	deps := deps("charts/app", "app")
	*P = pipe
	C.Versions = versions

	return fixtures.Cli(runner, tests.TaskListCli{
		AppName:     "pipe-helm",
		CommandName: "publish",
		TaskLists: []tests.TaskListFactory{
			func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
				tl := &plumber.TaskList{}

				return tl.New(p).
					SetRuntimeDepth(3).
					Set(func(tl *plumber.TaskList) plumber.Job {
						return plumber.JobSequence(task(tl, deps).Job())
					})
			},
		},
	}).Run()
}

func pipe() Pipe {
	return Pipe{HelmChart: HelmChart{
		Target:      "oci://registry.example.com/charts",
		Destination: "./dist/",
	}}
}

var _ = Describe("Helm package", func() {
	It("packages the chart in the directory the setup step resolved", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), []string{"1.0.0"}, HelmPackage)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("helm"))
		Expect(invocation.Args).To(Equal([]string{"package", "-d", "./dist/", ".", "--version", "1.0.0"}))
		Expect(invocation.Dir).To(Equal("charts/app"))
	})

	// The application version is what the chart reports as the version of the thing
	// it deploys, which most charts leave to whatever is committed in Chart.yaml.
	It("carries the application version only when one was given", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.HelmChart.AppVersion = "2.3.4"

		Expect(run(runner, p, []string{"1.0.0"}, HelmPackage)).To(Succeed())

		invocation, _ := runner.LastInvocation()
		Expect(invocation.Args).To(ContainElements("--app-version", "2.3.4"))
	})

	It("packages every version it was given", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), []string{"1.0.0", "1.0"}, HelmPackage)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"helm", "helm"}))
	})

	// Nothing to package is a pipeline whose version conditions selected nothing,
	// not a failure, so the task disables itself rather than running helm on it.
	It("runs nothing without a version", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), []string{}, HelmPackage)).To(Succeed())

		Expect(runner.InvocationNames()).To(BeEmpty())
	})
})

var _ = Describe("Helm publish", func() {
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
