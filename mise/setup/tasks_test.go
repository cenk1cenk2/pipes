package setup

import (
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Mise version", func() {
	probe := func(stdout string) error {
		GinkgoHelper()

		*C = Ctx{Executable: "/usr/local/bin/mise", Env: map[string]string{}}

		runner := fixtures.Runner(tests.TestingCommandResponse{
			Name:   "/usr/local/bin/mise",
			Args:   []string{"--version"},
			Stdout: stdout,
		})

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-mise",
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
		}).Run()
	}

	// the version is what the copy in the data directory is compared against, so
	// the platform and the build date that trail it are left out.
	It("takes the version out of the banner", func() {
		Expect(probe("2026.9.14 linux-x64 (2026-09-15)\n")).To(Succeed())

		Expect(C.Version).To(Equal("2026.9.14"))
	})

	It("fails when the binary answers with nothing", func() {
		Expect(probe("\n")).NotTo(Succeed())
	})
})

var _ = Describe("Mise environment", func() {
	run := func(dataDir string) {
		GinkgoHelper()

		P.DataDir = dataDir
		*C = Ctx{Env: map[string]string{}}

		Expect(fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-mise",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(environment(tl).Job())
						})
				},
			},
		}).Run()).To(Succeed())
	}

	// the shims link into the data directory by absolute path, so a relative one
	// would leave them pointing at wherever the job happened to run.
	It("resolves the data directory and the binary inside it to absolute paths", func() {
		run("./.mise/")

		dataDir, err := filepath.Abs(".mise")
		Expect(err).NotTo(HaveOccurred())

		Expect(C.DataDir).To(Equal(dataDir))
		Expect(C.Binary).To(Equal(filepath.Join(dataDir, "bin", "mise")))
		Expect(C.Env).To(HaveKeyWithValue("MISE_DATA_DIR", dataDir))
	})

	// mise links the shims to the first mise on the path, which would otherwise be
	// the one the image ships and no job restoring the directory has.
	It("puts the binary of the data directory first on the path", func() {
		GinkgoT().Setenv("PATH", "/usr/local/bin:/usr/bin")

		run("/builds/project/.mise")

		Expect(C.Env).To(HaveKeyWithValue("PATH", "/builds/project/.mise/bin:/usr/local/bin:/usr/bin"))
	})

	It("carries the rest of the process environment along", func() {
		GinkgoT().Setenv("MISE_GITHUB_TOKEN", "token")

		run("/builds/project/.mise")

		Expect(C.Env).To(HaveKeyWithValue("MISE_GITHUB_TOKEN", "token"))
	})
})
