package terraform_test

import (
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
					Action:    "create",
					Resources: []terraform.Resource{{Name: "one"}},
					Outputs:   []terraform.Output{{Name: "first"}},
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
		Expect(body).To(ContainSubstring("- `two` (moved from `old-two`)"))
	})

	It("says so when there is nothing to report", func() {
		body, err := terraform.RenderReport(terraform.Report{Title: "Example report"})
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("No changes detected."))
		Expect(body).NotTo(ContainSubstring("### Metadata"))
	})
})
