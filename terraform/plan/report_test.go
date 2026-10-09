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

		body, err := terraform.RenderReport(report, terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("## Terraform plan report"))
		Expect(body).To(ContainSubstring("| State | `production` |"))
		Expect(body).To(ContainSubstring("| Terraform version | `1.9.8` |"))
		Expect(body).To(ContainSubstring("[tf-plan](https://gitlab.example.test/project/-/jobs/1)"))
		Expect(body).To(ContainSubstring("Total planned actions: 5."))
		Expect(body).To(ContainSubstring("<summary><code>+ create</code> <code>aws_s3_bucket.logs</code></summary>"))
		Expect(body).To(ContainSubstring(
			"<summary><code>&gt; move</code> <code>aws_sqs_queue.jobs</code> (moved from <code>aws_sqs_queue.legacy_jobs</code>)</summary>",
		))
		Expect(body).To(ContainSubstring(
			"<summary><code>-/+ replace</code> <code>aws_db_instance.main</code> (destroy before create)</summary>\n\n```diff\n-/+ replace\n",
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

	DescribeTable("classifies every resource operation the plan carries",
		func(change string, actions []string, detail string, summary terraform.Summary) {
			report, err := parseTerraformShowPlan([]byte(`{
				"format_version": "1.2",
				"resource_changes": [{"address": "aws_instance.web", `+change+`}]
			}`), metadata())
			Expect(err).NotTo(HaveOccurred())

			names := []string{}
			for _, action := range report.Actions {
				names = append(names, action.Action)
			}

			Expect(names).To(Equal(actions))
			Expect(resource(report, "aws_instance.web").Detail).To(Equal(detail))
			Expect(terraform.Summarize(report)).To(Equal(summary))
		},
		Entry("a create", `"change": {"actions": ["create"]}`,
			[]string{"create"}, "", terraform.Summary{Create: 1}),
		Entry("an update", `"change": {"actions": ["update"]}`,
			[]string{"update"}, "", terraform.Summary{Update: 1}),
		Entry("a delete", `"change": {"actions": ["delete"]}`,
			[]string{"delete"}, "", terraform.Summary{Delete: 1}),
		Entry("a replacement that destroys first", `"change": {"actions": ["delete", "create"]}`,
			[]string{"replace"}, "destroy before create", terraform.Summary{Create: 1, Delete: 1}),
		Entry("a replacement that creates first", `"change": {"actions": ["create", "delete"]}`,
			[]string{"replace"}, "create before destroy", terraform.Summary{Create: 1, Delete: 1}),
		Entry("a replacement that forgets the prior object", `"change": {"actions": ["create", "forget"]}`,
			[]string{"create+forget"}, "", terraform.Summary{Create: 1}),
		Entry("an import", `"change": {"actions": ["no-op"], "importing": {"id": "i-123"}}`,
			[]string{"import"}, "imported from i-123", terraform.Summary{}),
		Entry("an import by identity", `"change": {"actions": ["no-op"], "importing": {"identity": {"id": "i-123"}}}`,
			[]string{"import"}, `imported by identity {"id": "i-123"}`, terraform.Summary{}),
		Entry("an import that updates", `"change": {"actions": ["update"], "importing": {"id": "i-123"}}`,
			[]string{"update", "import"}, "imported from i-123", terraform.Summary{Update: 1}),
		Entry("a move", `"previous_address": "aws_instance.old", "change": {"actions": ["no-op"]}`,
			[]string{"move"}, "", terraform.Summary{}),
		Entry("a move that updates", `"previous_address": "aws_instance.old", "change": {"actions": ["update"]}`,
			[]string{"update", "move"}, "", terraform.Summary{Update: 1}),
		Entry("a read", `"change": {"actions": ["read"]}`,
			[]string{"read"}, "", terraform.Summary{}),
		Entry("a forget", `"change": {"actions": ["forget"]}`,
			[]string{"forget"}, "", terraform.Summary{}),
	)

	It("labels an imported resource with its operation and where it comes from", func() {
		report, err := parseTerraformShowPlan([]byte(`{
			"format_version": "1.2",
			"resource_changes": [{
				"address": "aws_instance.web",
				"change": {
					"actions": ["update"],
					"before": {"instance_type": "t3.micro"},
					"after": {"instance_type": "t3.small"},
					"importing": {"id": "i-123"}
				}
			}]
		}`), metadata())
		Expect(err).NotTo(HaveOccurred())

		body, err := terraform.RenderReport(report, terraform.DiffPlan)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("| `= import` | 1 | 0 |"))
		Expect(body).To(ContainSubstring("#### `= import` (1)"))
		Expect(body).To(ContainSubstring(
			"<summary><code>= import</code> <code>aws_instance.web</code> (imported from i-123)</summary>\n\n```diff\n= import\n~ instance_type: \"t3.micro\" -> \"t3.small\"\n```",
		))
		Expect(body).To(ContainSubstring("<summary><code>~ update</code> <code>aws_instance.web</code> (imported from i-123)</summary>"))
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

		It("writes out what a value only known after apply replaces, attribute by attribute", func() {
			changes := resourceChanges(&tfjson.Change{
				Before: map[string]any{"metadata": map[string]any{
					"notes":    "line one\nline two\n",
					"revision": float64(146),
					"values":   `{"global":{"level":"error"},"cm":{"policy.csv":"g, admin\np, mcp\n"}}`,
				}},
				After:        map[string]any{},
				AfterUnknown: map[string]any{"metadata": true},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "metadata", Action: update, Children: []terraform.Change{
					{Name: "notes", Action: update, Before: "<<-EOT\n  line one\n  line two\nEOT", After: "(known after apply)"},
					{Name: "revision", Action: update, Before: "146", After: "(known after apply)"},
					{Name: "values", Action: update, Note: terraform.JSONEncoded, Children: []terraform.Change{
						{Name: "cm", Action: update, Children: []terraform.Change{
							{Name: "policy.csv", Action: update, Before: "<<-EOT\n  g, admin\n  p, mcp\nEOT", After: "(known after apply)"},
						}},
						{Name: "global", Action: update, Children: []terraform.Change{
							{Name: "level", Action: update, Before: `"error"`, After: "(known after apply)"},
						}},
					}},
				}},
			}))
		})

		It("compares a string holding a JSON document as the document", func() {
			changes := resourceChanges(&tfjson.Change{
				Before: map[string]any{"values": `{"replicas": 2, "url": "a"}`},
				After:  map[string]any{"values": `{"replicas": 3, "url": "a"}`},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "values", Action: update, Note: terraform.JSONEncoded, Children: []terraform.Change{
					{Name: "replicas", Action: update, Before: "2", After: "3"},
				}},
			}))
		})

		It("keeps a sensitive string holding a JSON document masked", func() {
			changes := resourceChanges(&tfjson.Change{
				Before:          map[string]any{"values": `{"password": "a"}`},
				After:           map[string]any{"values": `{"password": "b"}`},
				BeforeSensitive: map[string]any{"values": true},
				AfterSensitive:  map[string]any{"values": true},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "values", Action: update, Before: "(sensitive value)", After: "(sensitive value)"},
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
