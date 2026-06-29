package preview

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pulumi plan merge request report", func() {
	readFixture := func(name string) []byte {
		GinkgoHelper()

		data, err := os.ReadFile(filepath.Join("testdata", name))
		Expect(err).NotTo(HaveOccurred())

		return data
	}

	metadata := func() MergeRequestReportMetadata {
		return MergeRequestReportMetadata{
			Stack:          "dev",
			JobName:        "pulumi-preview",
			JobUrl:         "https://gitlab.example.test/project/-/jobs/1",
			PipelineId:     "42",
			PipelineUrl:    "https://gitlab.example.test/project/-/pipelines/42",
			CommitSha:      "0123456789abcdef",
			CommitShortSha: "01234567",
		}
	}

	It("summarizes an unwrapped Pulumi plan without rendering values", func() {
		report, err := parsePulumiPlanReport(readFixture("plan-unwrapped.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		Expect(report.PlanVersion).To(Equal(0))
		Expect(report.Manifest.Version).To(Equal("3.187.0"))
		Expect(report.Manifest.Time).To(Equal("2026-05-31T12:00:00Z"))
		Expect(report.TotalActions).To(Equal(2))
		Expect(report.Actions).To(HaveLen(2))
		Expect(report.Actions[0].Action).To(Equal("create"))
		Expect(report.Actions[0].Resources).To(ContainElement(PulumiPlanResource{
			Urn:  "urn:pulumi:dev::example::aws:s3/bucket:Bucket::logs",
			Type: "aws:s3/bucket:Bucket",
			Name: "logs",
		}))
		Expect(report.Actions[0].Outputs).To(ConsistOf(PulumiPlanOutput{
			Resource: PulumiPlanResource{
				Urn:  "urn:pulumi:dev::example::aws:s3/bucket:Bucket::logs",
				Type: "aws:s3/bucket:Bucket",
				Name: "logs",
			},
			Names: []string{"bucketName"},
		}))

		summary := summarizePulumiReport(report)
		Expect(summary).To(Equal(pulumiSummary{
			Create: 1,
			Update: 1,
			Delete: 0,
		}))

		summaryBody, err := renderSummary(summary)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(summaryBody)).To(Equal(`{
  "create": 1,
  "update": 1,
  "delete": 0
}
`))

		body, err := renderMergeRequestReport(report)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(ContainSubstring("Pulumi preview report"))
		Expect(body).To(ContainSubstring("`bucketName`"))
		Expect(body).To(ContainSubstring("`metadata`"))
		Expect(body).NotTo(ContainSubstring("secret-config-value"))
		Expect(body).NotTo(ContainSubstring("secret-input-value"))
		Expect(body).NotTo(ContainSubstring("secret-output-value"))
		Expect(body).NotTo(ContainSubstring("secret-old-value"))
		Expect(body).NotTo(ContainSubstring("secret-new-value"))
		Expect(body).NotTo(ContainSubstring("aws:iam/role:Role"))
	})

	It("summarizes a versioned Pulumi deployment plan wrapper", func() {
		report, err := parsePulumiPlanReport(readFixture("plan-versioned.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		Expect(report.PlanVersion).To(Equal(1))
		Expect(report.TotalActions).To(Equal(3))
		Expect(report.Actions).To(HaveLen(3))
		Expect(report.Actions[0].Action).To(Equal("replace"))
		Expect(report.Actions[1].Action).To(Equal("create-replacement"))
		Expect(report.Actions[2].Action).To(Equal("delete-replaced"))

		for _, action := range report.Actions {
			Expect(action.Resources).To(ContainElement(PulumiPlanResource{
				Urn:  "urn:pulumi:stage::example::aws:lambda/function:Function::worker",
				Type: "aws:lambda/function:Function",
				Name: "worker",
			}))
			Expect(action.Outputs).To(ConsistOf(PulumiPlanOutput{
				Resource: PulumiPlanResource{
					Urn:  "urn:pulumi:stage::example::aws:lambda/function:Function::worker",
					Type: "aws:lambda/function:Function",
					Name: "worker",
				},
				Names: []string{"arn"},
			}))
		}

		summary := summarizePulumiReport(report)
		Expect(summary).To(Equal(pulumiSummary{
			Create: 1,
			Update: 0,
			Delete: 1,
		}))

		body, err := renderMergeRequestReport(report)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(ContainSubstring("Plan schema version"))
		Expect(body).NotTo(ContainSubstring("secret-arn-value"))
	})
})
