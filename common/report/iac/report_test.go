package iac_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/common/report/iac"
)

var _ = Describe("IaC merge request report", func() {
	report := func(labels iac.Labels) iac.Report {
		return iac.Report{
			Title:  "Example report",
			Labels: labels,
			Metadata: iac.Metadata{
				Target:         "example",
				JobName:        "plan",
				JobUrl:         "https://gitlab.example.test/project/-/jobs/1",
				PipelineId:     "42",
				PipelineUrl:    "https://gitlab.example.test/project/-/pipelines/42",
				CommitShortSha: "01234567",
				ToolVersion:    "1.0.0",
			},
			Actions: []iac.Action{
				{
					Action:    "create",
					Resources: []iac.Resource{{Name: "one"}},
					Outputs:   []iac.Output{{Name: "first"}},
				},
				{
					Action:    "move",
					Resources: []iac.Resource{{Name: "two", PreviousName: "old-two"}},
				},
			},
		}
	}

	terraformLabels := iac.Labels{Target: "State", Outputs: "Outputs", ToolVersion: "Terraform version"}
	pulumiLabels := iac.Labels{Target: "Stack", Outputs: "Output properties", ToolVersion: "Pulumi version"}

	headings := func(body string) []string {
		found := []string{}
		for _, line := range strings.Split(body, "\n") {
			if strings.HasPrefix(line, "#") {
				found = append(found, strings.SplitN(line, " ", 2)[0])
			}
		}

		return found
	}

	// Both pipes render through this template, so the section skeleton is the
	// contract that keeps their reports readable side by side on one merge request.
	It("keeps one structure whichever labels it renders", func() {
		terraform, err := iac.RenderMergeRequestReport(report(terraformLabels))
		Expect(err).NotTo(HaveOccurred())

		pulumi, err := iac.RenderMergeRequestReport(report(pulumiLabels))
		Expect(err).NotTo(HaveOccurred())

		Expect(headings(terraform)).To(Equal(headings(pulumi)))
		Expect(headings(terraform)).To(Equal([]string{"##", "###", "###", "###", "####", "####", "###", "####"}))
	})

	It("names the tool specific concepts from the labels", func() {
		body, err := iac.RenderMergeRequestReport(report(pulumiLabels))
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
		body, err := iac.RenderMergeRequestReport(iac.Report{Title: "Example report"})
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("No changes detected."))
		Expect(body).NotTo(ContainSubstring("### Metadata"))
	})
})
