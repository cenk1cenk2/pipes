package sync

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

var _ = Describe("Report", func() {
	BeforeEach(func() {
		*P = Pipe{
			Target: Target{Project: "group/project", Branch: "next"},
			Patch:  "git-sync.patch",
			ReportMetadata: terraform.Metadata{
				JobName:        "generate",
				JobUrl:         "https://gitlab.test/group/project/-/jobs/1",
				PipelineId:     "2",
				PipelineUrl:    "https://gitlab.test/group/project/-/pipelines/2",
				CommitShortSha: "abcdef01",
			},
		}
		*C = Ctx{}
	})

	It("points back at the job and the pipeline", func() {
		Expect(
			render(),
		).To(ContainSubstring("Job: [generate](https://gitlab.test/group/project/-/jobs/1) | Pipeline: [#2](https://gitlab.test/group/project/-/pipelines/2) | Commit: `abcdef01`"))
	})

	It("counts the changed files", func() {
		C.Diff = "diff --git a/a b/a\n+a\ndiff --git a/b b/b\n-b\n"

		Expect(render()).To(ContainSubstring("2 file(s) changed."))
	})

	// a context line of a changed markdown file is indented by a single space, which
	// still closes a fence of the same length.
	It("fences the diff longer than any fence it carries", func() {
		C.Diff = "diff --git a/README.md b/README.md\n ````go\n+x\n ````\n"

		Expect(render()).To(ContainSubstring("`````diff\n" + C.Diff + "`````\n"))
	})

	It("cuts a diff GitLab would refuse at a line, and points at the full patch", func() {
		line := "+" + strings.Repeat("x", 99) + "\n"
		C.Diff = "diff --git a/a b/a\n" + strings.Repeat(line, reportDiffLimit/len(line)+10)

		body := render()

		Expect(len(body)).To(BeNumerically("<", 1_000_000))
		Expect(body).To(ContainSubstring(line + "```\n"))
		Expect(
			body,
		).To(ContainSubstring("The diff is truncated, the full patch `git-sync.patch` is in [the artifacts of the job](https://gitlab.test/group/project/-/jobs/1/artifacts/browse)."))
	})

	It("cuts a diff without a line break at the limit instead of dropping it", func() {
		C.Diff = strings.Repeat("x", reportDiffLimit+10)

		Expect(render()).To(ContainSubstring("```diff\n" + strings.Repeat("x", reportDiffLimit) + "```\n"))
	})
})
