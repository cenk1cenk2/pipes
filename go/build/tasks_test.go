package build

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var _ = Describe("Go build", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		// The build inherits the environment it runs in and the Taskfile exports
		// CGO_ENABLED for every task, so the spec below would assert against the
		// setting of whatever ran it rather than the one the pipe put there.
		tests.WithoutEnvironment("CGO_ENABLED")

		*P = pipe
		deps := Deps{Tool: &tool.Ctx{Cwd: "projects/api", Env: map[string]string{}}}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(GoBuild(tl, deps).Job())
						})
				},
			},
		}).Run()
	}

	pipe := func() Pipe {
		return Pipe{
			Output:         "./dist/",
			BinaryName:     "bin",
			BinaryTemplate: "{{ .name }}{{ if .os }}-{{ .os }}{{ end }}{{ if .arch }}-{{ .arch }}{{ end }}",
			BuildTargets:   []GoBuildTarget{{Os: "linux", Arch: "amd64"}},
		}
	}

	It("builds the vendored module into an artifact named after its target", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe())).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("go"))
		Expect(invocation.Dir).To(Equal("projects/api"))
		Expect(invocation.Args).
			To(Equal([]string{"build", "-mod=vendor", "-v", "-ldflags=", "-o", "dist/bin-linux-amd64"}))
	})

	// Cross compiling a target the runner has no toolchain for is what the pipe is
	// for, so the build stays static unless the pipeline asks for CGO.
	It("disables CGO unless the pipeline enabled it", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe())).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Env).To(ContainElement("CGO_ENABLED=0"))

		p := pipe()
		p.EnableCGO = true

		enabled := fixtures.Runner()
		Expect(run(enabled, p)).To(Succeed())

		invocation, ok = enabled.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Env).NotTo(ContainElement("CGO_ENABLED=0"))
	})

	It("builds every target that was asked for", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.BuildTargets = []GoBuildTarget{{Os: "linux", Arch: "amd64"}, {Os: "darwin", Arch: "arm64"}}

		Expect(run(runner, p)).To(Succeed())

		outputs := []string{}
		for _, invocation := range runner.Invocations() {
			outputs = append(outputs, invocation.Args[len(invocation.Args)-1])
		}

		Expect(outputs).To(ConsistOf("dist/bin-linux-amd64", "dist/bin-darwin-arm64"))
	})

	It("turns the build variables into linker flags", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.BuildVariables = map[string]string{"main.VERSION": "v1.0.0"}

		Expect(run(runner, p)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(ContainElement("-ldflags=-X main.VERSION=v1.0.0"))
	})

	It("passes the build tags as a single comma separated argument", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.BuildTags = []string{"netgo", "osusergo"}

		Expect(run(runner, p)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(ContainElements("-tags", "netgo,osusergo"))
	})
})
