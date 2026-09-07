package build

import (
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
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
		// The task reads the tool the setup resolved off its package level
		// instances, so a spec seeds those the same way it seeds its own.
		*setup.P = setup.Pipe{}
		*setup.C = setup.Ctx{Cwd: "projects/api", Env: map[string]string{}}
		*C = Ctx{}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(GoBuild(tl).Job())
						})
				},
			},
		}).Run()
	}

	commands := func(dirs ...string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{Name: "go", Stdout: strings.Join(dirs, "\n") + "\n"}
	}

	runWorkspace := func(runner *tests.TestingCommandRunner, pipe Pipe, workspace bool, modules ...string) error {
		GinkgoHelper()

		tests.WithoutEnvironment("CGO_ENABLED")

		*P = pipe
		*setup.P = setup.Pipe{Workspace: workspace}
		*setup.C = setup.Ctx{
			Cwd:       "projects/api",
			Env:       map[string]string{},
			Workspace: true,
			Modules:   modules,
		}
		*C = Ctx{}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(GoBuildPackages(tl).Job(), GoBuild(tl).Job())
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

	// The workspace build follows the flag rather than what the setup resolved:
	// the probe behind setup.C.Workspace is on for any invocation that merely sits
	// inside a workspace, and a child pipeline building one module must not start
	// building all of them.
	builds := func(runner *tests.TestingCommandRunner) []CommandInvocation {
		found := []CommandInvocation{}
		for _, invocation := range runner.Invocations() {
			if len(invocation.Args) > 0 && invocation.Args[0] == "build" {
				found = append(found, invocation)
			}
		}

		return found
	}

	It("builds every command of the workspace when the flag asked for it", func() {
		runner := fixtures.Runner(
			commands("/repository/_template"),
			commands("/repository/api"),
		)

		Expect(runWorkspace(runner, pipe(), true, "/repository/_template", "/repository/api")).To(Succeed())

		dirs := []string{}
		for _, invocation := range builds(runner) {
			dirs = append(dirs, invocation.Dir)
			Expect(invocation.Args).
				To(Equal([]string{"build", "-mod=vendor", "-v", "-ldflags=", "-o", "dist/bin-linux-amd64"}))
		}

		Expect(dirs).To(ConsistOf("/repository/_template", "/repository/api"))
	})

	// A workspace carries library modules next to the commands, and asking the go
	// tool what each module holds is what keeps those out of the build.
	It("leaves out a module that holds no command", func() {
		runner := fixtures.Runner(
			commands(""),
			commands("/repository/api"),
		)

		Expect(runWorkspace(runner, pipe(), true, "/repository/internal", "/repository/api")).To(Succeed())

		Expect(builds(runner)).To(HaveLen(1))
		Expect(builds(runner)[0].Dir).To(Equal("/repository/api"))
	})

	It("builds only the working directory when the workspace was merely detected", func() {
		runner := fixtures.Runner()

		Expect(runWorkspace(runner, pipe(), false, "/repository/_template", "/repository/api")).To(Succeed())

		Expect(runner.Invocations()).To(HaveLen(1))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Dir).To(Equal("projects/api"))
	})

	It("builds every command for every target it was asked for", func() {
		runner := fixtures.Runner(
			commands("/repository/api"),
			commands("/repository/worker"),
		)

		p := pipe()
		p.BuildTargets = []GoBuildTarget{{Os: "linux", Arch: "amd64"}, {Os: "darwin", Arch: "arm64"}}

		Expect(runWorkspace(runner, p, true, "/repository/api", "/repository/worker")).To(Succeed())

		built := []string{}
		for _, invocation := range builds(runner) {
			built = append(built, fmt.Sprintf("%s %s", invocation.Dir, invocation.Args[len(invocation.Args)-1]))
		}

		Expect(built).To(ConsistOf(
			"/repository/api dist/bin-linux-amd64",
			"/repository/api dist/bin-darwin-arm64",
			"/repository/worker dist/bin-linux-amd64",
			"/repository/worker dist/bin-darwin-arm64",
		))
	})
})
