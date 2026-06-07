package plan

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terraform report", func() {
	readFixture := func(name string) []byte {
		GinkgoHelper()

		data, err := os.ReadFile(filepath.Join("testdata", name))
		Expect(err).NotTo(HaveOccurred())

		return data
	}

	It("summarizes resource and output actions without rendering values", func() {
		plan, err := parseTerraformShowPlan(readFixture("plan.json"))
		Expect(err).NotTo(HaveOccurred())

		report := buildTerraformReport(plan)
		Expect(report.Summary).To(ConsistOf(
			reportSummary{Action: "create", ResourceCount: 1, OutputCount: 1},
			reportSummary{Action: "update", ResourceCount: 1, OutputCount: 1},
			reportSummary{Action: "delete", ResourceCount: 0, OutputCount: 1},
			reportSummary{Action: "replace", ResourceCount: 1, OutputCount: 0},
			reportSummary{Action: "read", ResourceCount: 1, OutputCount: 0},
		))
		Expect(report.ResourceGroups).To(ContainElement(reportGroup{
			Action: "replace",
			Items:  []string{"aws_db_instance.main"},
		}))
		Expect(report.OutputGroups).To(ContainElement(reportGroup{
			Action: "delete",
			Items:  []string{"password"},
		}))

		body, err := renderReport(report)
		Expect(err).NotTo(HaveOccurred())

		content := string(body)
		Expect(content).To(ContainSubstring("Terraform report"))
		Expect(content).To(ContainSubstring("aws_s3_bucket.logs"))
		Expect(content).To(ContainSubstring("bucket_name"))
		Expect(content).NotTo(ContainSubstring("secret-bucket-name"))
		Expect(content).NotTo(ContainSubstring("secret-old-user-data"))
		Expect(content).NotTo(ContainSubstring("secret-new-user-data"))
		Expect(content).NotTo(ContainSubstring("secret-old-password"))
		Expect(content).NotTo(ContainSubstring("secret-new-password"))
		Expect(content).NotTo(ContainSubstring("secret-output-bucket"))
		Expect(content).NotTo(ContainSubstring("secret-old-endpoint"))
		Expect(content).NotTo(ContainSubstring("secret-new-endpoint"))
	})

	It("renders Terraform summary counts from flattened resource actions", func() {
		plan, err := parseTerraformShowPlan(readFixture("plan.json"))
		Expect(err).NotTo(HaveOccurred())

		report := summarizeTerraformPlan(plan)
		Expect(report).To(Equal(terraformSummary{
			Create: 2,
			Update: 1,
			Delete: 1,
		}))

		body, err := renderSummary(report)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(Equal("{\n  \"create\": 2,\n  \"update\": 1,\n  \"delete\": 1\n}\n"))
	})
})
