package build

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"

	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

type fakeNotes struct {
	listed  []*clientgitlab.Note
	created []string
	updated []string
}

var _ gitlab.NotesAdapter = (*fakeNotes)(nil)

func (f *fakeNotes) ListMergeRequestNotes(
	_ any,
	_ int64,
	_ *clientgitlab.ListMergeRequestNotesOptions,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.Note, *clientgitlab.Response, error) {
	return f.listed, nil, nil
}

func (f *fakeNotes) CreateMergeRequestNote(
	_ any,
	_ int64,
	opt *clientgitlab.CreateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	f.created = append(f.created, *opt.Body)

	return &clientgitlab.Note{ID: 1}, nil, nil
}

func (f *fakeNotes) UpdateMergeRequestNote(
	_ any,
	_ int64,
	note int64,
	opt *clientgitlab.UpdateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	f.updated = append(f.updated, *opt.Body)

	return &clientgitlab.Note{ID: note}, nil, nil
}

var _ = Describe("merge request report", func() {
	var repo string
	var base string
	var notes *fakeNotes
	var mergeRequestIid int64
	var diffBase func() string
	var logEnabled bool
	var targetBranchSha string

	git := func(args ...string) string {
		GinkgoHelper()

		// a commit otherwise detaches a maintenance run, which can still be writing into
		// the repository when the temporary directory is removed.
		cmd := exec.Command("git", append([]string{"-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "-c", "user.name=spec", "-c", "user.email=spec@test"}, args...)...)
		cmd.Dir = repo
		out, err := cmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), string(out))

		return strings.TrimSpace(string(out))
	}

	write := func(path string, content string) {
		GinkgoHelper()

		Expect(os.MkdirAll(filepath.Dir(filepath.Join(repo, path)), 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(repo, path), []byte(content), 0o644)).To(Succeed())
	}

	configMap := func(name string, value string) string {
		return "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: " + name + "\n  namespace: apps\ndata:\n  value: " + value + "\n"
	}

	commit := func() string {
		GinkgoHelper()

		git("add", "-A")
		git("commit", "-q", "-m", "commit")

		return git("rev-parse", "HEAD")
	}

	BeforeEach(func() {
		repo = GinkgoT().TempDir()
		notes = &fakeNotes{}
		mergeRequestIid = 7
		diffBase = func() string { return base }
		logEnabled = true
		targetBranchSha = ""

		git("init", "-q", "-b", "main")
		write("app/kustomization.yaml", "resources:\n  - kept.yaml\n  - removed.yaml\n")
		write("app/kept.yaml", configMap("kept", "before"))
		write("app/removed.yaml", configMap("removed", "gone"))
		write("same/kustomization.yaml", "resources:\n  - same.yaml\n")
		write("same/same.yaml", configMap("same", "same"))
		base = commit()

		previous := newNotes
		newNotes = func(_ gitlab.MergeRequestReportConfig) (gitlab.NotesAdapter, error) { return notes, nil }
		DeferCleanup(func() { newNotes = previous })
	})

	run := func(overlays ...string) error {
		GinkgoHelper()

		*P = Pipe{
			LoadRestrictor: string(LoadRestrictorNone),
			// the pipe sees both shas of a detached merge request pipeline, the target
			// branch one set but empty.
			TargetBranchSha: targetBranchSha,
			DiffBaseSha:     diffBase(),
		}
		P.ReportLog.Enabled = logEnabled
		P.MergeRequestReport = gitlab.MergeRequestReportConfig{
			Enabled:         true,
			Token:           "glpat-report",
			ApiUrl:          "https://gitlab.test/api/v4",
			ProjectId:       "group/project",
			MergeRequestIid: mergeRequestIid,
		}
		P.ReportMetadata.JobName = "kustomize"
		*C = Ctx{}
		setup.C.Cwd = repo
		setup.C.Overlays = []string{}
		for _, overlay := range overlays {
			setup.C.Overlays = append(setup.C.Overlays, filepath.Join(repo, overlay))
		}

		return tests.NewTaskListCli(tests.TaskListCli{
			AppName:     "pipe-kustomize",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	It("notes the rendered resources the change adds, changes and removes per overlay", func() {
		write("app/kustomization.yaml", "resources:\n  - kept.yaml\n  - added.yaml\n")
		write("app/kept.yaml", configMap("kept", "after"))
		Expect(os.Remove(filepath.Join(repo, "app/removed.yaml"))).To(Succeed())
		write("app/added.yaml", configMap("added", "new"))
		write("fresh/kustomization.yaml", "resources:\n  - fresh.yaml\n")
		write("fresh/fresh.yaml", configMap("fresh", "new"))
		commit()

		Expect(run("app", "fresh", "same")).To(Succeed())

		Expect(notes.updated).To(BeEmpty())
		Expect(notes.created).To(HaveLen(1))
		body := notes.created[0]

		Expect(body).To(ContainSubstring("## Rendered manifests against `" + base[:8] + "`"))
		Expect(body).To(ContainSubstring("Job: kustomize"))
		Expect(body).To(ContainSubstring("| `app` | 1 | 1 | 1 |"))
		Expect(body).To(ContainSubstring("| `fresh` | 1 | 0 | 0 |"))
		Expect(body).To(ContainSubstring("| `same` | 0 | 0 | 0 |"))
		Expect(body).To(ContainSubstring("4 manifest(s) changed."))
		Expect(body).To(ContainSubstring("\n`app/apps/ConfigMap.core/kept.yaml`\n\n```diff\n@@"))
		Expect(body).To(ContainSubstring("-  value: before\n+  value: after\n"))
		Expect(body).To(ContainSubstring("\n`app/apps/ConfigMap.core/removed.yaml`\n\n```diff\ndeleted file mode"))
		Expect(body).NotTo(ContainSubstring("same.yaml"))
		Expect(body).To(ContainSubstring("<!-- gitlab-pipes:mr-report:kustomize:"))

		Expect(git("worktree", "list", "--porcelain")).NotTo(ContainSubstring("kustomize-base-"))
	})

	It("notes a change that renders the same manifests", func() {
		write("same/README.md", "not rendered\n")
		commit()

		Expect(run("app", "same")).To(Succeed())

		Expect(notes.updated).To(BeEmpty())
		Expect(notes.created).To(HaveLen(1))
		Expect(notes.created[0]).To(ContainSubstring("| `same` | 0 | 0 | 0 |"))
		Expect(notes.created[0]).To(ContainSubstring("No changes."))
	})

	It("updates the note an earlier push left in place", func() {
		write("same/README.md", "not rendered\n")
		commit()
		notes.listed = []*clientgitlab.Note{{ID: 3, Body: "earlier\n\n<!-- gitlab-pipes:mr-report:kustomize:" + repo + " -->"}}

		Expect(run("app", "same")).To(Succeed())

		Expect(notes.created).To(BeEmpty())
		Expect(notes.updated).To(HaveLen(1))
		Expect(notes.updated[0]).To(ContainSubstring("No changes."))
	})

	It("names an overlay the base can not render and still opens the note", func() {
		write("app/kustomization.yaml", "resources:\n  - missing.yaml\n")
		base = commit()
		write("app/kustomization.yaml", "resources:\n  - kept.yaml\n")
		commit()

		Expect(run("app")).To(Succeed())

		Expect(notes.created).To(HaveLen(1))
		Expect(notes.created[0]).To(ContainSubstring("| can not render the base |"))
		Expect(notes.created[0]).To(ContainSubstring("can not be rendered on the base, so it is left out of the diff"))
	})

	It("reports in the job log alone when the merge request note is disabled", func() {
		write("app/kept.yaml", configMap("kept", "after"))
		commit()
		mergeRequestIid = 0

		Expect(run("app")).To(Succeed())

		Expect(notes.created).To(BeEmpty())
		Expect(notes.updated).To(BeEmpty())
		Expect(git("worktree", "list", "--porcelain")).NotTo(ContainSubstring("kustomize-base-"))
	})

	It("diffs a merged results pipeline against the target branch", func() {
		write("app/kept.yaml", configMap("kept", "after"))
		target := commit()
		write("app/kept.yaml", configMap("kept", "later"))
		commit()
		targetBranchSha = target

		Expect(run("app")).To(Succeed())

		Expect(notes.created).To(HaveLen(1))
		Expect(notes.created[0]).To(ContainSubstring("## Rendered manifests against `" + target[:8] + "`"))
		Expect(notes.created[0]).To(ContainSubstring("-  value: after\n+  value: later\n"))
	})

	It("stays out of a pipeline that has no base to diff against", func() {
		write("app/kept.yaml", configMap("kept", "after"))
		commit()
		diffBase = func() string { return "" }

		Expect(run("app")).To(Succeed())

		Expect(notes.created).To(BeEmpty())
		Expect(git("worktree", "list", "--porcelain")).NotTo(ContainSubstring("kustomize-base-"))
	})
})
