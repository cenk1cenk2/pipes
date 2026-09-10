package markdown_test

import (
	"bytes"
	"log/slog"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/markdown"
)

var _ = Describe("Markdown", func() {
	// the shape the pipes report with: a heading per section, a table of metadata that
	// carries a link, and the resources as a list.
	document := strings.TrimSpace(`
## Terraform plan report

Only names, never values.

### Metadata

| Field | Value |
| --- | --- |
| State | ` + "`production`" + ` |
| Job | [terraform:plan](https://gitlab.example.test/project/-/jobs/1) |

### Summary

| Action | Resources |
| --- | ---: |
| ` + "`create`" + ` | 2 |

### Resources

#### ` + "`create`" + ` (2)

- ` + "`aws_s3_bucket.this`" + `
- ` + "`aws_iam_role.this`" + ` (moved from ` + "`aws_iam_role.old`" + `)
`)

	escapes := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	plain := func(body string) string {
		return escapes.ReplaceAllString(body, "")
	}

	It("carries the document through without its markup", func() {
		body, err := markdown.Render(document)
		Expect(err).NotTo(HaveOccurred())

		Expect(plain(body)).To(SatisfyAll(
			ContainSubstring("Terraform plan report"),
			ContainSubstring("production"),
			ContainSubstring("aws_iam_role.this"),
			ContainSubstring("moved from"),
			Not(ContainSubstring("##")),
			Not(ContainSubstring("`")),
		))
	})

	// the footnotes a table would collect are truncated against the word wrap, which the
	// renderer turns off, so a link that leaves the cell leaves its url behind with it.
	It("keeps the url of a link that sits in a table", func() {
		body, err := markdown.Render(document)
		Expect(err).NotTo(HaveOccurred())

		Expect(plain(body)).To(ContainSubstring("https://gitlab.example.test/project/-/jobs/1"))
	})

	It("leaves no line padded out to the width of the document", func() {
		body, err := markdown.Render(document)
		Expect(err).NotTo(HaveOccurred())

		for line := range strings.SplitSeq(body, "\n") {
			Expect(plain(line)).To(Equal(strings.TrimRight(plain(line), " \t")))
		}
	})

	// a line that ends mid style bleeds its color into every line the log viewer prints
	// after it.
	It("closes every styled line", func() {
		body, err := markdown.Render(document)
		Expect(err).NotTo(HaveOccurred())

		for line := range strings.SplitSeq(body, "\n") {
			if !strings.ContainsRune(line, '\x1b') {
				continue
			}

			Expect(line).To(HaveSuffix("\x1b[0m"))
		}
	})

	// the resources of a report fold on GitLab, and the log keeps the line the fold
	// opens with and the diff it hides, indented under it.
	It("strips a folded section down to its summary and the code block under it", func() {
		body, err := markdown.Render(strings.TrimSpace(`
<details>
<summary><code>+</code> <code>aws_s3_bucket.this</code></summary>

` + "```diff" + `
+ bucket: "logs"
- acl: "private"
` + "```" + `

</details>
`))
		Expect(err).NotTo(HaveOccurred())

		Expect(plain(body)).To(Equal("+ aws_s3_bucket.this\n  + bucket: \"logs\"\n  - acl: \"private\""))
		Expect(body).To(ContainSubstring("\x1b[32m+ bucket: \"logs\"\x1b[0m"))
		Expect(body).To(ContainSubstring("\x1b[31m- acl: \"private\"\x1b[0m"))
	})

	It("drops the styling when the environment asks for none", func() {
		GinkgoT().Setenv("NO_COLOR", "1")

		body, err := markdown.Render(document)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).NotTo(ContainSubstring("\x1b"))
	})

	It("logs the document one record per line", func() {
		out := &bytes.Buffer{}
		log := slog.New(slog.NewTextHandler(out, nil))

		body, err := markdown.Render(document)
		Expect(err).NotTo(HaveOccurred())

		Expect(markdown.Log(log, document)).To(Succeed())

		Expect(strings.Count(out.String(), "\n")).To(Equal(strings.Count(body, "\n") + 1))
	})
})
