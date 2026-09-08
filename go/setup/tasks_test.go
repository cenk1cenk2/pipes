package setup

import (
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Go version", func() {
	// the toolchain prints a banner no pattern narrows down, so what it answers is
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
	// the workspace is deliberately not reset, so a run overwrites what the last one decided.
	run := func(cwd string) error {
		GinkgoHelper()

		C.Cwd = cwd
		C.Env = map[string]string{}

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(workspace(tl).Job())
						})
				},
			},
		}).Run()
	}

	workspaced := func() string {
		GinkgoHelper()

		dir := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(dir, "go.work"), []byte("go 1.27\n"), 0o600)).To(Succeed())

		return dir
	}

	It("drives the modules as a workspace when the working directory holds the workspace file", func() {
		Expect(run(workspaced())).To(Succeed())

		Expect(C.Workspace).To(BeTrue())
	})

	// a module of a bigger workspace is built on its own, so the workspace file of
	// a parent directory is none of this run's business.
	It("stays on the single module when only a parent directory holds the workspace file", func() {
		child := filepath.Join(workspaced(), "api")
		Expect(os.Mkdir(child, 0o700)).To(Succeed())

		Expect(run(child)).To(Succeed())

		Expect(C.Workspace).To(BeFalse())
	})

	It("forgets a workspace that the current run did not detect", func() {
		Expect(run(workspaced())).To(Succeed())
		Expect(C.Workspace).To(BeTrue())

		Expect(run(GinkgoT().TempDir())).To(Succeed())
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
							return JobSequence(modules(tl).Job())
						})
				},
			},
		}).Run()
	}

	golist := func(stdout string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{Name: "go", Stdout: stdout}
	}

	It("resolves every module the workspace drives", func() {
		runner := fixtures.Runner(golist("/repository/api\n/repository/worker\n"))

		Expect(run(runner, true)).To(Succeed())

		Expect(C.Modules).To(Equal([]string{"/repository/api", "/repository/worker"}))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("go"))
		Expect(invocation.Args).To(Equal([]string{"list", "-m", "-f", "{{.Dir}}"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
		Expect(invocation.Env).To(ContainElement("GOPATH=/cache"))
	})

	It("asks the toolchain for nothing outside workspace mode", func() {
		runner := fixtures.Runner()

		Expect(run(runner, false)).To(Succeed())

		Expect(C.Modules).To(BeEmpty())
		Expect(runner.Invocations()).To(BeEmpty())
	})
})
