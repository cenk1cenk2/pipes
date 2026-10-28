package install

import (
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Node install", func() {
	BeforeEach(func() {
		GinkgoT().Setenv("COREPACK_HOME", "")
	})

	seed := func(packageManager string) {
		*setup.NodeCtx = node.Ctx{PackageManager: node.PackageManager{
			Exe:      packageManager,
			Commands: node.PackageManagers[packageManager],
		}}
	}

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, packageManager string) error {
		GinkgoHelper()

		*P = pipe
		seed(packageManager)

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(install(tl).Job(), prune(tl).Job())
						})
				},
			},
		}).Run()
	}

	pipe := func() Pipe {
		return Pipe{Install: Install{Cwd: "projects/web", Args: "--foo"}}
	}

	It("installs against the lockfile of the resolved package manager", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.UseLockFile = true

		Expect(run(runner, p, "pnpm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pnpm"))
		Expect(invocation.Args).To(Equal([]string{"i", "--frozen-lockfile", "--foo"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	It("installs without the lockfile when the pipeline asked for it", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), "npm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"i", "--unsafe-perm", "--foo"}))
	})

	// the cache directory is named after the package manager, and the flag it is
	// passed under differs between them.
	It("points the package manager at a cache named after it", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.Cache = true

		Expect(run(runner, p, "yarn")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"install", "--foo", "--prefer-offline", "--cache-folder", ".yarn"}))
	})

	It("prunes the pnpm store it installed into", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.Cache = true

		Expect(run(runner, p, "pnpm")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"pnpm", "pnpm"}))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"store", "prune", "--store-dir", ".pnpm"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	It("leaves the pnpm store alone when nothing is cached", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), "pnpm")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"pnpm"}))
	})

	Describe("corepack cache", func() {
		var home string
		var cacheDir string
		var runner *tests.TestingCommandRunner
		var p Pipe

		BeforeEach(func() {
			home = GinkgoT().TempDir()
			cacheDir = filepath.Join(home, "v1", "pnpm")
			GinkgoT().Setenv("COREPACK_HOME", home)

			for _, version := range []string{"9.15.0", "10.12.1", "10.13.0"} {
				Expect(os.MkdirAll(filepath.Join(cacheDir, version), 0700)).To(Succeed())
				Expect(os.WriteFile(filepath.Join(cacheDir, version, "pnpm-native"), []byte("binary"), 0600)).To(Succeed())
			}

			runner = fixtures.Runner(tests.TestingCommandResponse{
				Name:   "pnpm",
				Args:   []string{"--version"},
				Stdout: " 10.12.1 \n",
				Stderr: "corepack warning\n",
			}, tests.TestingCommandResponse{}, tests.TestingCommandResponse{})
			p = pipe()
			p.Install.Cache = true
		})

		It("keeps the install project's version and removes obsolete version trees", func() {
			Expect(run(runner, p, "pnpm")).To(Succeed())

			Expect(filepath.Join(cacheDir, "10.12.1", "pnpm-native")).To(BeAnExistingFile())
			Expect(filepath.Join(cacheDir, "9.15.0")).NotTo(BeADirectory())
			Expect(filepath.Join(cacheDir, "10.13.0")).NotTo(BeADirectory())

			invocations := runner.Invocations()
			Expect(invocations).To(HaveLen(3))
			Expect(invocations[1].Args).To(Equal([]string{"store", "prune", "--store-dir", ".pnpm"}))
			Expect(invocations[2].Args).To(Equal([]string{"--version"}))
			Expect(invocations[2].Dir).To(Equal("projects/web"))
		})

		DescribeTable("leaves versions alone without an explicit corepack home",
			func(unset bool) {
				GinkgoT().Setenv("COREPACK_HOME", "")
				if unset {
					Expect(os.Unsetenv("COREPACK_HOME")).To(Succeed())
				}

				Expect(run(runner, p, "pnpm")).To(Succeed())

				Expect(os.ReadDir(cacheDir)).To(HaveLen(3))
				Expect(runner.Invocations()).To(HaveLen(2))
			},
			Entry("unset", true),
			Entry("empty", false),
		)

		It("leaves versions alone when the cache directory is missing", func() {
			Expect(os.Rename(cacheDir, filepath.Join(home, "versions"))).To(Succeed())

			Expect(run(runner, p, "pnpm")).To(Succeed())

			Expect(os.ReadDir(filepath.Join(home, "versions"))).To(HaveLen(3))
			Expect(cacheDir).NotTo(BeADirectory())
		})

		DescribeTable("leaves versions alone when the resolved version is unknown",
			func(output string) {
				runner = fixtures.Runner(tests.TestingCommandResponse{
					Name:   "pnpm",
					Args:   []string{"--version"},
					Stdout: output,
				}, tests.TestingCommandResponse{}, tests.TestingCommandResponse{})

				Expect(run(runner, p, "pnpm")).To(Succeed())

				Expect(os.ReadDir(cacheDir)).To(HaveLen(3))
			},
			Entry("empty", ""),
			Entry("whitespace", " \n"),
			Entry("unparsable", "not-a-version\n"),
			Entry("ambiguous", "10.12.1\n10.13.0\n"),
			Entry("unmatched", "10.14.0\n"),
		)

		It("leaves versions alone when the version probe fails", func() {
			failure := tests.TestingCommandFailure(1)
			runner = fixtures.Runner(tests.TestingCommandResponse{
				Name:   "pnpm",
				Args:   []string{"--version"},
				Stdout: "10.12.1\n",
				Result: &failure,
			}, tests.TestingCommandResponse{}, tests.TestingCommandResponse{})

			Expect(run(runner, p, "pnpm")).To(Succeed())

			Expect(os.ReadDir(cacheDir)).To(HaveLen(3))
		})

		It("leaves symlinks, files, other package managers and corepack metadata alone", func() {
			outside := GinkgoT().TempDir()
			Expect(os.Symlink(outside, filepath.Join(cacheDir, "linked"))).To(Succeed())
			Expect(os.WriteFile(filepath.Join(cacheDir, "file"), []byte("keep"), 0600)).To(Succeed())
			Expect(os.MkdirAll(filepath.Join(home, "v1", "yarn", "1.22.0"), 0700)).To(Succeed())
			Expect(os.Mkdir(filepath.Join(home, "temp"), 0700)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(home, "lastKnownGood.json"), []byte("{}"), 0600)).To(Succeed())

			Expect(run(runner, p, "pnpm")).To(Succeed())

			Expect(os.Readlink(filepath.Join(cacheDir, "linked"))).To(Equal(outside))
			Expect(outside).To(BeADirectory())
			Expect(filepath.Join(cacheDir, "file")).To(BeAnExistingFile())
			Expect(filepath.Join(home, "v1", "yarn", "1.22.0")).To(BeADirectory())
			Expect(filepath.Join(home, "temp")).To(BeADirectory())
			Expect(filepath.Join(home, "lastKnownGood.json")).To(BeAnExistingFile())
		})

		DescribeTable("leaves versions alone when the cache layout contains a symlink",
			func(component string) {
				target := filepath.Join(home, component)
				outside := filepath.Join(GinkgoT().TempDir(), "cache")
				Expect(os.Rename(target, outside)).To(Succeed())
				Expect(os.Symlink(outside, target)).To(Succeed())
				GinkgoT().Setenv("COREPACK_HOME", home+string(os.PathSeparator))

				Expect(run(runner, p, "pnpm")).To(Succeed())

				Expect(os.ReadDir(cacheDir)).To(HaveLen(3))
			},
			Entry("home", "."),
			Entry("layout", "v1"),
			Entry("pnpm", "v1/pnpm"),
		)

		It("does not treat a symlink as the matching version", func() {
			Expect(os.RemoveAll(filepath.Join(cacheDir, "10.12.1"))).To(Succeed())
			Expect(os.Symlink(GinkgoT().TempDir(), filepath.Join(cacheDir, "10.12.1"))).To(Succeed())

			Expect(run(runner, p, "pnpm")).To(Succeed())

			Expect(os.ReadDir(cacheDir)).To(HaveLen(3))
		})

		DescribeTable("leaves versions alone when pruning is disabled",
			func(cache bool, packageManager string) {
				p.Install.Cache = cache

				Expect(run(runner, p, packageManager)).To(Succeed())

				Expect(os.ReadDir(cacheDir)).To(HaveLen(3))
				Expect(runner.Invocations()).To(HaveLen(1))
			},
			Entry("cache disabled", false, "pnpm"),
			Entry("npm", true, "npm"),
			Entry("yarn", true, "yarn"),
		)
	})

	// only pnpm prunes what nothing references, the other caches are left as they are.
	DescribeTable("leaves the cache of the other package managers alone",
		func(packageManager string) {
			runner := fixtures.Runner()

			p := pipe()
			p.Install.Cache = true

			Expect(run(runner, p, packageManager)).To(Succeed())

			Expect(runner.InvocationNames()).To(Equal([]string{packageManager}))
		},
		Entry("npm", "npm"),
		Entry("yarn", "yarn"),
	)
})
