package plan

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gitlab.kilic.dev/devops/pipes/common/report/iac"
)

var _ = Describe("Terraform merge request report", func() {
	readFixture := func(name string) []byte {
		GinkgoHelper()

		data, err := os.ReadFile(filepath.Join("testdata", name))
		Expect(err).NotTo(HaveOccurred())

		return data
	}

	metadata := func() iac.Metadata {
		return iac.Metadata{
			Target:         "production",
			Cwd:            ".",
			JobName:        "tf-plan",
			JobUrl:         "https://gitlab.example.test/project/-/jobs/1",
			PipelineId:     "42",
			PipelineUrl:    "https://gitlab.example.test/project/-/pipelines/42",
			CommitSha:      "0123456789abcdef",
			CommitShortSha: "01234567",
		}
	}

	It("summarizes resource and output actions without rendering values", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		Expect(report.Metadata.ToolVersion).To(Equal("1.9.8"))
		Expect(report.Metadata.PlanTime).To(Equal("2026-05-31T12:00:00Z"))
		Expect(report.Metadata.PlanSchema).To(Equal("1.2"))
		Expect(report.Total()).To(Equal(5))

		Expect(report.Actions).To(ConsistOf(
			iac.Action{
				Action:    "create",
				Resources: []iac.Resource{{Name: "aws_s3_bucket.logs"}},
				Outputs:   []iac.Output{{Name: "bucket_name"}},
			},
			iac.Action{
				Action:    "update",
				Resources: []iac.Resource{{Name: "aws_instance.web"}},
				Outputs:   []iac.Output{{Name: "endpoint"}},
			},
			iac.Action{
				Action:    "delete",
				Resources: nil,
				Outputs:   []iac.Output{{Name: "password"}},
			},
			iac.Action{
				Action:    "replace",
				Resources: []iac.Resource{{Name: "aws_db_instance.main"}},
				Outputs:   nil,
			},
			iac.Action{
				Action:    "move",
				Resources: []iac.Resource{{Name: "aws_sqs_queue.jobs", PreviousName: "aws_sqs_queue.legacy_jobs"}},
				Outputs:   nil,
			},
			iac.Action{
				Action:    "read",
				Resources: []iac.Resource{{Name: "data.aws_caller_identity.current"}},
				Outputs:   nil,
			},
		))
	})

	It("orders actions by the terraform step vocabulary", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		names := []string{}
		for _, action := range report.Actions {
			names = append(names, action.Action)
		}

		Expect(names).To(Equal([]string{"create", "update", "delete", "replace", "move", "read"}))
	})

	It("renders the report without leaking values", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		body, err := iac.RenderMergeRequestReport(report)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("## Terraform plan report"))
		Expect(body).To(ContainSubstring("| State | `production` |"))
		Expect(body).To(ContainSubstring("| Terraform version | `1.9.8` |"))
		Expect(body).To(ContainSubstring("[tf-plan](https://gitlab.example.test/project/-/jobs/1)"))
		Expect(body).To(ContainSubstring("Total planned actions: 5."))
		Expect(body).To(ContainSubstring("`aws_s3_bucket.logs`"))
		Expect(body).To(ContainSubstring("`aws_sqs_queue.jobs` (moved from `aws_sqs_queue.legacy_jobs`)"))
		Expect(body).To(ContainSubstring("`bucket_name`"))
		Expect(body).NotTo(ContainSubstring("null_resource.noop"))

		Expect(body).NotTo(ContainSubstring("secret-bucket-name"))
		Expect(body).NotTo(ContainSubstring("secret-old-user-data"))
		Expect(body).NotTo(ContainSubstring("secret-new-user-data"))
		Expect(body).NotTo(ContainSubstring("secret-old-password"))
		Expect(body).NotTo(ContainSubstring("secret-new-password"))
		Expect(body).NotTo(ContainSubstring("secret-output-bucket"))
		Expect(body).NotTo(ContainSubstring("secret-old-endpoint"))
		Expect(body).NotTo(ContainSubstring("secret-new-endpoint"))
		Expect(body).NotTo(ContainSubstring("secret-queue-name"))
	})

	It("keeps the summary artifact independent of the report grouping", func() {
		summary, err := summarizeTerraformShowPlan(readFixture("plan.json"))
		Expect(err).NotTo(HaveOccurred())
		Expect(summary).To(Equal(terraformSummary{
			Create: 2,
			Update: 1,
			Delete: 1,
		}))

		body, err := renderSummary(summary)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(Equal(`{
  "create": 2,
  "update": 1,
  "delete": 1
}
`))
	})
})
