package preview

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

var _ = Describe("Pulumi plan merge request report", func() {
	readFixture := func(name string) []byte {
		GinkgoHelper()

		data, err := os.ReadFile(filepath.Join("testdata", name))
		Expect(err).NotTo(HaveOccurred())

		return data
	}

	metadata := func() terraform.Metadata {
		return terraform.Metadata{
			Target:         "dev",
			Cwd:            ".",
			JobName:        "pulumi-preview",
			JobUrl:         "https://gitlab.example.test/project/-/jobs/1",
			PipelineId:     "42",
			PipelineUrl:    "https://gitlab.example.test/project/-/pipelines/42",
			CommitSha:      "0123456789abcdef",
			CommitShortSha: "01234567",
		}
	}

	resource := func(action terraform.Action, id string) terraform.Resource {
		GinkgoHelper()

		for _, resource := range action.Resources {
			if resource.Id == id {
				return resource
			}
		}

		Fail("resource not in action: " + id)

		return terraform.Resource{}
	}

	create := terraform.ChangeCreate

	It("summarizes an unwrapped Pulumi plan without rendering secrets", func() {
		report, err := parsePulumiPlanReport(readFixture("plan-unwrapped.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		Expect(report.Metadata.PlanSchema).To(BeEmpty())
		Expect(report.Metadata.ToolVersion).To(Equal("3.187.0"))
		Expect(report.Metadata.PlanTime).To(Equal("2026-05-31T12:00:00Z"))
		Expect(report.Total()).To(Equal(2))
		Expect(report.Actions).To(HaveLen(2))
		Expect(report.Actions[0].Action).To(Equal("create"))

		logs := resource(report.Actions[0], "urn:pulumi:dev::example::aws:s3/bucket:Bucket::logs")
		Expect(logs.Name).To(Equal("aws:s3/bucket:Bucket/logs"))
		// a secret is masked whether the plan holds its ciphertext or its plaintext.
		Expect(logs.Changes).To(Equal([]terraform.Change{
			{Name: "acl", Action: create, After: `"private"`},
			{Name: "region", Action: create, After: "[unknown]"},
			{Name: "rules", Action: create, Children: []terraform.Change{
				{Name: "[0]", Action: create, Children: []terraform.Change{
					{Name: "days", Action: create, After: "30"},
					{Name: "name", Action: create, After: `"expire"`},
				}},
			}},
			{Name: "tags", Action: create, Children: []terraform.Change{
				{Name: "env", Action: create, After: `"dev"`},
				{Name: "owner", Action: create, After: "[secret]"},
			}},
			{Name: "token", Action: create, After: "[secret]"},
		}))
		Expect(report.Actions[0].Outputs).To(ConsistOf(terraform.Output{
			Name:   "aws:s3/bucket:Bucket/logs",
			Fields: []string{"bucketName"},
		}))

		Expect(report.Actions[1].Action).To(Equal("update"))
		settings := resource(report.Actions[1], "urn:pulumi:dev::example::kubernetes:core/v1:ConfigMap::settings")
		Expect(settings.Changes).To(Equal([]terraform.Change{
			{Name: "data", Action: terraform.ChangeUpdate, Children: []terraform.Change{
				{Name: "level", After: `"debug"`},
				{Name: "password", After: "[secret]"},
			}},
			{Name: "immutable", Action: terraform.ChangeDelete},
		}))

		summary := terraform.Summarize(report)
		Expect(summary).To(Equal(terraform.Summary{
			Create: 1,
			Update: 1,
			Delete: 0,
		}))

		Expect(terraform.RenderSummary(summary)).To(Equal(`{
  "create": 1,
  "update": 1,
  "delete": 0
}
`))

		body, err := terraform.RenderReport(report)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(ContainSubstring("## Pulumi preview report"))
		Expect(body).To(ContainSubstring("| Stack | `dev` |"))
		Expect(body).To(ContainSubstring("| Pulumi version | `3.187.0` |"))
		Expect(body).To(ContainSubstring("### Output properties"))
		Expect(body).To(ContainSubstring("Total planned actions: 2."))
		Expect(body).To(ContainSubstring("+ acl: \"private\""))
		Expect(body).To(ContainSubstring("+ token: [secret]"))
		Expect(body).To(ContainSubstring("- immutable\n"))
		Expect(body).To(ContainSubstring("`bucketName`"))
		Expect(body).To(ContainSubstring("`metadata`"))
		Expect(body).NotTo(ContainSubstring("secret-config-value"))
		Expect(body).NotTo(ContainSubstring("secret-input-value"))
		Expect(body).NotTo(ContainSubstring("secret-owner-value"))
		Expect(body).NotTo(ContainSubstring("secret-update-value"))
		Expect(body).NotTo(ContainSubstring("secret-output-value"))
		Expect(body).NotTo(ContainSubstring("secret-new-value"))
		Expect(body).NotTo(ContainSubstring("ciphertext"))
		Expect(body).NotTo(ContainSubstring("aws:iam/role:Role"))
	})

	It("summarizes a versioned Pulumi deployment plan wrapper", func() {
		report, err := parsePulumiPlanReport(readFixture("plan-versioned.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		Expect(report.Metadata.PlanSchema).To(Equal("1"))
		Expect(report.Total()).To(Equal(3))
		Expect(report.Actions).To(HaveLen(3))
		Expect(report.Actions[0].Action).To(Equal("replace"))
		Expect(report.Actions[1].Action).To(Equal("create-replacement"))
		Expect(report.Actions[2].Action).To(Equal("delete-replaced"))

		for _, action := range report.Actions {
			worker := resource(action, "urn:pulumi:stage::example::aws:lambda/function:Function::worker")
			Expect(worker.Name).To(Equal("aws:lambda/function:Function/worker"))
			Expect(worker.Changes).To(Equal([]terraform.Change{
				{Name: "environment", Action: terraform.ChangeUpdate, After: "[secret]"},
				{Name: "runtime", Action: terraform.ChangeUpdate, After: `"nodejs22.x"`},
			}))
			Expect(action.Outputs).To(ConsistOf(terraform.Output{
				Name:   "aws:lambda/function:Function/worker",
				Fields: []string{"arn"},
			}))
		}

		summary := terraform.Summarize(report)
		Expect(summary).To(Equal(terraform.Summary{
			Create: 1,
			Update: 0,
			Delete: 1,
		}))

		body, err := terraform.RenderReport(report)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(ContainSubstring("Plan schema version"))
		Expect(body).To(ContainSubstring("~ runtime: \"nodejs22.x\""))
		Expect(body).NotTo(ContainSubstring("secret-arn-value"))
		Expect(body).NotTo(ContainSubstring("secret-env-value"))
	})

	Describe("attribute diff", func() {
		It("masks a secret however deep it sits and keeps a resource reference as its urn", func() {
			changes := resourceChanges(&apitype.GoalV1{
				InputDiff: apitype.PlanDiffV1{
					Adds: map[string]any{
						"members": []any{
							"a",
							map[string]any{
								"4dabf18193072939515e22adb298388d": "1b47061264138c4ac30d75fd1eb44270",
								"plaintext":                        `"secret-member"`,
							},
						},
						"role": map[string]any{
							"4dabf18193072939515e22adb298388d": "5cf8f73096256a8f31e491e813e4eb8e",
							"urn":                              "urn:pulumi:dev::example::aws:iam/role:Role::task",
							"id":                               "task",
						},
					},
				},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "members", Action: create, After: `["a", [secret]]`},
				{Name: "role", Action: create, After: `"urn:pulumi:dev::example::aws:iam/role:Role::task"`},
			}))
		})

		It("carries nothing for a resource without a goal", func() {
			Expect(resourceChanges(nil)).To(BeNil())
		})
	})
})
