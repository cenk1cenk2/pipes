package git_test

import (
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Repository", func() {
	var remote string

	run := func(dir string, args ...string) string {
		GinkgoHelper()

		command := exec.Command("git", args...)
		command.Dir = dir
		command.Env = append(
			command.Environ(),
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1",
		)

		output, err := command.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), string(output))

		return strings.TrimSpace(string(output))
	}

	// a bare remote with the target branch the sync writes into, reached over file://
	// so the shallow, partial clone behaves like it does against a server.
	BeforeEach(func() {
		// a git hook runs with the variables that locate the repository it fired in,
		// which every git call of a spec would otherwise act on instead of its own.
		for _, name := range strings.Fields(run("", "rev-parse", "--local-env-vars")) {
			if value, ok := os.LookupEnv(name); ok {
				Expect(os.Unsetenv(name)).To(Succeed())
				DeferCleanup(os.Setenv, name, value)
			}
		}

		base := GinkgoT().TempDir()
		bare := filepath.Join(base, "remote.git")
		seed := filepath.Join(base, "seed")

		run(base, "init", "--quiet", "--bare", "--initial-branch", "next", bare)
		run(bare, "config", "uploadpack.allowFilter", "true")
		run(base, "init", "--quiet", "--initial-branch", "next", seed)

		write(filepath.Join(seed, ".gitignore"), "generated/\n")
		write(filepath.Join(seed, "out", "kept.txt"), "kept")
		write(filepath.Join(seed, "out", "removed.txt"), "removed")
		write(filepath.Join(seed, "other", "untouched.txt"), "untouched")
		run(seed, "add", "--all")
		run(seed, "commit", "--quiet", "--message", "seed")
		run(seed, "push", "--quiet", bare, "next")

		remote = "file://" + bare
	})

	clone := func(ctx SpecContext, destinations ...string) git.Repository {
		GinkgoHelper()

		repository := git.Repository{Dir: filepath.Join(GinkgoT().TempDir(), "target")}
		Expect(repository.Clone(ctx, fixtures.Task(), remote, "next", destinations)).To(Succeed())

		return repository
	}

	staged := func(repository git.Repository) []string {
		GinkgoHelper()

		return strings.Split(run(repository.Dir, "diff", "--cached", "--name-status"), "\n")
	}

	It("checks out only the destinations of the target branch", func(ctx SpecContext) {
		repository := clone(ctx, "out")

		Expect(filepath.Join(repository.Dir, "out", "kept.txt")).To(BeAnExistingFile())
		Expect(filepath.Join(repository.Dir, "other")).NotTo(BeAnExistingFile())
	})

	It("checks out and stages a destination that is a single file", func(ctx SpecContext) {
		repository := clone(ctx, "out/kept.txt")
		Expect(filepath.Join(repository.Dir, "out", "removed.txt")).NotTo(BeAnExistingFile())

		source := filepath.Join(GinkgoT().TempDir(), "kept.txt")
		write(source, "changed")

		Expect(git.Sync(fixtures.Task(), repository.Dir, []git.Path{{Source: source, Destination: "out/kept.txt"}}, false)).To(Succeed())
		Expect(repository.Stage(ctx, fixtures.Task(), []string{"out/kept.txt"})).To(Succeed())

		Expect(staged(repository)).To(ConsistOf("M\tout/kept.txt"))
	})

	It("fails to clone a branch the target does not have", func(ctx SpecContext) {
		repository := git.Repository{Dir: filepath.Join(GinkgoT().TempDir(), "target")}

		Expect(repository.Clone(ctx, fixtures.Task(), remote, "missing", []string{"out"})).
			To(MatchError(ContainSubstring("missing")))
	})

	It("stages a file the source dropped as a deletion", func(ctx SpecContext) {
		repository := clone(ctx, "out")
		source := GinkgoT().TempDir()
		write(filepath.Join(source, "kept.txt"), "changed")

		Expect(git.Sync(fixtures.Task(), repository.Dir, []git.Path{{Source: source, Destination: "out"}}, false)).To(Succeed())
		Expect(repository.Stage(ctx, fixtures.Task(), []string{"out"})).To(Succeed())

		Expect(staged(repository)).To(ConsistOf("M\tout/kept.txt", "D\tout/removed.txt"))
	})

	// the target ignoring the destination is the normal case for generated code, and
	// a plain add stages nothing there without failing.
	It("stages a destination the target ignores", func(ctx SpecContext) {
		repository := clone(ctx, "generated")
		source := GinkgoT().TempDir()
		write(filepath.Join(source, "sdk.txt"), "sdk")

		Expect(git.Sync(fixtures.Task(), repository.Dir, []git.Path{{Source: source, Destination: "generated"}}, false)).To(Succeed())
		Expect(repository.Stage(ctx, fixtures.Task(), []string{"generated"})).To(Succeed())

		Expect(staged(repository)).To(ConsistOf("A\tgenerated/sdk.txt"))
	})

	It("writes the staged changes as a patch that applies onto the target", func(ctx SpecContext) {
		repository := clone(ctx, "out")
		source := GinkgoT().TempDir()
		write(filepath.Join(source, "kept.txt"), "changed")
		write(filepath.Join(source, "binary.bin"), "\x00\x01\x02")

		Expect(git.Sync(fixtures.Task(), repository.Dir, []git.Path{{Source: source, Destination: "out"}}, false)).To(Succeed())
		Expect(repository.Stage(ctx, fixtures.Task(), []string{"out"})).To(Succeed())

		patch := filepath.Join(GinkgoT().TempDir(), "changes.patch")
		Expect(repository.Patch(ctx, fixtures.Task(), patch)).To(Succeed())
		Expect(read(patch)).To(ContainSubstring("GIT binary patch"))

		other := clone(ctx, "out")
		run(other.Dir, "apply", "--index", patch)
		Expect(staged(other)).To(ConsistOf("M\tout/kept.txt", "D\tout/removed.txt", "A\tout/binary.bin"))
	})

	It("writes an empty patch when nothing changed", func(ctx SpecContext) {
		repository := clone(ctx, "out")
		Expect(repository.Stage(ctx, fixtures.Task(), []string{"out"})).To(Succeed())

		patch := filepath.Join(GinkgoT().TempDir(), "changes.patch")
		Expect(repository.Patch(ctx, fixtures.Task(), patch)).To(Succeed())
		Expect(read(patch)).To(BeEmpty())
	})

	It("checks the target branch out next to an existing checkout", func(ctx SpecContext) {
		checkout := filepath.Join(GinkgoT().TempDir(), "checkout")
		run(filepath.Dir(checkout), "clone", "--quiet", remote, checkout)

		worktree, err := git.Repository{Dir: checkout}.
			Worktree(ctx, fixtures.Task(), "origin", "next", filepath.Join(GinkgoT().TempDir(), "worktree"))
		Expect(err).NotTo(HaveOccurred())

		Expect(read(filepath.Join(worktree.Dir, "out", "kept.txt"))).To(Equal("kept"))
		Expect(run(worktree.Dir, "rev-parse", "HEAD")).To(Equal(run(checkout, "rev-parse", "origin/next")))
	})

	Describe("Lease", func() {
		It("is empty for a branch the remote does not have", func(ctx SpecContext) {
			repository := clone(ctx, "out")

			Expect(repository.Lease(ctx, fixtures.Task(), remote, "sync/missing")).To(BeEmpty())
		})

		It("is the commit the branch points at", func(ctx SpecContext) {
			repository := clone(ctx, "out")

			Expect(repository.Lease(ctx, fixtures.Task(), remote, "next")).
				To(Equal(run(repository.Dir, "rev-parse", "HEAD")))
		})
	})

	Describe("Push", func() {
		commit := func(repository git.Repository, content string) {
			GinkgoHelper()

			write(filepath.Join(repository.Dir, "out", "kept.txt"), content)
			run(repository.Dir, "commit", "--quiet", "--all", "--message", content)
		}

		It("creates the branch on an empty lease", func(ctx SpecContext) {
			repository := clone(ctx, "out")
			commit(repository, "first")

			Expect(repository.Push(ctx, fixtures.Task(), remote, "sync/branch", "")).To(Succeed())
			Expect(repository.Lease(ctx, fixtures.Task(), remote, "sync/branch")).
				To(Equal(run(repository.Dir, "rev-parse", "HEAD")))
		})

		// another run pushing in between moves the branch off the lease, and its result
		// is the one that has to survive.
		It("refuses to overwrite a branch that moved off the lease", func(ctx SpecContext) {
			first := clone(ctx, "out")
			commit(first, "first")
			Expect(first.Push(ctx, fixtures.Task(), remote, "sync/branch", "")).To(Succeed())

			second := clone(ctx, "out")
			commit(second, "second")
			Expect(second.Push(ctx, fixtures.Task(), remote, "sync/branch", "")).NotTo(Succeed())

			lease, err := second.Lease(ctx, fixtures.Task(), remote, "sync/branch")
			Expect(err).NotTo(HaveOccurred())
			Expect(second.Push(ctx, fixtures.Task(), remote, "sync/branch", lease)).To(Succeed())
		})
	})

	Describe("credentials", func() {
		var runner *tests.TestingCommandRunner
		var t *Task

		BeforeEach(func() {
			runner = fixtures.Runner()
			t = NewTaskList(tests.NewPlumber().Plumber.SetRuntime(Runtime{CommandRunner: runner.Runner()})).CreateTask("git")
		})

		// CI_REPOSITORY_URL carries the job token in its userinfo, and the remote is
		// printed with every command.
		It("refuses a remote that carries credentials of its own", func(ctx SpecContext) {
			repository := git.Repository{Dir: GinkgoT().TempDir(), Token: "glpat-secret"}
			remote := "https://gitlab-ci-token:job-secret@gitlab.example.com/group/project.git"

			Expect(repository.Clone(ctx, t, remote, "next", []string{"out"})).To(MatchError(Not(ContainSubstring("job-secret"))))
			_, err := repository.Lease(ctx, t, remote, "sync/branch")
			Expect(err).To(HaveOccurred())
			Expect(repository.Push(ctx, t, remote, "sync/branch", "")).NotTo(Succeed())

			Expect(runner.Invocations()).To(BeEmpty())
		})

		// the inherited environment wins over the one a command sets, so the header
		// would be dropped and the push would fail as unauthenticated much later.
		It("refuses to run when git config is already passed through the environment", func(ctx SpecContext) {
			GinkgoT().Setenv("GIT_CONFIG_COUNT", "1")
			repository := git.Repository{Dir: GinkgoT().TempDir(), Token: "glpat-secret"}

			Expect(repository.Stage(ctx, t, []string{"out"})).To(MatchError(ContainSubstring("GIT_CONFIG_COUNT")))
			Expect(runner.Invocations()).To(BeEmpty())
		})
	})

	// an argument list shows up in a process listing and in the command log, and a URL
	// lands in the remote config of the clone. Without a username the token goes out as
	// the oauth2 user.
	It("keeps the token out of the arguments of every git call", func(ctx SpecContext) {
		token := "glpat-secret"
		runner := fixtures.Runner()
		p := tests.NewPlumber().Plumber.SetRuntime(Runtime{CommandRunner: runner.Runner()})
		t := NewTaskList(p).CreateTask("git")
		repository := git.Repository{Dir: GinkgoT().TempDir(), Token: token}

		Expect(repository.Clone(ctx, t, "https://gitlab.example.com/group/project.git", "next", []string{"out"})).To(Succeed())
		_, err := repository.Worktree(ctx, t, "origin", "next", GinkgoT().TempDir())
		Expect(err).NotTo(HaveOccurred())
		_, err = repository.Lease(ctx, t, "origin", "sync/branch")
		Expect(err).NotTo(HaveOccurred())
		Expect(repository.Push(ctx, t, "origin", "sync/branch", "")).To(Succeed())
		Expect(repository.Stage(ctx, t, []string{"out"})).To(Succeed())
		Expect(repository.Patch(ctx, t, filepath.Join(GinkgoT().TempDir(), "changes.patch"))).To(Succeed())

		encoded := base64.StdEncoding.EncodeToString([]byte("oauth2:" + token))
		invocations := runner.Invocations()
		Expect(invocations).To(HaveLen(8))

		for _, invocation := range invocations {
			Expect(invocation.Name).To(Equal("git"))
			Expect(strings.Join(invocation.Args, " ")).NotTo(Or(ContainSubstring(token), ContainSubstring(encoded)))
			Expect(invocation.Env).To(ContainElement("GIT_CONFIG_VALUE_0=Authorization: Basic " + encoded))
		}
	})
})
