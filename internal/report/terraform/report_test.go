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
		terraformBody, err := terraform.RenderReport(report(terraformLabels))
		Expect(err).NotTo(HaveOccurred())

		pulumiBody, err := terraform.RenderReport(report(pulumiLabels))
		Expect(err).NotTo(HaveOccurred())

		Expect(headings(terraformBody)).NotTo(BeEmpty())
		Expect(headings(terraformBody)).To(Equal(headings(pulumiBody)))
	})

	It("names the tool specific concepts from the labels", func() {
		body, err := terraform.RenderReport(report(pulumiLabels))
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
		body, err := terraform.RenderReport(report(terraformLabels))
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring(strings.TrimSpace(`
<details>
<summary><code>+</code> <code>one</code></summary>

` + "```diff" + `
+ name: "one"
` + "```" + `

</details>
`)))
		Expect(body).To(ContainSubstring(strings.TrimSpace(`
<details>
<summary><code>&gt;</code> <code>two</code> (moved from <code>old-two</code>)</summary>

` + "```diff" + `
No attribute changes.
` + "```" + `

</details>
`)))
	})

	It("says so when there is nothing to report", func() {
		body, err := terraform.RenderReport(terraform.Report{Title: "Example report"})
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
			})
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
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(strings.Count(body, "+ attr499: 1")).To(Equal(len(resources)))
			Expect(body).NotTo(ContainSubstring("more lines"))
		})

		It("keeps a value with backticks inside the code block", func() {
			body := render(terraform.Change{Name: "script", Action: terraform.ChangeCreate, After: "\"echo ```\""})

			Expect(body).To(ContainSubstring("````diff\n+ script: \"echo ```\"\n````"))
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

		It("keeps the value of a deletion on the before side", func() {
			Expect(terraform.Expand(terraform.ChangeDelete, "name", "x")).To(Equal(terraform.Change{
				Name:   "name",
				Action: terraform.ChangeDelete,
				Before: `"x"`,
			}))
		})
	})

	It("writes a nested value on one line with its keys sorted and its markers bare", func() {
		Expect(terraform.FormatValue(map[string]any{
			"b": []any{true, float64(1.5), "x"},
			"a": terraform.Marker("[secret]"),
		})).To(Equal(`{"a": [secret], "b": [true, 1.5, "x"]}`))
	})

	It("leads a resource with the sign of its action", func() {
		Expect(terraform.Glyph("create")).To(Equal("+"))
		Expect(terraform.Glyph("delete-replaced")).To(Equal("-"))
		Expect(terraform.Glyph("replace")).To(Equal("-/+"))
		Expect(terraform.Glyph("something-else")).To(Equal("?"))
	})
})
