package test

import (
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Go test", func() {
	var dir string

	BeforeEach(func() {
		dir = GinkgoT().TempDir()
	})

	runWith := func(runner *tests.TestingCommandRunner, pipe Pipe, ctx setup.Ctx) error {
		GinkgoHelper()

		*P = pipe
		*setup.P = setup.Pipe{}
		ctx.Env = map[string]string{}
		*setup.C = ctx
		*C = Ctx{}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "test",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		return runWith(runner, pipe, setup.Ctx{Cwd: dir})
	}

	It("lets ginkgo write the JUnit report itself", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Runner: string(RunnerGinkgo), JUnit: "junit.xml"})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(ContainElement("--junit-report=junit.xml"))
	})

	It("runs the go runner through gotestsum so it writes the JUnit report", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Runner: string(RunnerGo), JUnit: "junit.xml"})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("gotestsum"))
		Expect(invocation.Args).To(Equal([]string{"--junitfile=junit.xml", "--format=standard-verbose", "--", "-v", "-failfast", "./..."}))
	})

	It("runs the go runner bare when no JUnit report is asked for", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Runner: string(RunnerGo)})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("go"))
	})

	// a failing run is the one the report matters for.
	It("converts the coverage profile even when the tests fail", func() {
		Expect(os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/mod\n"), 0o600)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(dir, "coverage.out"), []byte("mode: set\nexample.com/mod/a.go:1.1,2.2 1 0\n"), 0o600)).To(Succeed())

		failure := tests.TestingCommandFailure(1)
		runner := fixtures.Runner(tests.TestingCommandResponse{Name: "go", Result: &failure})

		Expect(run(runner, Pipe{Runner: string(RunnerGinkgo), Coverage: "coverage.out", Cobertura: "cobertura.xml"})).NotTo(Succeed())

		Expect(filepath.Join(dir, "cobertura.xml")).To(BeAnExistingFile())
	})

	It("converts the coverage profile into a Cobertura report", func() {
		Expect(os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/mod\n"), 0o600)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(dir, "coverage.out"), []byte("mode: set\nexample.com/mod/a.go:1.1,2.2 1 1\n"), 0o600)).To(Succeed())

		Expect(run(fixtures.Runner(), Pipe{Runner: string(RunnerGinkgo), Coverage: "coverage.out", Cobertura: "cobertura.xml"})).To(Succeed())

		Expect(filepath.Join(dir, "cobertura.xml")).To(BeAnExistingFile())
	})

	It("skips the Cobertura report when the run left no coverage profile", func() {
		Expect(run(fixtures.Runner(), Pipe{Runner: string(RunnerGinkgo), Coverage: "coverage.out", Cobertura: "cobertura.xml"})).To(Succeed())

		Expect(filepath.Join(dir, "cobertura.xml")).NotTo(BeAnExistingFile())
	})

	// two modules with the same report names, which each have to keep their own.
	workspace := func() []string {
		GinkgoHelper()

		modules := []string{filepath.Join(dir, "a"), filepath.Join(dir, "b")}
		for _, module := range modules {
			name := filepath.Base(module)
			Expect(os.MkdirAll(module, 0o755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.com/"+name+"\n"), 0o600)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(module, "coverage.out"), []byte("mode: set\nexample.com/"+name+"/"+name+".go:1.1,2.2 1 1\n"), 0o600)).To(Succeed())
		}

		return modules
	}

	It("tests every module of a workspace in its own directory and keeps each report there", func() {
		modules := workspace()
		runner := fixtures.Runner()

		Expect(runWith(runner, Pipe{Runner: string(RunnerGinkgo), Coverage: "coverage.out", JUnit: "junit.xml", Cobertura: "cobertura.xml"}, setup.Ctx{
			Cwd:       dir,
			Workspace: true,
			Modules:   modules,
		})).To(Succeed())

		dirs := []string{}
		for _, invocation := range runner.Invocations() {
			Expect(invocation.Args).To(ContainElements("--junit-report=junit.xml", "--coverprofile=coverage.out"))
			dirs = append(dirs, invocation.Dir)
		}
		Expect(dirs).To(Equal(modules))

		for _, module := range modules {
			body, err := os.ReadFile(filepath.Join(module, "cobertura.xml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(ContainSubstring("<source>" + module + "</source>"))
			Expect(string(body)).To(ContainSubstring(`filename="` + filepath.Base(module) + `.go"`))
		}
	})

	It("refuses a report path that leaves the modules of a workspace before running anything", func() {
		modules := workspace()

		for _, p := range []Pipe{
			{Runner: string(RunnerGinkgo), JUnit: filepath.Join(dir, "junit.xml")},
			{Runner: string(RunnerGinkgo), Coverage: "coverage.out", Cobertura: "../cobertura.xml"},
			{Runner: string(RunnerGinkgo), Coverage: filepath.Join(dir, "coverage.out")},
		} {
			runner := fixtures.Runner()

			Expect(runWith(runner, p, setup.Ctx{Cwd: dir, Workspace: true, Modules: modules})).
				To(MatchError(ContainSubstring("otherwise the modules overwrite each other's report")))
			Expect(runner.Invocations()).To(BeEmpty())
		}
	})

	It("takes a report path anywhere for a single module", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Runner: string(RunnerGinkgo), JUnit: filepath.Join(dir, "reports", "junit.xml")})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(ContainElement("--junit-report=" + filepath.Join(dir, "reports", "junit.xml")))
	})
})
