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
		report, err := parsePulumiPlanReport(readFixture("plan-unwrapped.json"), nil, metadata())
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
				{Name: "level", Action: terraform.ChangeUpdate, After: `"debug"`},
				{Name: "password", Action: terraform.ChangeUpdate, After: "[secret]"},
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
		report, err := parsePulumiPlanReport(readFixture("plan-versioned.json"), nil, metadata())
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
			}, nil)

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "members", Action: create, After: `["a", [secret]]`},
				{Name: "role", Action: create, After: `"urn:pulumi:dev::example::aws:iam/role:Role::task"`},
			}))
		})

		update := terraform.ChangeUpdate
		remove := terraform.ChangeDelete

		It("marks only the lines that changed in an update against the stack state", func() {
			changes := resourceChanges(&apitype.GoalV1{
				InputDiff: apitype.PlanDiffV1{Updates: map[string]any{
					"policy": "path \"a\" {\n  capabilities = [\"read\"]\n}\npath \"b\" {\n  capabilities = [\"read\"]\n}\n",
				}},
			}, map[string]any{
				"policy": "path \"a\" {\n  capabilities = [\"read\"]\n}\n",
			})

			Expect(changes).To(Equal([]terraform.Change{{
				Name:   "policy",
				Action: update,
				Before: "<<-EOT\n  path \"a\" {\n    capabilities = [\"read\"]\n  }\nEOT",
				After:  "<<-EOT\n  path \"a\" {\n    capabilities = [\"read\"]\n  }\n  path \"b\" {\n    capabilities = [\"read\"]\n  }\nEOT",
			}}))

			body, err := terraform.RenderReport(terraform.Report{Title: "x", Actions: []terraform.Action{{
				Action:    "update",
				Resources: []terraform.Resource{{Name: "policy", Changes: changes}},
			}}})
			Expect(err).NotTo(HaveOccurred())
			Expect(
				body,
			).To(ContainSubstring("~ policy: <<-EOT\n    path \"a\" {\n      capabilities = [\"read\"]\n    }\n+   path \"b\" {\n+     capabilities = [\"read\"]\n+   }\n  EOT"))
		})

		It("compares an updated object key by key against the stack state", func() {
			changes := resourceChanges(&apitype.GoalV1{
				InputDiff: apitype.PlanDiffV1{Updates: map[string]any{
					"data": map[string]any{"level": "debug", "mode": "a"},
				}},
			}, map[string]any{"data": map[string]any{"level": "info", "mode": "a", "old": "x"}})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "data", Action: update, Children: []terraform.Change{
					{Name: "level", Action: update, Before: `"info"`, After: `"debug"`},
					{Name: "old", Action: remove, Before: `"x"`},
				}},
			}))
		})

		It("shows what a value only known after apply replaces", func() {
			changes := resourceChanges(&apitype.GoalV1{
				InputDiff: apitype.PlanDiffV1{Updates: map[string]any{"arn": computedValuePlaceholder}},
			}, map[string]any{"arn": "arn:old"})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "arn", Action: update, Before: `"arn:old"`, After: "[unknown]"},
			}))
		})

		It("keeps an update whose sides read the same once masked", func() {
			secret := map[string]any{"4dabf18193072939515e22adb298388d": "1b47061264138c4ac30d75fd1eb44270", "ciphertext": "v1:x"}

			changes := resourceChanges(&apitype.GoalV1{
				InputDiff: apitype.PlanDiffV1{Updates: map[string]any{"token": secret}},
			}, map[string]any{"token": secret})

			Expect(changes).To(Equal([]terraform.Change{{Name: "token", Action: update, After: "[secret]"}}))
		})

		It("writes a deleted property from the value the stack state holds", func() {
			changes := resourceChanges(&apitype.GoalV1{
				InputDiff: apitype.PlanDiffV1{Deletes: []string{"immutable", "unknown"}},
			}, map[string]any{"immutable": true})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "immutable", Action: remove, Before: "true"},
				{Name: "unknown", Action: remove},
			}))
		})

		It("writes a deleted resource out from the stack state without the provider bookkeeping", func() {
			changes := resourceChanges(nil, map[string]any{
				"name":       "keel-system/default",
				"__defaults": []any{},
				"policies":   []any{"a"},
			})

			Expect(changes).To(Equal([]terraform.Change{
				{Name: "name", Action: remove, Before: `"keel-system/default"`},
				{Name: "policies", Action: remove, Before: `["a"]`},
			}))
		})

		It("falls back to the plan values when the stack state does not know the resource", func() {
			changes := resourceChanges(&apitype.GoalV1{
				InputDiff: apitype.PlanDiffV1{Updates: map[string]any{"runtime": "nodejs22.x"}},
			}, nil)

			Expect(changes).To(Equal([]terraform.Change{{Name: "runtime", Action: update, After: `"nodejs22.x"`}}))
		})

		It("carries nothing for a resource without a goal", func() {
			Expect(resourceChanges(nil, nil)).To(BeNil())
		})
	})
})

var _ = Describe("Pulumi stack state", func() {
	It("reads the inputs of every resource by urn and skips the ones pending deletion", func() {
		state, err := parseStackState([]byte(`{
			"version": 3,
			"deployment": {
				"manifest": {"time": "2026-10-01T12:00:00Z"},
				"resources": [
					{"urn": "urn:a", "inputs": {"name": "new"}, "outputs": {"id": "1"}},
					{"urn": "urn:a", "inputs": {"name": "old"}, "delete": true},
					{"urn": "urn:b", "inputs": {"policy": "x"}}
				]
			}
		}`))

		Expect(err).NotTo(HaveOccurred())
		Expect(state).To(Equal(map[string]map[string]any{
			"urn:a": {"name": "new"},
			"urn:b": {"policy": "x"},
		}))
	})

	It("fails on a state that is not JSON", func() {
		_, err := parseStackState([]byte("not json"))

		Expect(err).To(MatchError(ContainSubstring("parse Pulumi stack state")))
	})
})
