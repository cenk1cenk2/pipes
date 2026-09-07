package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Go version", func() {
	// The toolchain prints a banner no pattern narrows down, so what it answers is
	// reported whole once the surrounding whitespace is off it.
	It("reports the whole banner the toolchain printed", func() {
		*C = Ctx{Env: map[string]string{}}

		runner := fixtures.Runner(tests.TestingCommandResponse{
			Name:   "go",
			Args:   []string{"version"},
			Stdout: "go version go1.27.0 linux/amd64\n",
		})

		Expect(fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(version(tl).Job())
						})
				},
			},
		}).Run()).To(Succeed())

		Expect(C.Version).To(Equal("go version go1.27.0 linux/amd64"))
	})
})

var _ = Describe("Go workspace", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		// Only the resolved tool is rebuilt, so what the previous run decided about
		// the workspace is still there for the run under test to overwrite.
		*P = pipe
		C.Cwd = "projects/api"
		C.Env = map[string]string{}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(GoWorkspace(tl).Job())
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

var _ = Describe("Go modules", func() {
	run := func(runner *tests.TestingCommandRunner, workspace bool) error {
		GinkgoHelper()

		*C = Ctx{Cwd: "projects/api", Env: map[string]string{"GOPATH": "/cache"}, Workspace: workspace}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(GoModules(tl).Job())
						})
				},
			},
		}).Run()
	}

	modules := func(stdout string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{Name: "go", Stdout: stdout}
	}

	It("resolves every module the workspace drives", func() {
		runner := fixtures.Runner(modules("/repository/api\n/repository/worker\n"))

		Expect(run(runner, true)).To(Succeed())

		Expect(C.Modules).To(Equal([]string{"/repository/api", "/repository/worker"}))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("go"))
		Expect(invocation.Args).To(Equal([]string{"list", "-m", "-f", "{{.Dir}}"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
		Expect(invocation.Env).To(ContainElement("GOPATH=/cache"))
	})

	// The scaffold module lives under a directory the go tool drops out of package
	// patterns, so it is kept here and linted from inside its own directory.
	It("keeps the modules the go tool hides from package patterns", func() {
		runner := fixtures.Runner(modules("/repository/_template\n/repository/api\n"))

		Expect(run(runner, true)).To(Succeed())

		Expect(C.Modules).To(Equal([]string{"/repository/_template", "/repository/api"}))
	})

	It("asks the toolchain for nothing outside workspace mode", func() {
		runner := fixtures.Runner()

		Expect(run(runner, false)).To(Succeed())

		Expect(C.Modules).To(BeEmpty())
		Expect(runner.Invocations()).To(BeEmpty())
	})
})
