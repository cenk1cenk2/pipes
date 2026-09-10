package plan

import (
	"os"
	"path/filepath"

	tfjson "github.com/hashicorp/terraform-json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

var _ = Describe("Terraform merge request report", func() {
	readFixture := func(name string) []byte {
		GinkgoHelper()

		data, err := os.ReadFile(filepath.Join("testdata", name))
		Expect(err).NotTo(HaveOccurred())

		return data
	}

	metadata := func() terraform.Metadata {
		return terraform.Metadata{
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

	resource := func(report terraform.Report, name string) terraform.Resource {
		GinkgoHelper()

		for _, action := range report.Actions {
			for _, resource := range action.Resources {
				if resource.Name == name {
					return resource
				}
			}
		}

		Fail("resource not in report: " + name)

		return terraform.Resource{}
	}

	It("summarizes resource and output actions", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		Expect(report.Metadata.ToolVersion).To(Equal("1.9.8"))
		Expect(report.Metadata.PlanTime).To(Equal("2026-05-31T12:00:00Z"))
		Expect(report.Metadata.PlanSchema).To(Equal("1.2"))
		Expect(report.Total()).To(Equal(5))

		resources := map[string][]string{}
		outputs := map[string][]string{}
		for _, action := range report.Actions {
			for _, resource := range action.Resources {
				resources[action.Action] = append(resources[action.Action], resource.Name)
			}
			for _, output := range action.Outputs {
				outputs[action.Action] = append(outputs[action.Action], output.Name)
			}
		}

		Expect(resources).To(Equal(map[string][]string{
			"create":  {"aws_s3_bucket.logs"},
			"update":  {"aws_instance.web"},
			"replace": {"aws_db_instance.main"},
			"move":    {"aws_sqs_queue.jobs"},
			"read":    {"data.aws_caller_identity.current"},
		}))
		Expect(outputs).To(Equal(map[string][]string{
			"create": {"bucket_name"},
			"update": {"endpoint"},
			"delete": {"password"},
		}))
		Expect(resource(report, "aws_sqs_queue.jobs").PreviousName).To(Equal("aws_sqs_queue.legacy_jobs"))
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

	It("carries the attribute changes with the masks applied", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		create := terraform.ChangeCreate
		update := terraform.ChangeUpdate

		Expect(resource(report, "aws_s3_bucket.logs").Changes).To(Equal([]terraform.Change{
			{Name: "arn", Action: create, After: "(known after apply)"},
			{Name: "bucket", Action: create, After: "(sensitive value)"},
			{Name: "cors_rule", Action: create, Children: []terraform.Change{
				{Name: "[0]", Action: create, Children: []terraform.Change{
					{Name: "allowed_methods", Action: create, After: `["GET"]`},
					{Name: "id", Action: create, After: "(known after apply)"},
					{Name: "max_age_seconds", Action: create, After: "300"},
				}},
			}},
			{Name: "force_destroy", Action: create, After: "false"},
			{Name: "id", Action: create, After: "(known after apply)"},
			{Name: "lifecycle_rule", Action: create, After: "[]"},
			{Name: "tags", Action: create, Children: []terraform.Change{
				{Name: "env", Action: create, After: `"dev"`},
				{Name: "team", Action: create, After: `"platform"`},
			}},
		}))

		// a sensitive value that changes is still a change, only without its values.
		Expect(resource(report, "aws_instance.web").Changes).To(Equal([]terraform.Change{
			{Name: "instance_type", Action: update, Before: `"t3.micro"`, After: `"t3.small"`},
			{Name: "user_data", Action: update, Before: "(sensitive value)", After: "(sensitive value)"},
		}))

		Expect(resource(report, "aws_db_instance.main").Changes).To(Equal([]terraform.Change{
			{Name: "engine_version", Action: update, Before: `"15"`, After: `"16"`, Note: "forces replacement"},
			{Name: "password", Action: update, Before: "(sensitive value)", After: "(sensitive value)"},
		}))

		Expect(resource(report, "data.aws_caller_identity.current").Changes).To(Equal([]terraform.Change{
			{Name: "account_id", Action: create, After: "(sensitive value)"},
		}))

		Expect(resource(report, "aws_sqs_queue.jobs").Changes).To(BeNil())
	})

	It("renders the report without leaking values", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		body, err := terraform.RenderReport(report)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("## Terraform plan report"))
		Expect(body).To(ContainSubstring("| State | `production` |"))
		Expect(body).To(ContainSubstring("| Terraform version | `1.9.8` |"))
		Expect(body).To(ContainSubstring("[tf-plan](https://gitlab.example.test/project/-/jobs/1)"))
		Expect(body).To(ContainSubstring("Total planned actions: 5."))
		Expect(body).To(ContainSubstring("<summary><code>+</code> <code>aws_s3_bucket.logs</code></summary>"))
		Expect(body).To(ContainSubstring(
			"<summary><code>&gt;</code> <code>aws_sqs_queue.jobs</code> (moved from <code>aws_sqs_queue.legacy_jobs</code>)</summary>",
		))
		Expect(body).To(ContainSubstring("~ instance_type: \"t3.micro\" -> \"t3.small\""))
		Expect(body).To(ContainSubstring("~ engine_version: \"15\" -> \"16\" # forces replacement"))
		Expect(body).To(ContainSubstring("+ arn: (known after apply)"))
		Expect(body).To(ContainSubstring("+ bucket: (sensitive value)"))
		Expect(body).To(ContainSubstring("`bucket_name`"))
		Expect(body).NotTo(ContainSubstring("null_resource.noop"))

		Expect(body).NotTo(ContainSubstring("secret-bucket-name"))
		Expect(body).NotTo(ContainSubstring("secret-old-user-data"))
		Expect(body).NotTo(ContainSubstring("secret-new-user-data"))
		Expect(body).NotTo(ContainSubstring("secret-old-password"))
		Expect(body).NotTo(ContainSubstring("secret-new-password"))
		Expect(body).NotTo(ContainSubstring("secret-account-id"))
		Expect(body).NotTo(ContainSubstring("secret-noop-value"))
		Expect(body).NotTo(ContainSubstring("secret-output-bucket"))
		Expect(body).NotTo(ContainSubstring("secret-old-endpoint"))
		Expect(body).NotTo(ContainSubstring("secret-new-endpoint"))
		Expect(body).NotTo(ContainSubstring("secret-output-password"))
		Expect(body).NotTo(ContainSubstring("secret-unchanged-output"))
		Expect(body).NotTo(ContainSubstring("secret-queue-name"))
	})

	It("keeps the summary artifact independent of the report grouping", func() {
		report, err := parseTerraformShowPlan(readFixture("plan.json"), metadata())
		Expect(err).NotTo(HaveOccurred())

		summary := terraform.Summarize(report)
		Expect(summary).To(Equal(terraform.Summary{
			Create: 2,
			Update: 1,
			Delete: 1,
		}))

		Expect(terraform.RenderSummary(summary)).To(Equal(`{
  "create": 2,
  "update": 1,
  "delete": 1
}
`))
	})

	Describe("attribute diff", func() {
		update := terraform.ChangeUpdate

		It("masks a whole subtree that is sensitive on either side", func() {
			changes := resourceChanges(&tfjson.Change{
				Before:          map[string]any{"env": map[string]any{"A": "1"}, "port": float64(80)},
				After:           map[string]any{"env": map[string]any{"A": "2"}, "port": float64(80)},
				AfterSensitive:  map[string]any{"env": true},
				BeforeSensitive: map[string]any{},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "env", Action: update, Before: `{"A": "1"}`, After: "(sensitive value)"},
			}))
		})

		// the value did not change, but who may read it did, which Terraform reports
		// as a change too.
		It("shows a value that only became sensitive", func() {
			changes := resourceChanges(&tfjson.Change{
				Before:         map[string]any{"token": "x"},
				After:          map[string]any{"token": "x"},
				AfterSensitive: map[string]any{"token": true},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "token", Action: update, Before: `"x"`, After: "(sensitive value)"},
			}))
		})

		It("marks what is only known after apply, in lists and objects alike", func() {
			changes := resourceChanges(&tfjson.Change{
				Before: map[string]any{"ids": []any{"a", "b"}, "nested": map[string]any{"id": "1"}, "same": "x"},
				After:  map[string]any{"ids": []any{"a", nil}, "nested": map[string]any{}, "same": "x"},
				AfterUnknown: map[string]any{
					"ids":    []any{false, true},
					"nested": map[string]any{"id": true},
				},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "ids", Action: update, Before: `["a", "b"]`, After: `["a", (known after apply)]`},
				{Name: "nested", Action: update, Children: []terraform.Change{
					{Name: "id", Action: update, Before: `"1"`, After: "(known after apply)"},
				}},
			}))
		})

		It("writes a value that changed its shape on one line", func() {
			changes := resourceChanges(&tfjson.Change{
				Before: map[string]any{"value": "x"},
				After:  map[string]any{"value": map[string]any{"b": float64(2), "a": float64(1)}},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "value", Action: update, Before: `"x"`, After: `{"a": 1, "b": 2}`},
			}))
		})

		It("compares lists of objects element by element", func() {
			changes := resourceChanges(&tfjson.Change{
				Before: map[string]any{"rules": []any{map[string]any{"port": float64(80)}}},
				After: map[string]any{"rules": []any{
					map[string]any{"port": float64(80)},
					map[string]any{"port": float64(443)},
				}},
				ReplacePaths: []any{[]any{"rules", float64(1)}},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "rules", Action: update, Children: []terraform.Change{
					{Name: "[1]", Action: terraform.ChangeCreate, Note: "forces replacement", Children: []terraform.Change{
						{Name: "port", Action: terraform.ChangeCreate, After: "443"},
					}},
				}},
			}))
		})

		It("writes a deletion from what was there before", func() {
			changes := resourceChanges(&tfjson.Change{
				Before:          map[string]any{"name": "x", "secret": "y"},
				BeforeSensitive: map[string]any{"secret": true},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "name", Action: terraform.ChangeDelete, Before: `"x"`},
				{Name: "secret", Action: terraform.ChangeDelete, Before: "(sensitive value)"},
			}))
		})
	})
})
