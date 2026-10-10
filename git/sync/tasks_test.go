package sync

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"

	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

type fakeNotes struct {
	bodies map[int64]string
}

var _ gitlab.NotesAdapter = (*fakeNotes)(nil)

func (f *fakeNotes) ListMergeRequestNotes(
	_ any,
	_ int64,
	_ *clientgitlab.ListMergeRequestNotesOptions,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.Note, *clientgitlab.Response, error) {
	return nil, nil, nil
}

func (f *fakeNotes) CreateMergeRequestNote(
	_ any,
	mergeRequest int64,
	opt *clientgitlab.CreateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	f.bodies[mergeRequest] = *opt.Body

	return &clientgitlab.Note{ID: mergeRequest}, nil, nil
}

func (f *fakeNotes) UpdateMergeRequestNote(
	_ any,
	_ int64,
	_ int64,
	_ *clientgitlab.UpdateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	Fail("no note is listed, so none can be updated")

	return nil, nil, nil
}

// keeps the merge requests GitLab would hold, so a run that finds one open updates it
// and a create racing another run can be answered the way GitLab answers it.
type fakeMergeRequests struct {
	open     []*clientgitlab.BasicMergeRequest
	listed   []*clientgitlab.ListProjectMergeRequestsOptions
	created  []*clientgitlab.CreateMergeRequestOptions
	updated  []*clientgitlab.UpdateMergeRequestOptions
	conflict bool
}

var _ gitlab.MergeRequestsAdapter = (*fakeMergeRequests)(nil)

func (f *fakeMergeRequests) ListProjectMergeRequests(
	_ any,
	opt *clientgitlab.ListProjectMergeRequestsOptions,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.BasicMergeRequest, *clientgitlab.Response, error) {
	f.listed = append(f.listed, opt)

	return f.open, &clientgitlab.Response{}, nil
}

func (f *fakeMergeRequests) CreateMergeRequest(
	_ any,
	opt *clientgitlab.CreateMergeRequestOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.MergeRequest, *clientgitlab.Response, error) {
	f.created = append(f.created, opt)

	if f.conflict {
		f.conflict = false
		f.open = append(f.open, &clientgitlab.BasicMergeRequest{IID: 42, State: "opened"})

		return nil, nil, &clientgitlab.ErrorResponse{StatusCode: http.StatusConflict, Message: "Another open merge request already exists for this source branch"}
	}

	f.open = append(f.open, &clientgitlab.BasicMergeRequest{IID: 1, State: "opened"})

	return &clientgitlab.MergeRequest{BasicMergeRequest: clientgitlab.BasicMergeRequest{IID: 1, WebURL: "https://gitlab.test/-/merge_requests/1"}}, nil, nil
}

func (f *fakeMergeRequests) UpdateMergeRequest(
	_ any,
	mergeRequest int64,
	opt *clientgitlab.UpdateMergeRequestOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.MergeRequest, *clientgitlab.Response, error) {
	f.updated = append(f.updated, opt)

	return &clientgitlab.MergeRequest{BasicMergeRequest: clientgitlab.BasicMergeRequest{IID: mergeRequest}}, nil, nil
}

// runs every command for real and records it, with a hook to act on the remote
// right before a command does.
type passthrough struct {
	invocations []CommandInvocation
	before      func(invocation CommandInvocation)
}

func (r *passthrough) Run(ctx context.Context, invocation CommandInvocation, runtime CommandRuntime) (CommandResult, error) {
	r.invocations = append(r.invocations, invocation)

	if r.before != nil {
		r.before(invocation)
	}

	return NewCommandRunner().Run(ctx, invocation, runtime)
}

func (r *passthrough) ran(args ...string) bool {
	return slices.ContainsFunc(r.invocations, func(invocation CommandInvocation) bool {
		return len(invocation.Args) >= len(args) && slices.Equal(invocation.Args[:len(args)], args)
	})
}

var _ = Describe("Git sync tasks", func() {
	var base, source, checkout string
	var runner *passthrough
	var notes *fakeNotes
	var mergeRequests *fakeMergeRequests

	invoke := func(dir string, args ...string) string {
		GinkgoHelper()

		command := exec.Command("git", args...)
		command.Dir = dir
		command.Env = append(command.Environ(), "GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com", "GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com")

		output, err := command.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), string(output))

		return strings.TrimSpace(string(output))
	}

	write := func(path string, content string) {
		GinkgoHelper()

		Expect(os.MkdirAll(filepath.Dir(path), 0o755)).To(Succeed())
		Expect(os.WriteFile(path, []byte(content), 0o644)).To(Succeed())
	}

	remote := func(project string) string {
		return filepath.Join(base, project+".git")
	}

	// a bare project on the remote with the default branch the pipeline runs for and the
	// target branch the sync writes into, reached over file:// like a server.
	project := func(path string) {
		GinkgoHelper()

		seed := filepath.Join(GinkgoT().TempDir(), "seed")

		invoke(base, "init", "--quiet", "--bare", "--initial-branch", "main", remote(path))
		invoke(remote(path), "config", "uploadpack.allowFilter", "true")
		invoke(base, "init", "--quiet", "--initial-branch", "main", seed)

		write(filepath.Join(seed, "README.md"), "readme")
		invoke(seed, "add", "--all")
		invoke(seed, "commit", "--quiet", "--message", "main")
		invoke(seed, "push", "--quiet", remote(path), "main")

		invoke(seed, "checkout", "--quiet", "-b", "next")
		write(filepath.Join(seed, "out", "kept.txt"), "kept")
		write(filepath.Join(seed, "out", "removed.txt"), "removed")
		write(filepath.Join(seed, "other", "untouched.txt"), "untouched")
		invoke(seed, "add", "--all")
		invoke(seed, "commit", "--quiet", "--message", "next")
		invoke(seed, "push", "--quiet", remote(path), "next")
	}

	head := func(project string, branch string) string {
		GinkgoHelper()

		output := invoke(base, "ls-remote", "--heads", remote(project), "refs/heads/"+branch)
		sha, _, _ := strings.Cut(output, "\t")

		return sha
	}

	BeforeEach(func() {
		// a git hook runs with the variables that locate the repository it fired in,
		// which every git call of a spec would otherwise act on instead of its own.
		for _, name := range strings.Fields(invoke("", "rev-parse", "--local-env-vars")) {
			if value, ok := os.LookupEnv(name); ok {
				Expect(os.Unsetenv(name)).To(Succeed())
				DeferCleanup(os.Setenv, name, value)
			}
		}

		GinkgoT().Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
		GinkgoT().Setenv("GIT_CONFIG_NOSYSTEM", "1")

		base = GinkgoT().TempDir()
		project("group/project")

		checkout = filepath.Join(base, "checkout")
		invoke(base, "clone", "--quiet", "--branch", "main", remote("group/project"), checkout)

		source = GinkgoT().TempDir()
		write(filepath.Join(source, "kept.txt"), "changed")

		runner = &passthrough{}
		notes = &fakeNotes{bodies: map[int64]string{}}
		mergeRequests = &fakeMergeRequests{}

		previousNotes, previousMergeRequests := newNotes, newMergeRequests
		newNotes = func(_ gitlab.MergeRequestReportConfig) (gitlab.NotesAdapter, error) { return notes, nil }
		newMergeRequests = func(_ gitlab.MergeRequestConfig) (gitlab.MergeRequestsAdapter, error) { return mergeRequests, nil }
		DeferCleanup(func() {
			newNotes, newMergeRequests = previousNotes, previousMergeRequests
		})

		*C = Ctx{}
	})

	pipe := func(mode string) Pipe {
		return Pipe{
			Mode:          mode,
			Target:        Target{Project: "group/project", Branch: "next"},
			Paths:         []git.Path{{Source: source, Destination: "out"}},
			Branch:        "sync/group-project-generate",
			CommitMessage: "chore: sync generated files\n\nGenerated by the pipeline.",
			Patch:         filepath.Join(GinkgoT().TempDir(), "git-sync.patch"),
			Token:         "glpat-test",
			Author:        Identity{Name: "Pipes", Email: "gitlab+pipes@kilic.dev"},
			Committer:     Identity{Name: "Committer", Email: "committer@kilic.dev"},
			Project: Project{
				Path:          "group/project",
				Dir:           checkout,
				ServerUrl:     "file://" + base,
				DefaultBranch: "main",
			},
			Refs: git.Refs{Branch: "main"},
			MergeRequestReport: gitlab.MergeRequestReportConfig{
				Enabled:   true,
				Token:     "glpat-report",
				ApiUrl:    "https://gitlab.test/api/v4",
				ProjectId: "group/project",
			},
			ReportMetadata: terraform.Metadata{JobName: "generate", CommitSha: head("group/project", "main")},
		}
	}

	run := func(p Pipe) error {
		GinkgoHelper()

		*P = p

		return tests.NewTaskListCli(tests.TaskListCli{
			AppName:     "pipe-git",
			CommandName: "sync",
			Runtime:     Runtime{CommandRunner: runner},
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	Describe("report", func() {
		It("reports the changes on the merge request of the pipeline without pushing", func() {
			p := pipe(ModeReport)
			p.Refs.Branch = "feature"
			p.MergeRequestReport.MergeRequestIid = 7

			Expect(run(p)).To(Succeed())

			Expect(runner.ran("worktree", "add")).To(BeTrue())
			Expect(runner.ran("push")).To(BeFalse())
			Expect(head("group/project", p.Branch)).To(BeEmpty())

			Expect(notes.bodies).To(HaveKey(int64(7)))
			Expect(notes.bodies[7]).To(ContainSubstring("## Sync into `group/project@next`"))
			Expect(notes.bodies[7]).To(ContainSubstring("-removed"))
			Expect(notes.bodies[7]).To(ContainSubstring("+changed"))

			patch, err := os.ReadFile(p.Patch)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(patch)).To(ContainSubstring("diff --git a/out/removed.txt b/out/removed.txt"))

			Expect(filepath.Join(checkout, "out")).NotTo(BeAnExistingFile())
			Expect(mergeRequests.created).To(BeEmpty())
		})

		It("clones a target in another project with only the destinations", func() {
			project("group/other")

			p := pipe(ModeReport)
			p.Target.Project = "group/other"
			p.MergeRequestReport.MergeRequestIid = 7

			Expect(run(p)).To(Succeed())

			Expect(runner.ran("clone")).To(BeTrue())
			Expect(runner.ran("worktree", "add")).To(BeFalse())
			Expect(notes.bodies[7]).To(ContainSubstring("## Sync into `group/other@next`"))
			Expect(notes.bodies[7]).To(ContainSubstring("+changed"))
		})

		DescribeTable(
			"posts the report on the merge requests open from the branch",
			func(open string, listed []int64, expected []int64) {
				for _, iid := range listed {
					mergeRequests.open = append(mergeRequests.open, &clientgitlab.BasicMergeRequest{IID: iid})
				}

				p := pipe(ModeReport)
				p.Refs.Branch = "feature"
				p.Project.OpenMergeRequests = open

				Expect(run(p)).To(Succeed())

				reported := []int64{}
				for iid := range notes.bodies {
					reported = append(reported, iid)
				}

				Expect(reported).To(ConsistOf(expected))
			},
			Entry("none", "", nil, []int64{}),
			Entry("the ones the pipeline lists", "group/project!1,group/project!2", nil, []int64{1, 2}),
			Entry(
				"the ones GitLab lists once the pipeline may have cut the list",
				"group/project!1,group/project!2,group/project!3,group/project!4",
				[]int64{1, 2, 3, 4, 5},
				[]int64{1, 2, 3, 4, 5},
			),
		)

		It("lists the merge requests open from the branch of the pipeline once the list may be cut", func() {
			p := pipe(ModeReport)
			p.Refs.Branch = "feature"
			p.Project.OpenMergeRequests = "group/project!1,group/project!2,group/project!3,group/project!4"

			Expect(run(p)).To(Succeed())

			Expect(mergeRequests.listed).To(HaveLen(1))
			Expect(*mergeRequests.listed[0].SourceBranch).To(Equal("feature"))
			Expect(*mergeRequests.listed[0].State).To(Equal("opened"))
		})

		It("refuses an open merge request it can not parse", func() {
			p := pipe(ModeReport)
			p.Project.OpenMergeRequests = "group/project#1"

			Expect(run(p)).To(MatchError(ContainSubstring("Can not parse the open merge request: group/project#1")))
		})
	})

	It("skips the sync when the pipeline runs on the target branch itself", func() {
		p := pipe(ModePublish)
		p.Refs.Branch = "next"

		Expect(run(p)).To(Succeed())

		Expect(runner.invocations).To(BeEmpty())
		Expect(notes.bodies).To(BeEmpty())
	})

	It("reports no changes and leaves the merge request alone when the target already matches", func() {
		write(filepath.Join(source, "kept.txt"), "kept")
		write(filepath.Join(source, "removed.txt"), "removed")

		p := pipe(ModePublish)
		p.MergeRequestReport.MergeRequestIid = 7

		Expect(run(p)).To(Succeed())

		Expect(notes.bodies[7]).To(ContainSubstring("No changes."))
		Expect(runner.ran("commit")).To(BeFalse())
		Expect(runner.ran("push")).To(BeFalse())
		Expect(mergeRequests.listed).To(BeEmpty())
	})

	Describe("publish", func() {
		It("commits onto the tip of the target with its own identity and opens a merge request", func() {
			GinkgoT().Setenv("GIT_AUTHOR_NAME", "job")
			GinkgoT().Setenv("GIT_COMMITTER_NAME", "job")

			p := pipe(ModePublish)
			tip := head("group/project", "next")

			Expect(run(p)).To(Succeed())

			pushed := head("group/project", p.Branch)
			Expect(pushed).NotTo(BeEmpty())

			invoke(checkout, "fetch", "--quiet", "origin", p.Branch)
			Expect(invoke(checkout, "rev-parse", "FETCH_HEAD^")).To(Equal(tip))
			Expect(invoke(checkout, "log", "-1", "--format=%an <%ae>|%cn <%ce>|%B", "FETCH_HEAD")).
				To(Equal("Pipes <gitlab+pipes@kilic.dev>|Committer <committer@kilic.dev>|chore: sync generated files\n\nGenerated by the pipeline."))
			Expect(strings.Fields(invoke(checkout, "diff", "--name-status", tip, "FETCH_HEAD"))).
				To(Equal([]string{"M", "out/kept.txt", "D", "out/removed.txt"}))

			Expect(mergeRequests.created).To(HaveLen(1))
			Expect(*mergeRequests.created[0].Title).To(Equal("chore: sync generated files"))
			Expect(*mergeRequests.created[0].SourceBranch).To(Equal(p.Branch))
			Expect(*mergeRequests.created[0].TargetBranch).To(Equal("next"))
			Expect(*mergeRequests.created[0].Description).To(ContainSubstring("+changed"))
		})

		It("replaces what an earlier run pushed and updates the merge request it opened", func() {
			Expect(run(pipe(ModePublish))).To(Succeed())
			first := head("group/project", "sync/group-project-generate")

			write(filepath.Join(source, "kept.txt"), "changed again")
			*C = Ctx{}

			Expect(run(pipe(ModePublish))).To(Succeed())

			Expect(head("group/project", "sync/group-project-generate")).NotTo(Or(BeEmpty(), Equal(first)))
			Expect(mergeRequests.created).To(HaveLen(1))
			Expect(mergeRequests.updated).To(HaveLen(1))
			Expect(*mergeRequests.updated[0].Description).To(ContainSubstring("+changed again"))
		})

		// another run pushing between the lease and the push moves the branch off the
		// lease, and what that run pushed is what has to survive.
		It("refuses to overwrite a branch another run pushed to in the meantime", func() {
			p := pipe(ModePublish)
			racer := filepath.Join(GinkgoT().TempDir(), "racer")

			runner.before = func(invocation CommandInvocation) {
				if len(invocation.Args) == 0 || invocation.Args[0] != "push" {
					return
				}

				invoke(base, "clone", "--quiet", "--branch", "next", remote("group/project"), racer)
				write(filepath.Join(racer, "out", "kept.txt"), "racer")
				invoke(racer, "commit", "--quiet", "--all", "--message", "racer")
				invoke(racer, "push", "--quiet", "origin", "HEAD:"+p.Branch)
			}

			Expect(run(p)).To(MatchError(ContainSubstring("Can not push the branch")))

			Expect(head("group/project", p.Branch)).To(Equal(invoke(racer, "rev-parse", "HEAD")))
			Expect(mergeRequests.created).To(BeEmpty())
		})

		It("updates the merge request another run opened between the listing and the creation", func() {
			mergeRequests.conflict = true

			Expect(run(pipe(ModePublish))).To(Succeed())

			Expect(mergeRequests.created).To(HaveLen(1))
			Expect(mergeRequests.listed).To(HaveLen(2))
			Expect(mergeRequests.updated).To(HaveLen(1))
		})

		It("fails the publish once the commit of the pipeline is no longer the head of the default branch", func() {
			p := pipe(ModePublish)
			p.ReportMetadata.CommitSha = strings.Repeat("0", 40)
			p.MergeRequestReport.MergeRequestIid = 7

			Expect(run(p)).To(MatchError(ContainSubstring("no longer the head")))

			Expect(notes.bodies[7]).To(ContainSubstring("+changed"))
			Expect(runner.ran("commit")).To(BeFalse())
			Expect(head("group/project", p.Branch)).To(BeEmpty())
			Expect(mergeRequests.created).To(BeEmpty())
		})

		It("needs a token to publish", func() {
			p := pipe(ModePublish)
			p.Token = ""

			Expect(run(p)).To(MatchError(ContainSubstring("Publishing needs a GitLab token")))
			Expect(runner.invocations).To(BeEmpty())
		})

		It("credits the author of the pipeline commit as a co-author while committing as itself", func() {
			p := pipe(ModePublish)
			p.Project.CommitAuthor = "Jane Doe <jane@example.com>"

			Expect(run(p)).To(Succeed())

			invoke(checkout, "fetch", "--quiet", "origin", p.Branch)
			Expect(invoke(checkout, "log", "-1", "--format=%an <%ae>|%B", "FETCH_HEAD")).
				To(Equal("Pipes <gitlab+pipes@kilic.dev>|chore: sync generated files\n\nGenerated by the pipeline.\n\nCo-authored-by: Jane Doe <jane@example.com>"))
		})

		DescribeTable("refuses to publish without the identity of the commit, naming what is missing",
			func(unset func(p *Pipe), variable string) {
				p := pipe(ModePublish)
				unset(&p)

				Expect(run(p)).To(MatchError(ContainSubstring(variable + " is not set")))
				Expect(runner.invocations).To(BeEmpty())
			},
			Entry("the author name", func(p *Pipe) { p.Author.Name = "" }, "GIT_PIPES_AUTHOR_NAME"),
			Entry("the author email", func(p *Pipe) { p.Author.Email = "" }, "GIT_PIPES_AUTHOR_EMAIL"),
			Entry("the committer name", func(p *Pipe) { p.Committer.Name = "" }, "GIT_PIPES_COMMITTER_NAME"),
			Entry("the committer email", func(p *Pipe) { p.Committer.Email = "" }, "GIT_PIPES_COMMITTER_EMAIL"),
		)

		It("reports without the identity of the commit, since a report never commits", func() {
			p := pipe(ModeReport)
			p.Author, p.Committer = Identity{}, Identity{}
			p.MergeRequestReport.MergeRequestIid = 7

			Expect(run(p)).To(Succeed())
			Expect(notes.bodies[7]).To(ContainSubstring("+changed"))
		})
	})
})
