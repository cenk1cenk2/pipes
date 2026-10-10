package generate

import (
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Pulumi generate tasks", func() {
	var cwd string

	BeforeEach(func() {
		cwd = GinkgoT().TempDir()
		*setup.C = setup.Ctx{Cwd: cwd, Env: map[string]string{}}
		*stack.P = stack.Pipe{Stack: "main"}
		*C = Ctx{}
	})

	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		*P = pipe

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-pulumi",
			CommandName: "generate",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	whoami := func(stdout string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{
			Name:   "pulumi",
			Args:   []string{"whoami", "--output", "json"},
			Stdout: stdout + "\n",
		}
	}

	backend := func(url string) tests.TestingCommandResponse {
		return whoami(`{ "user": "runner", "url": "` + url + `" }`)
	}

	// stands in for the Pulumi program, writing the given files into the working
	// directory once the stack is deployed.
	writes := func(name string, files ...string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{
			Name: name,
			Match: func(_ CommandInvocation) bool {
				for _, file := range files {
					Expect(os.MkdirAll(filepath.Dir(filepath.Join(cwd, file)), 0o755)).To(Succeed())
					Expect(os.WriteFile(filepath.Join(cwd, file), []byte("generated"), 0o600)).To(Succeed())
				}

				return true
			},
		}
	}

	seed := func(files ...string) {
		GinkgoHelper()

		for _, file := range files {
			Expect(os.MkdirAll(filepath.Dir(filepath.Join(cwd, file)), 0o755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(cwd, file), []byte("stale"), 0o600)).To(Succeed())
		}
	}

	It("deploys a fresh stack on the local backend in the working directory", func() {
		runner := fixtures.Runner(backend("file://."), tests.TestingCommandResponse{Name: "pulumi"}, writes("pulumi"))

		Expect(run(runner, Pipe{})).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).To(HaveLen(3))
		Expect(invocations[1].Args).To(Equal([]string{"stack", "init", "-s", "main"}))
		Expect(invocations[2].Args).To(Equal([]string{"up", "--diff", "--yes", "-f", "--stack", "main"}))

		for _, invocation := range invocations {
			Expect(invocation.Dir).To(Equal(cwd))
			Expect(invocation.Env).To(ContainElement("PULUMI_CONFIG_PASSPHRASE="))
		}
	})

	It("builds the program before checking the backend, and runs the given command once the stack is initialized", func() {
		runner := fixtures.Runner(tests.TestingCommandResponse{Name: "sh"}, backend("file://."))

		Expect(run(runner, Pipe{Build: "npm run build", Command: "pulumi up --yes --refresh"})).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"sh", "pulumi", "pulumi", "sh"}))

		invocations := runner.Invocations()
		Expect(invocations[0].Args).To(Equal([]string{"-c", "npm run build"}))
		Expect(invocations[1].Args).To(Equal([]string{"whoami", "--output", "json"}))
		Expect(invocations[2].Args).To(Equal([]string{"stack", "init", "-s", "main"}))
		Expect(invocations[3].Args).To(Equal([]string{"-c", "pulumi up --yes --refresh"}))
		Expect(invocations[3].Dir).To(Equal(cwd))
	})

	It("replaces the state and the generated paths, so the files a regeneration drops are gone", func() {
		seed(
			".pulumi/stacks/probe/main.json",
			"manifests/removed/deployment.yaml",
			"manifests/service.yaml",
			"crds/stale.yaml",
			"untouched.yaml",
		)

		runner := fixtures.Runner(
			backend("file://."),
			tests.TestingCommandResponse{Name: "pulumi"},
			writes("pulumi", "manifests/service.yaml"),
		)

		Expect(run(runner, Pipe{Paths: []string{"manifests", "crds"}})).To(Succeed())

		Expect(filepath.Join(cwd, ".pulumi")).NotTo(BeAnExistingFile())
		Expect(filepath.Join(cwd, "manifests/removed")).NotTo(BeAnExistingFile())
		Expect(filepath.Join(cwd, "crds")).NotTo(BeAnExistingFile())
		Expect(os.ReadFile(filepath.Join(cwd, "manifests/service.yaml"))).To(Equal([]byte("generated")))
		Expect(filepath.Join(cwd, "untouched.yaml")).To(BeARegularFile())
	})

	DescribeTable("accepts the file:// backend of the working directory",
		func(url func(cwd string) string) {
			runner := fixtures.Runner(backend(url(cwd)))

			Expect(run(runner, Pipe{})).To(Succeed())
			Expect(runner.Invocations()).To(HaveLen(3))
		},
		Entry("relative to it", func(_ string) string { return "file://." }),
		Entry("as an absolute path", func(cwd string) string { return "file://" + cwd }),
	)

	DescribeTable("refuses a backend that is not the working directory, before touching anything",
		func(url string) {
			seed(".pulumi/meta.yaml", "manifests/service.yaml")

			runner := fixtures.Runner(tests.TestingCommandResponse{Name: "sh"}, backend(url))

			Expect(run(runner, Pipe{Build: "npm run build", Paths: []string{"manifests"}})).
				To(MatchError(ContainSubstring("has to use the file:// backend of the working directory: resolves to: " + url)))
			Expect(runner.InvocationNames()).To(Equal([]string{"sh", "pulumi"}))
			Expect(filepath.Join(cwd, ".pulumi/meta.yaml")).To(BeARegularFile())
			Expect(filepath.Join(cwd, "manifests/service.yaml")).To(BeARegularFile())
		},
		Entry("Pulumi Cloud", "https://api.pulumi.com"),
		Entry("an object storage", "s3://state-bucket"),
		Entry("the home directory", "file://~"),
		Entry("another absolute directory", "file:///elsewhere"),
		Entry("a parent directory", "file://.."),
	)

	DescribeTable("fails on a backend it can not read, before touching anything",
		func(response tests.TestingCommandResponse, message string) {
			seed(".pulumi/meta.yaml")

			runner := fixtures.Runner(response)

			Expect(run(runner, Pipe{})).To(MatchError(ContainSubstring(message)))
			Expect(runner.Invocations()).To(HaveLen(1))
			Expect(filepath.Join(cwd, ".pulumi/meta.yaml")).To(BeARegularFile())
		},
		Entry("an unparsable answer", whoami("not logged in"), "parse Pulumi backend"),
		Entry("a failing command", tests.TestingCommandResponse{
			Name:   "pulumi",
			Result: new(tests.TestingCommandFailure(1)),
		}, "command exited with code 1: $ pulumi whoami --output json"),
	)

	DescribeTable("validates the generated paths",
		func(paths []string, valid bool) {
			runner := fixtures.Runner(backend("file://."))

			err := run(runner, Pipe{Paths: paths})

			if valid {
				Expect(err).NotTo(HaveOccurred())

				return
			}

			Expect(err).To(MatchError(ContainSubstring("Generated path has to be inside the working directory")))
			Expect(runner.Invocations()).To(BeEmpty())
		},
		Entry("the working directory itself", []string{"."}, false),
		Entry("an empty path", []string{""}, false),
		Entry("a parent", []string{"../x"}, false),
		Entry("a directory inside", []string{"apps"}, true),
		Entry("a directory inside reached through another", []string{"apps/../namespaces"}, true),
	)
})
