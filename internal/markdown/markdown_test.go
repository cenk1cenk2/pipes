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

// what a reader sees as the color of a line.
var escapes = regexp.MustCompile(`\x1b\[[0-9;]*m`)

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

var _ = Describe("Plan diff lexer", func() {
	const dimmed = "(sensitive value)"

	BeforeEach(func() {
		markdown.RegisterDiffLexer(dimmed, "No attribute changes.")
	})

	render := func(lines ...string) string {
		body, err := markdown.Render("```diff\n" + strings.Join(lines, "\n") + "\n```\n")
		Expect(err).NotTo(HaveOccurred())

		return body
	}

	// the escape sequence a line carries is what a reader sees as its color; the specs
	// below only care that two of them differ, never which color chroma settled on.
	color := func(body string, contains string) string {
		for line := range strings.SplitSeq(body, "\n") {
			if !strings.Contains(line, contains) {
				continue
			}

			if found := escapes.FindString(line); found != "" {
				return found
			}

			return ""
		}

		Fail("no line carrying " + contains)

		return ""
	}

	// chroma's own diff lexer reads every one of these as context and leaves them
	// unstyled, which is what this lexer replaces it for.
	It("colors a line by the glyph it opens with", func() {
		body := render(`+ created: "a"`, `- deleted: "b"`, `~ changed: "c" -> "d"`)

		created := color(body, "created")
		deleted := color(body, "deleted")
		changed := color(body, "changed")

		Expect(created).NotTo(BeEmpty())
		Expect(deleted).NotTo(BeEmpty())
		Expect(changed).NotTo(BeEmpty())
		Expect([]string{created, deleted, changed}).To(HaveLen(3))
		Expect(created).NotTo(Equal(deleted))
		Expect(created).NotTo(Equal(changed))
		Expect(deleted).NotTo(Equal(changed))
	})

	It("dims what the pipe stands in for a value it will not show", func() {
		body := render(`~ password: ` + dimmed + ` -> ` + dimmed)

		line := ""
		for candidate := range strings.SplitSeq(body, "\n") {
			if strings.Contains(candidate, "password") {
				line = candidate
			}
		}

		Expect(escapes.FindAllString(line, -1)).To(ContainElement(Not(Equal(color(body, "password")))))
		Expect(line).To(ContainSubstring(dimmed))
	})

	It("dims a resource that carries nothing to show", func() {
		body := render("No attribute changes.")

		Expect(color(body, "No attribute changes.")).NotTo(BeEmpty())
	})

	It("marks what forces a resource to be replaced", func() {
		body := render(`~ engine: "15" -> "16" # forces replacement`)

		line := ""
		for candidate := range strings.SplitSeq(body, "\n") {
			if strings.Contains(candidate, "engine") {
				line = candidate
			}
		}

		Expect(len(escapes.FindAllString(line, -1))).To(BeNumerically(">", 2))
		Expect(line).To(ContainSubstring("# forces replacement"))
	})

	It("leaves a nested line to read as the one above it", func() {
		body := render(`~ data:`, `    level: "debug"`)

		Expect(color(body, "level")).To(BeEmpty())
	})

	It("drops the styling when the environment asks for none", func() {
		GinkgoT().Setenv("NO_COLOR", "1")

		Expect(render(`~ changed: "c" -> "d"`)).NotTo(ContainSubstring("\x1b"))
	})
})
