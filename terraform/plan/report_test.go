package plan

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terraform merge request report", func() {
	readFixture := func(name string) []byte {
		GinkgoHelper()

		data, err := os.ReadFile(filepath.Join("testdata", name))
		Expect(err).NotTo(HaveOccurred())

		return data
	}

	It("summarizes resource and output actions without rendering values", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"))
		Expect(err).NotTo(HaveOccurred())

		Expect(report.Summary).To(ConsistOf(
			mergeRequestReportSummary{Action: "create", ResourceCount: 1, OutputCount: 1},
			mergeRequestReportSummary{Action: "update", ResourceCount: 1, OutputCount: 1},
			mergeRequestReportSummary{Action: "delete", ResourceCount: 0, OutputCount: 1},
			mergeRequestReportSummary{Action: "replace", ResourceCount: 1, OutputCount: 0},
			mergeRequestReportSummary{Action: "read", ResourceCount: 1, OutputCount: 0},
		))
		Expect(report.ResourceGroups).To(ContainElement(mergeRequestReportGroup{
			Action: "replace",
			Items:  []string{"aws_db_instance.main"},
		}))
		Expect(report.OutputGroups).To(ContainElement(mergeRequestReportGroup{
			Action: "delete",
			Items:  []string{"password"},
		}))

		summary, err := summarizeTerraformShowPlan(readFixture("plan.json"))
		Expect(err).NotTo(HaveOccurred())
		Expect(summary).To(Equal(terraformSummary{
			Create: 2,
			Update: 1,
			Delete: 1,
		}))

		summaryBody, err := renderSummary(summary)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(summaryBody)).To(Equal(`{
  "create": 2,
  "update": 1,
  "delete": 1
}
`))

		body, err := renderMergeRequestReport(report)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("Terraform plan report"))
		Expect(body).To(ContainSubstring("`aws_s3_bucket.logs`"))
		Expect(body).To(ContainSubstring("`bucket_name`"))
		Expect(body).NotTo(ContainSubstring("secret-bucket-name"))
		Expect(body).NotTo(ContainSubstring("secret-old-user-data"))
		Expect(body).NotTo(ContainSubstring("secret-new-user-data"))
		Expect(body).NotTo(ContainSubstring("secret-old-password"))
		Expect(body).NotTo(ContainSubstring("secret-new-password"))
		Expect(body).NotTo(ContainSubstring("secret-output-bucket"))
		Expect(body).NotTo(ContainSubstring("secret-old-endpoint"))
		Expect(body).NotTo(ContainSubstring("secret-new-endpoint"))
	})
})
