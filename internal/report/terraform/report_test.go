package terraform_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

var _ = Describe("Terraform merge request report", func() {
	report := func(labels terraform.Labels) terraform.Report {
		return terraform.Report{
			Title:  "Example report",
			Labels: labels,
			Metadata: terraform.Metadata{
				Target:         "example",
				JobName:        "plan",
				JobUrl:         "https://gitlab.example.test/project/-/jobs/1",
				PipelineId:     "42",
				PipelineUrl:    "https://gitlab.example.test/project/-/pipelines/42",
				CommitShortSha: "01234567",
				ToolVersion:    "1.0.0",
			},
			Actions: []terraform.Action{
				{
					Action: "create",
					Resources: []terraform.Resource{{
						Name: "one",
						Changes: []terraform.Change{
							{Name: "name", Action: terraform.ChangeCreate, After: `"one"`},
						},
					}},
					Outputs: []terraform.Output{{Name: "first"}},
				},
				{
					Action:    "move",
					Resources: []terraform.Resource{{Name: "two", PreviousName: "old-two"}},
				},
			},
		}
	}

	terraformLabels := terraform.Labels{Target: "State", Outputs: "Outputs", ToolVersion: "Terraform version"}
	pulumiLabels := terraform.Labels{Target: "Stack", Outputs: "Output properties", ToolVersion: "Pulumi version"}

	headings := func(body string) []string {
		found := []string{}
		for line := range strings.SplitSeq(body, "\n") {
			if strings.HasPrefix(line, "#") {
				found = append(found, strings.SplitN(line, " ", 2)[0])
			}
		}

		return found
	}

	// both pipes render through this template, so the section skeleton is the
	// contract that keeps their reports readable side by side on one merge request.
	It("keeps one structure whichever labels it renders", func() {
		terraformBody, err := terraform.RenderReport(report(terraformLabels), terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		pulumiBody, err := terraform.RenderReport(report(pulumiLabels), terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		Expect(headings(terraformBody)).NotTo(BeEmpty())
		Expect(headings(terraformBody)).To(Equal(headings(pulumiBody)))
	})

	It("names the tool specific concepts from the labels", func() {
		body, err := terraform.RenderReport(report(pulumiLabels), terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("## Example report"))
		Expect(body).To(ContainSubstring("| Stack | `example` |"))
		Expect(body).To(ContainSubstring("| Pulumi version | `1.0.0` |"))
		Expect(body).To(ContainSubstring("| Job | [plan](https://gitlab.example.test/project/-/jobs/1) |"))
		Expect(body).To(ContainSubstring("### Output properties"))
		Expect(body).To(ContainSubstring("Total planned actions: 2."))
	})

	// the note is read on GitLab, which folds the html, and in the job log, which
	// strips it down to the summary line and the code block under it.
	It("folds every resource into a section of its own", func() {
		body, err := terraform.RenderReport(report(terraformLabels), terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring(strings.TrimSpace(`
<details>
<summary><code>+ create</code> <code>one</code></summary>

` + "```diff" + `
+ create
+ name: "one"
` + "```" + `

</details>
`)))
		Expect(body).To(ContainSubstring(strings.TrimSpace(`
<details>
<summary><code>&gt; move</code> <code>two</code> (moved from <code>old-two</code>)</summary>

` + "```diff" + `
> move
No attribute changes.
` + "```" + `

</details>
`)))
	})

	It("keeps a pipe in a metadata value inside its table cell", func() {
		body, err := terraform.RenderReport(terraform.Report{
			Title:    "Example report",
			Metadata: terraform.Metadata{JobName: "plan: [a | b]"},
		}, terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("| Job | `plan: [a \\| b]` |"))
	})

	It("says so when there is nothing to report", func() {
		body, err := terraform.RenderReport(terraform.Report{Title: "Example report"}, terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("No changes detected."))
		Expect(body).NotTo(ContainSubstring("### Metadata"))
	})

	Describe("attribute diff", func() {
		render := func(changes ...terraform.Change) string {
			GinkgoHelper()

			body, err := terraform.RenderReport(terraform.Report{
				Title: "Example report",
				Actions: []terraform.Action{{
					Action:    "update",
					Resources: []terraform.Resource{{Name: "one", Changes: changes}},
				}},
			}, terraform.DiffPlan)
			Expect(err).NotTo(HaveOccurred())

			return body
		}

		It("writes a change per line with its children indented under it", func() {
			body := render(
				terraform.Change{Name: "size", Action: terraform.ChangeUpdate, Before: "1", After: "2"},
				terraform.Change{Name: "tags", Action: terraform.ChangeCreate, Children: []terraform.Change{
					{Name: "env", Action: terraform.ChangeCreate, After: `"dev"`},
				}},
				terraform.Change{Name: "old", Action: terraform.ChangeDelete, Before: `"x"`},
				terraform.Change{Name: "gone", Action: terraform.ChangeDelete},
				terraform.Change{Name: "engine", Action: terraform.ChangeUpdate, Before: `"a"`, After: `"b"`, Note: "forces replacement"},
				terraform.Change{Name: "data", Action: terraform.ChangeUpdate, Children: []terraform.Change{
					{Name: "level", After: `"debug"`},
				}},
			)

			Expect(body).To(ContainSubstring(strings.TrimSpace(`
~ size: 1 -> 2
+ tags:
+   env: "dev"
- old: "x"
- gone
~ engine: "a" -> "b" # forces replacement
~ data:
    level: "debug"
`)))
		})

		// a code host only highlights the lines of a diff block that open with a plus or
		// a minus, which an update on one line marked with a tilde never does.
		It("writes an update as its old line removed and its new line added for a code host", func() {
			body, err := terraform.RenderReport(terraform.Report{
				Title: "Example report",
				Actions: []terraform.Action{{
					Action: "update",
					Resources: []terraform.Resource{{Name: "one", Changes: []terraform.Change{
						{Name: "size", Action: terraform.ChangeUpdate, Before: "1", After: "2"},
						{Name: "engine", Action: terraform.ChangeUpdate, Before: `"a"`, After: `"b"`, Note: "forces replacement"},
						{Name: "data", Action: terraform.ChangeUpdate, Children: []terraform.Change{
							{Name: "level", Action: terraform.ChangeUpdate, Before: `"info"`, After: `"debug"`},
							{Name: "mode", Action: terraform.ChangeUpdate, After: `"a"`},
						}},
						{
							Name:   "policy",
							Action: terraform.ChangeUpdate,
							Before: terraform.FormatValue("a\nb\n"),
							After:  terraform.FormatValue("a\nc\n"),
						},
						{Name: "tags", Action: terraform.ChangeCreate, After: `"x"`},
					}}},
				}},
			}, terraform.DiffUnified)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(ContainSubstring(strings.TrimSpace(`
- size: 1
+ size: 2
- engine: "a"
+ engine: "b" # forces replacement
  data:
-   level: "info"
+   level: "debug"
+   mode: "a"
  policy: <<-EOT
    a
-   b
+   c
  EOT
+ tags: "x"
`)))
			Expect(body).NotTo(MatchRegexp(`(?m)^~ .*:`))
		})

		// the report is what a plan is read through, so a diff it shortened would send
		// its reader back to the raw plan output for exactly the resource they came for.
		It("writes a long value out whole", func() {
			body := render(terraform.Change{Name: "blob", Action: terraform.ChangeCreate, After: strings.Repeat("x", 4096)})

			Expect(body).To(ContainSubstring("+ blob: " + strings.Repeat("x", 4096) + "\n"))
		})

		It("writes every attribute of every resource out whole", func() {
			resources := make([]terraform.Resource, 20)
			for index := range resources {
				resources[index].Name = fmt.Sprintf("resource-%02d", index)
				for attr := range 500 {
					resources[index].Changes = append(resources[index].Changes, terraform.Change{
						Name:   fmt.Sprintf("attr%03d", attr),
						Action: terraform.ChangeCreate,
						After:  "1",
					})
				}
			}

			body, err := terraform.RenderReport(terraform.Report{
				Title:   "Example report",
				Actions: []terraform.Action{{Action: "create", Resources: resources}},
			}, terraform.DiffPlan)
			Expect(err).NotTo(HaveOccurred())

			Expect(strings.Count(body, "+ attr499: 1")).To(Equal(len(resources)))
			Expect(body).NotTo(ContainSubstring("more lines"))
		})

		It("writes a value that spans lines as a heredoc carrying the action on every line", func() {
			body := render(terraform.Expand(terraform.ChangeCreate, "policy", "path \"a\" {\n  capabilities = [\"read\"]\n}\n"))

			Expect(body).To(ContainSubstring(strings.TrimSpace(`
+ policy: <<-EOT
+   path "a" {
+     capabilities = ["read"]
+   }
+ EOT
`)))
		})

		It("marks only the lines that changed between two values that span lines", func() {
			body := render(terraform.Change{
				Name:   "policy",
				Action: terraform.ChangeUpdate,
				Before: terraform.FormatValue("path \"a\" {\n  capabilities = [\"read\"]\n}\n"),
				After:  terraform.FormatValue("path \"a\" {\n  capabilities = [\"read\", \"update\"]\n}\n"),
				Note:   "forces replacement",
			})

			Expect(body).To(ContainSubstring(strings.TrimSpace(`
~ policy: <<-EOT # forces replacement
    path "a" {
-     capabilities = ["read"]
+     capabilities = ["read", "update"]
    }
  EOT
`)))
		})

		It("writes both sides whole when only one of them spans lines", func() {
			body := render(terraform.Change{
				Name:   "script",
				Action: terraform.ChangeUpdate,
				Before: terraform.FormatValue(nil),
				After:  terraform.FormatValue("a\nb"),
			})

			Expect(body).To(ContainSubstring(strings.TrimSpace(`
~ script:
-   null
+   <<-EOT
+     a
+     b
+   EOT
`)))
		})

		It("writes both sides whole when they are too long to line up", func() {
			body := render(terraform.Change{
				Name:   "blob",
				Action: terraform.ChangeUpdate,
				Before: terraform.FormatValue("shared\n" + strings.Repeat("a\n", 1100)),
				After:  terraform.FormatValue("shared\n" + strings.Repeat("b\n", 1100)),
			})

			Expect(body).To(ContainSubstring("\n-   shared\n"))
			Expect(body).To(ContainSubstring("\n+   shared\n"))
			Expect(strings.LastIndex(body, "-   a")).To(BeNumerically("<", strings.Index(body, "+   shared")))
		})

		It("keeps a value with backticks inside the code block", func() {
			body := render(terraform.Change{Name: "script", Action: terraform.ChangeCreate, After: "\"echo ```\""})

			Expect(body).To(ContainSubstring("````diff\n~ update\n+ script: \"echo ```\"\n````"))
		})
	})

	Describe("Compare", func() {
		update := terraform.ChangeUpdate

		It("says nothing changed between equal values", func() {
			_, changed := terraform.Compare("a", map[string]any{"b": "c"}, map[string]any{"b": "c"})

			Expect(changed).To(BeFalse())
		})

		It("compares objects and lists of objects child by child", func() {
			change, changed := terraform.Compare("root",
				map[string]any{"keep": "x", "gone": "y", "rules": []any{map[string]any{"port": float64(80)}}},
				map[string]any{"keep": "x", "new": "z", "rules": []any{map[string]any{"port": float64(443)}}},
			)

			Expect(changed).To(BeTrue())
			Expect(change).To(Equal(terraform.Change{Name: "root", Action: update, Children: []terraform.Change{
				{Name: "gone", Action: terraform.ChangeDelete, Before: `"y"`},
				{Name: "new", Action: terraform.ChangeCreate, After: `"z"`},
				{Name: "rules", Action: update, Children: []terraform.Change{
					{Name: "[0]", Action: update, Children: []terraform.Change{
						{Name: "port", Action: update, Before: "80", After: "443"},
					}},
				}},
			}}))
		})

		It("compares two strings holding JSON documents as the documents", func() {
			change, _ := terraform.Compare("values", `{"a": 1, "b": 2}`, `{"a": 1, "b": 3}`)

			Expect(change).To(Equal(terraform.Change{Name: "values", Action: update, Note: terraform.JSONEncoded, Children: []terraform.Change{
				{Name: "b", Action: update, Before: "2", After: "3"},
			}}))
		})

		It("writes out what a marker replaces attribute by attribute", func() {
			change, _ := terraform.Compare("meta", map[string]any{"a": "x"}, terraform.Marker("[unknown]"))

			Expect(change).To(Equal(terraform.Change{Name: "meta", Action: update, Children: []terraform.Change{
				{Name: "a", Action: update, Before: `"x"`, After: "[unknown]"},
			}}))
		})
	})

	Describe("Expand", func() {
		It("writes objects and lists of objects out one child per line", func() {
			change := terraform.Expand(terraform.ChangeCreate, "root", map[string]any{
				"cidr":  []any{"10.0.0.0/8", "10.1.0.0/16"},
				"empty": map[string]any{},
				"rules": []any{map[string]any{"port": float64(80)}},
				"key":   terraform.Marker("(sensitive value)"),
				"none":  nil,
			})

			Expect(change).To(Equal(terraform.Change{
				Name:   "root",
				Action: terraform.ChangeCreate,
				Children: []terraform.Change{
					{Name: "cidr", Action: terraform.ChangeCreate, After: `["10.0.0.0/8", "10.1.0.0/16"]`},
					{Name: "empty", Action: terraform.ChangeCreate, After: "{}"},
					{Name: "key", Action: terraform.ChangeCreate, After: "(sensitive value)"},
					{Name: "none", Action: terraform.ChangeCreate, After: "null"},
					{Name: "rules", Action: terraform.ChangeCreate, Children: []terraform.Change{
						{Name: "[0]", Action: terraform.ChangeCreate, Children: []terraform.Change{
							{Name: "port", Action: terraform.ChangeCreate, After: "80"},
						}},
					}},
				},
			}))
		})

		It("writes a string holding a JSON document out as the document", func() {
			Expect(terraform.Expand(terraform.ChangeCreate, "values", `{"cm": {"url": "a"}}`)).To(Equal(terraform.Change{
				Name:   "values",
				Action: terraform.ChangeCreate,
				Note:   terraform.JSONEncoded,
				Children: []terraform.Change{
					{Name: "cm", Action: terraform.ChangeCreate, Children: []terraform.Change{
						{Name: "url", Action: terraform.ChangeCreate, After: `"a"`},
					}},
				},
			}))
		})

		It("keeps a string holding a JSON scalar list or an invalid document as the string", func() {
			Expect(terraform.Expand(terraform.ChangeCreate, "a", `["x"]`).After).To(Equal(`"[\"x\"]"`))
			Expect(terraform.Expand(terraform.ChangeCreate, "b", `{broken`).After).To(Equal(`"{broken"`))
		})

		It("writes a list holding a value that spans lines one element per line", func() {
			Expect(terraform.Expand(terraform.ChangeCreate, "lines", []any{"a", "b\nc"})).To(Equal(terraform.Change{
				Name:   "lines",
				Action: terraform.ChangeCreate,
				Children: []terraform.Change{
					{Name: "[0]", Action: terraform.ChangeCreate, After: `"a"`},
					{Name: "[1]", Action: terraform.ChangeCreate, After: "<<-EOT\n  b\n  c\nEOT"},
				},
			}))
		})

		It("keeps the value of a deletion on the before side", func() {
			Expect(terraform.Expand(terraform.ChangeDelete, "name", "x")).To(Equal(terraform.Change{
				Name:   "name",
				Action: terraform.ChangeDelete,
				Before: `"x"`,
			}))
		})
	})

	It("keeps a string that spans lines escaped inside a nested value", func() {
		Expect(terraform.FormatValue([]any{"a\nb"})).To(Equal(`["a\nb"]`))
	})

	It("writes a nested value on one line with its keys sorted and its markers bare", func() {
		Expect(terraform.FormatValue(map[string]any{
			"b": []any{true, float64(1.5), "x"},
			"a": terraform.Marker("[secret]"),
		})).To(Equal(`{"a": [secret], "b": [true, 1.5, "x"]}`))
	})

	DescribeTable("leads a resource with the sign of its action",
		func(action string, glyph string) {
			Expect(terraform.Glyph(action)).To(Equal(glyph))
		},
		Entry(nil, "create", "+"),
		Entry(nil, "create-replacement", "+"),
		Entry(nil, "create+forget", "+"),
		Entry(nil, "update", "~"),
		Entry(nil, "update-replacement", "~"),
		Entry(nil, "delete", "-"),
		Entry(nil, "delete-replaced", "-"),
		Entry(nil, "discard", "-"),
		Entry(nil, "discard-replaced", "-"),
		Entry(nil, "remove-pending-replace", "-"),
		Entry(nil, "replace", "-/+"),
		Entry(nil, "import", "="),
		Entry(nil, "import-replacement", "="),
		Entry(nil, "read", "<="),
		Entry(nil, "read-replacement", "<="),
		Entry(nil, "refresh", "<="),
		Entry(nil, "move", ">"),
		Entry(nil, "forget", "."),
		Entry(nil, "something-else", "?"),
	)

	It("labels every resource with its operation in the summary, its section and the diff it opens", func() {
		body, err := terraform.RenderReport(terraform.Report{
			Title: "Example report",
			Actions: []terraform.Action{{
				Action:    "replace",
				Resources: []terraform.Resource{{Name: "one", Detail: "create before destroy"}},
			}},
		}, terraform.DiffUnified)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("| `-/+ replace` | 1 | 0 |"))
		Expect(body).To(ContainSubstring("#### `-/+ replace` (1)"))
		Expect(body).To(ContainSubstring(strings.TrimSpace(`
<summary><code>-/+ replace</code> <code>one</code> (create before destroy)</summary>

` + "```diff" + `
-/+ replace
No attribute changes.
` + "```" + `
`)))
	})
})
