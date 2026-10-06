package install

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/mise/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Mise binary", func() {
	var executable string

	BeforeEach(func() {
		dir := GinkgoT().TempDir()

		executable = filepath.Join(dir, "image", "mise")
		Expect(os.MkdirAll(filepath.Dir(executable), 0o700)).To(Succeed())
		Expect(os.WriteFile(executable, []byte("image"), 0o600)).To(Succeed())

		*setup.C = setup.Ctx{
			Version:    "2026.9.14",
			Executable: executable,
			DataDir:    filepath.Join(dir, ".mise"),
			Binary:     filepath.Join(dir, ".mise", "bin", "mise"),
			Env:        map[string]string{},
		}
	})

	run := func(runner *tests.TestingCommandRunner) error {
		GinkgoHelper()

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-mise",
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(binary(tl).Job())
						})
				},
			},
		}).Run()
	}

	cached := func(contents string) {
		GinkgoHelper()

		Expect(os.MkdirAll(filepath.Dir(setup.C.Binary), 0o700)).To(Succeed())
		Expect(os.WriteFile(setup.C.Binary, []byte(contents), 0o600)).To(Succeed())
	}

	answers := func(stdout string) tests.TestingCommandResponse {
		return tests.TestingCommandResponse{Name: setup.C.Binary, Args: []string{"--version"}, Stdout: stdout}
	}

	expectCopied := func() {
		GinkgoHelper()

		contents, err := os.ReadFile(setup.C.Binary)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal("image"))

		info, err := os.Stat(setup.C.Binary)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Mode().Perm()).To(Equal(os.FileMode(0o755)))
	}

	It("copies the binary of the image into an empty data directory", func() {
		runner := fixtures.Runner()

		Expect(run(runner)).To(Succeed())

		expectCopied()
		Expect(runner.Invocations()).To(BeEmpty())
	})

	It("keeps a copy that is already at the version of the image", func() {
		cached("cached")

		Expect(run(fixtures.Runner(answers("2026.9.14 linux-x64 (2026-09-15)\n")))).To(Succeed())

		contents, err := os.ReadFile(setup.C.Binary)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal("cached"))
	})

	It("replaces a copy of another version", func() {
		cached("cached")

		Expect(run(fixtures.Runner(answers("2026.8.1 linux-x64 (2026-08-02)\n")))).To(Succeed())

		expectCopied()
	})

	// a data directory restored from a cache built on another platform holds a
	// binary that does not run at all.
	It("replaces a copy that does not run", func() {
		cached("cached")

		runner := fixtures.Runner(tests.TestingCommandResponse{
			Name: setup.C.Binary,
			Err:  errors.New("exec format error"),
		})

		Expect(run(runner)).To(Succeed())

		expectCopied()
	})

	It("leaves nothing but the binary in its directory", func() {
		cached("cached")

		Expect(run(fixtures.Runner(answers("2026.8.1 linux-x64 (2026-08-02)\n")))).To(Succeed())

		entries, err := os.ReadDir(filepath.Dir(setup.C.Binary))
		Expect(err).NotTo(HaveOccurred())

		names := []string{}
		for _, entry := range entries {
			names = append(names, entry.Name())
		}

		Expect(names).To(Equal([]string{"mise"}))
	})

	// a runner node may still execute the copy restored from the cache, and
	// linux refuses to open a binary for writing while it runs.
	It("replaces a copy that is still running", func() {
		if runtime.GOOS != "linux" {
			Skip("text file busy is specific to linux")
		}

		sleep, err := exec.LookPath("sleep")
		Expect(err).NotTo(HaveOccurred())

		contents, err := os.ReadFile(sleep)
		Expect(err).NotTo(HaveOccurred())

		Expect(os.MkdirAll(filepath.Dir(setup.C.Binary), 0o700)).To(Succeed())
		Expect(os.WriteFile(setup.C.Binary, contents, 0o700)).To(Succeed())

		// a multi-call binary picks the applet from the name it was started as.
		cmd := &exec.Cmd{Path: setup.C.Binary, Args: []string{"sleep", "30"}}
		Expect(cmd.Start()).To(Succeed())
		DeferCleanup(func() {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		})

		Expect(run(fixtures.Runner(answers("2026.8.1 linux-x64 (2026-08-02)\n")))).To(Succeed())

		expectCopied()
		Expect(cmd.Process.Signal(syscall.Signal(0))).To(Succeed())
	})
})

var _ = Describe("Mise install", func() {
	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		*P = pipe
		*setup.C = setup.Ctx{
			Cwd:    "projects/api",
			Binary: "/builds/project/.mise/bin/mise",
			Env:    map[string]string{"MISE_DATA_DIR": "/builds/project/.mise"},
		}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-mise",
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(
								install(tl).Job(),
								prune(tl).Job(),
								list(tl).Job(),
							)
						})
				},
			},
		}).Run()
	}

	// the shims link to the mise that runs the install, so every command goes
	// through the copy in the data directory.
	It("installs, prunes and lists through the binary of the data directory", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Prune: true})).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).To(HaveLen(3))

		args := [][]string{}
		for _, invocation := range invocations {
			Expect(invocation.Name).To(Equal("/builds/project/.mise/bin/mise"))
			Expect(invocation.Dir).To(Equal("projects/api"))
			Expect(invocation.Env).To(ContainElement("MISE_DATA_DIR=/builds/project/.mise"))

			args = append(args, invocation.Args)
		}

		Expect(args).To(Equal([][]string{{"install"}, {"prune", "--yes"}, {"ls", "--current"}}))
	})

	It("appends the arguments to the install command", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Args: "--jobs 2"})).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).NotTo(BeEmpty())
		Expect(invocations[0].Args).To(Equal([]string{"install", "--jobs", "2"}))
	})

	It("keeps the unused versions when the pipeline turned pruning off", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{})).To(Succeed())

		Expect(runner.Invocations()).To(HaveLen(2))
		for _, invocation := range runner.Invocations() {
			Expect(invocation.Args).NotTo(ContainElement("prune"))
		}
	})
})
