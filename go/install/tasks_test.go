package install

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

var _ = Describe("Go install", func() {
	seed := func(pipe Pipe, workspace bool, env map[string]string) {
		GinkgoHelper()

		*P = pipe
		*setup.C = setup.Ctx{
			Cwd:       "projects/api",
			Env:       env,
			Workspace: workspace,
		}
	}

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	cached := func(runner *tests.TestingCommandRunner, pipe Pipe, workspace bool, env map[string]string) error {
		GinkgoHelper()

		seed(pipe, workspace, env)

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-go",
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(
								vendor(tl).Job(),
								verify(tl).Job(),
								prune(tl).Job(),
							)
						})
				},
			},
		}).Run()
	}

	run := func(runner *tests.TestingCommandRunner, pipe Pipe, workspace bool) error {
		GinkgoHelper()

		return cached(runner, pipe, workspace, map[string]string{"GOPATH": "/cache"})
	}

	disabled := func(workspace bool, task func(*TaskList) *Task) bool {
		GinkgoHelper()

		seed(Pipe{}, workspace, map[string]string{})

		tl := &TaskList{}
		tl.New(NewPlumber(func(_ *Plumber) *cli.Command {
			return &cli.Command{Name: "test"}
		}))

		return task(tl).IsDisabled()
	}

	It("vendors the single module in the working directory", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{}, false)).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Name).To(Equal("go"))
		Expect(invocations[0].Args).To(Equal([]string{"mod", "vendor"}))
		Expect(invocations[0].Dir).To(Equal("projects/api"))
		Expect(invocations[0].Env).To(ContainElement("GOPATH=/cache"))
	})

	It("vendors the whole workspace when the setup resolved one", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{}, true)).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Args).To(Equal([]string{"work", "vendor"}))
	})

	It("appends the arguments to whichever vendor command it picked", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Args: "-e"}, true)).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Args).To(Equal([]string{"work", "vendor", "-e"}))

		single := fixtures.Runner()

		Expect(run(single, Pipe{Args: "-e"}, false)).To(Succeed())

		invocations = single.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Args).To(Equal([]string{"mod", "vendor", "-e"}))
	})

	// the two vendor tasks are the halves of the workspace condition, so exactly
	// one of them runs under the parent whichever way the setup resolved.
	DescribeTable(
		"vendors through exactly one of the two tasks",
		func(workspace bool, module, workspaceDisabled bool) {
			Expect(disabled(workspace, vendorModule)).To(Equal(module))
			Expect(disabled(workspace, vendorWorkspace)).To(Equal(workspaceDisabled))
		},
		Entry("no workspace", false, false, true),
		Entry("workspace", true, true, false),
	)

	It("verifies the modules unless the pipeline turned it off", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Verify: true}, false)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"mod", "verify"}))
		Expect(invocation.Dir).To(Equal("projects/api"))

		disabled := fixtures.Runner()
		Expect(run(disabled, Pipe{}, false)).To(Succeed())

		Expect(disabled.Invocations()).To(HaveLen(1))
	})

	// the cache holds what an older dependency graph fetched beside what this one
	// selects, a replacement, the checksum database and a toolchain go switched to.
	populate := func(download string) {
		GinkgoHelper()

		for _, file := range []string{
			"golang.org/x/text/@v/list",
			"golang.org/x/text/@v/v0.20.0.info",
			"golang.org/x/text/@v/v0.20.0.mod",
			"golang.org/x/text/@v/v0.20.0.zip",
			"golang.org/x/text/@v/v0.20.0.ziphash",
			"golang.org/x/text/@v/v0.20.0.lock",
			"golang.org/x/text/@v/v0.21.0.info",
			"golang.org/x/text/@v/v0.21.0.mod",
			"golang.org/x/text/@v/v0.21.0.zip",
			"golang.org/x/text/@v/v0.21.0.ziphash",
			"golang.org/x/text/@v/v0.21.0.lock",
			"github.com/!burnt!sushi/toml/@v/v1.4.0.zip",
			"github.com/!burnt!sushi/toml/@v/v1.5.0.zip",
			"example.com/fork/@v/v1.0.0.zip",
			"example.com/upstream/@v/v1.0.0.zip",
			"sumdb/sum.golang.org/lookup/golang.org/x/text@v0.20.0",
			"golang.org/toolchain/@v/v0.0.1-go1.26.0.linux-amd64.zip",
		} {
			file = filepath.Join(download, filepath.FromSlash(file))

			Expect(os.MkdirAll(filepath.Dir(file), 0o755)).To(Succeed())
			Expect(os.WriteFile(file, []byte("cached"), 0o644)).To(Succeed())
		}
	}

	listed := func(workspace bool) *tests.TestingCommandRunner {
		GinkgoHelper()

		return fixtures.Runner(
			tests.TestingCommandResponse{Name: "go", Args: []string{"mod", "vendor"}},
			tests.TestingCommandResponse{Name: "go", Args: []string{"work", "vendor"}},
			tests.TestingCommandResponse{
				Name: "go",
				Args: []string{"list", "-mod=readonly", "-m", "-json", "all"},
				Stdout: `{
	"Path": "projects/api",
	"Main": true,
	"Dir": "/builds/projects/api"
}
{
	"Path": "golang.org/x/text",
	"Version": "v0.21.0"
}
{
	"Path": "github.com/BurntSushi/toml",
	"Version": "v1.5.0"
}
{
	"Path": "example.com/upstream",
	"Version": "v1.0.0",
	"Replace": {
		"Path": "example.com/fork",
		"Version": "v1.0.0"
	}
}
{
	"Path": "example.com/local",
	"Version": "v1.0.0",
	"Replace": {
		"Path": "../local",
		"Dir": "/builds/projects/local"
	}
}
`,
			},
		)
	}

	DescribeTable(
		"prunes the cache it installed from",
		func(workspace bool) {
			cache := GinkgoT().TempDir()
			download := filepath.Join(cache, "cache", "download")
			populate(download)

			runner := listed(workspace)

			Expect(cached(runner, Pipe{}, workspace, map[string]string{"GOMODCACHE": cache})).To(Succeed())

			invocation, ok := runner.LastInvocation()
			Expect(ok).To(BeTrue())
			Expect(invocation.Args).To(Equal([]string{"list", "-mod=readonly", "-m", "-json", "all"}))
			Expect(invocation.Dir).To(Equal("projects/api"))
			Expect(invocation.Env).To(ContainElement("GOMODCACHE=" + cache))

			for _, file := range []string{
				"golang.org/x/text/@v/v0.20.0.zip",
				"golang.org/x/text/@v/v0.20.0.ziphash",
				"golang.org/x/text/@v/v0.20.0.lock",
				"github.com/!burnt!sushi/toml/@v/v1.4.0.zip",
			} {
				Expect(filepath.Join(download, filepath.FromSlash(file))).NotTo(BeAnExistingFile())
			}

			for _, file := range []string{
				"golang.org/x/text/@v/list",
				"golang.org/x/text/@v/v0.20.0.info",
				"golang.org/x/text/@v/v0.20.0.mod",
				"golang.org/x/text/@v/v0.21.0.zip",
				"golang.org/x/text/@v/v0.21.0.ziphash",
				"golang.org/x/text/@v/v0.21.0.lock",
				"github.com/!burnt!sushi/toml/@v/v1.5.0.zip",
				"example.com/fork/@v/v1.0.0.zip",
				"example.com/upstream/@v/v1.0.0.zip",
				"sumdb/sum.golang.org/lookup/golang.org/x/text@v0.20.0",
				"golang.org/toolchain/@v/v0.0.1-go1.26.0.linux-amd64.zip",
			} {
				Expect(filepath.Join(download, filepath.FromSlash(file))).To(BeAnExistingFile())
			}
		},
		Entry("module", false),
		Entry("workspace", true),
	)

	It("finds nothing to prune in a cache nothing was downloaded into", func() {
		runner := listed(false)

		Expect(cached(runner, Pipe{}, false, map[string]string{"GOMODCACHE": GinkgoT().TempDir()})).To(Succeed())

		Expect(runner.Invocations()).To(HaveLen(2))
	})

	It("leaves the cache of the environment alone", func() {
		runner := fixtures.Runner()

		Expect(cached(runner, Pipe{}, false, map[string]string{})).To(Succeed())

		Expect(runner.Invocations()).To(HaveLen(1))
	})
})
