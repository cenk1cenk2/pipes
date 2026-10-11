package sync

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

var _ = Describe("Report", func() {
	readFixture := func(name string) string {
		GinkgoHelper()

		data, err := os.ReadFile(filepath.Join("testdata", name))
		Expect(err).NotTo(HaveOccurred())

		return string(data)
	}

	BeforeEach(func() {
		*P = Pipe{
			Target: Target{Project: "group/project", Branch: "main"},
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

	Context("with a patch of several files", func() {
		BeforeEach(func() {
			C.Diff = readFixture("files.patch")
		})

		It("renders every file as a section of its own under its name", func() {
			Expect(render()).To(Equal(readFixture("files.md")))
		})

		It("leaves out the header lines the name above the section already says", func() {
			body := render()

			Expect(body).NotTo(ContainSubstring("diff --git "))
			Expect(body).NotTo(MatchRegexp(`(?m)^index `))
			Expect(body).NotTo(MatchRegexp(`(?m)^--- (a/|/dev/null)`))
			Expect(body).NotTo(MatchRegexp(`(?m)^\+\+\+ (b/|/dev/null)`))
		})

		It("keeps the rename lines git marks under the name the file ends up with", func() {
			Expect(render()).To(ContainSubstring(
				"\n`system/namespace-legacy.yaml`\n\n```diff\n" +
					"similarity index 100%\nrename from system/namespace.yaml\nrename to system/namespace-legacy.yaml\n```\n",
			))
		})

		It("keeps the binary patch of a binary file under its name", func() {
			Expect(render()).To(ContainSubstring("\n`system/logo.png`\n\n```diff\nnew file mode 100644\nGIT binary patch\nliteral 10\n"))
		})

		// a context line of a changed markdown file is indented by a single space, which
		// still closes a fence of the same length.
		It("fences the diff of a file longer than any fence it carries", func() {
			Expect(render()).To(ContainSubstring("\n`README.md`\n\n`````diff\n@@ -1,5 +1,5 @@\n"))
		})
	})

	Context("with a patch past the limit", func() {
		BeforeEach(func() {
			C.Diff = readFixture("cap.patch")

			previous := reportDiffLimit
			reportDiffLimit = 186
			DeferCleanup(func() {
				reportDiffLimit = previous
			})
		})

		It("shows the files that fit, cuts the next one at a line and counts the rest", func() {
			Expect(render()).To(Equal(readFixture("cap.md")))
		})

		It("points at the full patch in the artifacts of the job", func() {
			Expect(render()).To(HaveSuffix(
				"\n2 more file(s) are not shown.\n" +
					"\nThe diff is truncated, the full patch `git-sync.patch` is in [the artifacts of the job](https://gitlab.test/group/project/-/jobs/1/artifacts/browse).\n",
			))
		})
	})

	Context("with the branch of the pipeline to compare against", func() {
		const a = "diff --git a/a b/a\n+a\n"
		const b = "diff --git a/b b/b\n+b\n"

		BeforeEach(func() {
			P.Refs.Branch = "feature"
		})

		DescribeTable("groups the files by whether the branch already generated them",
			func(diff string, pending map[string]bool, expected string) {
				C.Diff = diff
				C.Pending = pending

				Expect(render()).To(HaveSuffix(expected))
			},
			Entry("leading with what the branch still lacks", a+b, map[string]bool{"b": true},
				"### Not yet generated on `feature`\n\n1 file(s) changed.\n\n`b`\n\n```diff\n+b\n```\n"+
					"\n### Already generated on `feature`\n\n1 file(s) changed.\n\n`a`\n\n```diff\n+a\n```\n",
			),
			Entry("as one list once the branch generated none of them", a+b, map[string]bool{"a": true, "b": true},
				"Commit: `abcdef01`\n\n2 file(s) changed.\n\n`a`\n\n```diff\n+a\n```\n\n`b`\n\n```diff\n+b\n```\n",
			),
			Entry("saying the branch lacks nothing once it generated all of them", a+b, map[string]bool{},
				"### Not yet generated on `feature`\n\nNo changes beyond what `feature` already carries.\n"+
					"\n### Already generated on `feature`\n\n2 file(s) changed.\n\n`a`\n\n```diff\n+a\n```\n\n`b`\n\n```diff\n+b\n```\n",
			),
			Entry("as one list without a branch to compare against", a+b, nil,
				"Commit: `abcdef01`\n\n2 file(s) changed.\n\n`a`\n\n```diff\n+a\n```\n\n`b`\n\n```diff\n+b\n```\n",
			),
			Entry("with a rename the branch lacks the old name of", "diff --git a/old b/new\nrename from old\nrename to new\n"+b, map[string]bool{"old": true},
				"### Not yet generated on `feature`\n\n1 file(s) changed.\n\n`new`\n\n```diff\nrename from old\nrename to new\n```\n"+
					"\n### Already generated on `feature`\n\n1 file(s) changed.\n\n`b`\n\n```diff\n+b\n```\n",
			),
			Entry("with a quoted name that can not be matched", "diff --git \"a/\\303\\251\" \"b/\\303\\251\"\n+e\n"+b, map[string]bool{},
				"### Not yet generated on `feature`\n\n1 file(s) changed.\n\n`\"a/\\303\\251\" \"b/\\303\\251\"`\n\n```diff\n+e\n```\n"+
					"\n### Already generated on `feature`\n\n1 file(s) changed.\n\n`b`\n\n```diff\n+b\n```\n",
			),
		)

		It("names the branch generically once the pipeline has none", func() {
			P.Refs.Branch = ""
			C.Diff = a + b
			C.Pending = map[string]bool{"b": true}

			Expect(render()).To(ContainSubstring("### Not yet generated on the branch\n"))
			Expect(render()).To(ContainSubstring("### Already generated on the branch\n"))
		})

		Context("past the limit", func() {
			BeforeEach(func() {
				C.Diff = readFixture("cap.patch")

				previous := reportDiffLimit
				reportDiffLimit = 186
				DeferCleanup(func() {
					reportDiffLimit = previous
				})
			})

			It("gives what the branch lacks the limit first and counts the rest of each group", func() {
				C.Pending = map[string]bool{"many/file-1.yaml": true}

				Expect(render()).To(Equal(readFixture("cap-groups.md")))
			})

			It("shows nothing of what the branch already generated once what it lacks is cut", func() {
				C.Pending = map[string]bool{"many/file-1.yaml": true, "many/file-2.yaml": true, "many/file-3.yaml": true}

				body := render()

				Expect(body).To(ContainSubstring("\n### Already generated on `feature`\n\n2 file(s) changed.\n\n2 more file(s) are not shown.\n"))
				Expect(body).NotTo(ContainSubstring("`many/file-4.yaml`"))
				Expect(strings.Count(body, "The diff is truncated")).To(Equal(1))
			})

			It("points at the full patch once only what the branch already generated is cut", func() {
				C.Pending = map[string]bool{"many/file-1.yaml": true}

				body := render()

				Expect(body).To(ContainSubstring("### Not yet generated on `feature`\n\n1 file(s) changed.\n\n`many/file-1.yaml`\n"))
				Expect(
					body,
				).To(HaveSuffix("The diff is truncated, the full patch `git-sync.patch` is in [the artifacts of the job](https://gitlab.test/group/project/-/jobs/1/artifacts/browse).\n"))
			})
		})
	})

	Context("at the size GitLab accepts", func() {
		It("cuts a file GitLab would refuse at a line", func() {
			line := "+" + strings.Repeat("x", 99) + "\n"
			C.Diff = "diff --git a/a b/a\n" + strings.Repeat(line, reportDiffLimit/len(line)+10) + "diff --git a/b b/b\n+b\n"

			body := render()

			Expect(len(body)).To(BeNumerically("<", 1_000_000))
			Expect(body).To(ContainSubstring(line + "```\n"))
			Expect(body).NotTo(ContainSubstring("`b`"))
			Expect(body).To(ContainSubstring("1 more file(s) are not shown."))
		})

		It("counts the sections of many small files against the limit", func() {
			C.Diff = strings.Repeat("diff --git a/file b/file\n+x\n", reportDiffLimit/16)

			body := render()

			Expect(len(body)).To(BeNumerically("<", 1_000_000))
			Expect(body).To(MatchRegexp(`\n\d+ more file\(s\) are not shown\.\n`))
			Expect(body).To(ContainSubstring("The diff is truncated"))
		})
	})
})
