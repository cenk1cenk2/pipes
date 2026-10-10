package terraform_test

import (
	"slices"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

var (
	sensitive = terraform.Marker("(sensitive value)")
	computed  = terraform.Marker("(known after apply)")
	secret    = terraform.Marker("[secret]")
	unknown   = terraform.Marker("[unknown]")
	asset     = terraform.Marker("[asset]")
)

var _ = BeforeSuite(func() {
	terraform.RegisterUnknowns(computed, unknown)
})

var _ = Describe("Resources without property changes", func() {
	update := terraform.ChangeUpdate

	real := terraform.Change{Name: "size", Action: update, Before: "1", After: "2"}
	pending := terraform.Change{Name: "arn", Action: terraform.ChangeCreate, After: string(computed)}

	DescribeTable("tells whether the plan changes a property of the resource",
		func(changes []terraform.Change, expected bool) {
			Expect(terraform.Resource{Name: "one", Changes: changes}.HasPropertyChanges()).To(Equal(expected))
		},
		Entry("no changes at all", nil, false),
		Entry("a value only known after apply", []terraform.Change{pending}, false),
		Entry("a value unknown on both sides", []terraform.Change{
			{Name: "arn", Action: update, Before: string(computed), After: string(computed)},
		}, false),
		Entry("a value unknown to another tool's marker", []terraform.Change{
			{Name: "arn", Action: update, After: string(unknown)},
		}, false),
		Entry("a secret that changed on both sides", []terraform.Change{
			{Name: "password", Action: update, Before: string(sensitive), After: string(sensitive)},
		}, true),
		Entry("a pulumi secret that changed on both sides", []terraform.Change{
			{Name: "token", Action: update, Before: string(secret), After: string(secret)},
		}, true),
		Entry("a secret the plan only shows the new side of", []terraform.Change{
			{Name: "token", Action: update, After: string(secret)},
		}, true),
		Entry("a secret next to a visible value", []terraform.Change{
			{Name: "password", Action: update, Before: string(sensitive), After: `"b"`},
		}, true),
		Entry("a visible value turned into a secret", []terraform.Change{
			{Name: "password", Action: update, Before: `"a"`, After: string(sensitive)},
		}, true),
		Entry("a secret that becomes known after apply", []terraform.Change{
			{Name: "password", Action: update, Before: string(sensitive), After: string(computed)},
		}, true),
		Entry("a visible value that becomes known after apply", []terraform.Change{
			{Name: "subnet_id", Action: update, Before: `"subnet-old"`, After: string(computed)},
		}, true),
		Entry("an asset", []terraform.Change{
			{Name: "source", Action: update, Before: string(asset), After: string(asset)},
		}, true),
		Entry("an attribute that forces the replacement while pending", []terraform.Change{
			{Name: "password", Action: update, Before: string(sensitive), After: string(computed), Note: terraform.ForcesReplacement},
			pending,
		}, true),
		Entry("an unknown attribute that forces the replacement", []terraform.Change{
			{Name: "arn", Action: update, Before: string(computed), After: string(computed), Note: terraform.ForcesReplacement},
		}, true),
		Entry("values known after apply nested in an object", []terraform.Change{
			{Name: "data", Action: update, Children: []terraform.Change{pending, pending}},
		}, false),
		Entry("a real value next to a pending one", []terraform.Change{pending, real}, true),
		Entry("a real value nested in an object", []terraform.Change{
			{Name: "data", Action: update, Children: []terraform.Change{pending, real}},
		}, true),
		Entry("a property removed without its value", []terraform.Change{
			{Name: "immutable", Action: terraform.ChangeDelete},
		}, true),
	)

	actionRows := func(body string) []string {
		return slices.DeleteFunc(strings.Split(body, "\n"), func(line string) bool {
			return !strings.HasPrefix(line, "| `")
		})
	}

	report := func(actions ...terraform.Action) terraform.Report {
		return terraform.Report{Title: "Example report", Actions: actions}
	}

	It("lists the updates without property changes as a group of their own", func() {
		body, err := terraform.RenderReport(report(terraform.Action{
			Action: "update",
			Resources: []terraform.Resource{
				{Name: "changed", Changes: []terraform.Change{real, pending}},
				{Name: "pending", Changes: []terraform.Change{pending}},
				{Name: "empty"},
			},
			Outputs: []terraform.Output{{Name: "endpoint"}},
		}), terraform.DiffUnified)
		Expect(err).NotTo(HaveOccurred())

		Expect(actionRows(body)).To(Equal([]string{"| `~ update` | 1 | 1 |", "| `~ update` without property changes | 2 | 0 |"}))
		Expect(body).To(ContainSubstring("Total planned actions: 3."))
		Expect(body).To(ContainSubstring("#### `~ update` (1)\n\n<details>\n<summary><code>~ update</code> <code>changed</code></summary>"))
		Expect(body).To(ContainSubstring(
			"#### `~ update` without property changes (2)\n\n<details>\n<summary><code>~ update</code> <code>pending</code></summary>",
		))
		Expect(body).To(ContainSubstring("<summary><code>~ update</code> <code>empty</code></summary>"))
		Expect(body).To(ContainSubstring("#### `~ update` (1)\n- `endpoint`"))
	})

	It("keeps the row of an action whose resources all change nothing when it carries outputs", func() {
		body, err := terraform.RenderReport(report(terraform.Action{
			Action:    "update",
			Resources: []terraform.Resource{{Name: "pending", Changes: []terraform.Change{pending}}},
			Outputs:   []terraform.Output{{Name: "endpoint"}},
		}), terraform.DiffUnified)
		Expect(err).NotTo(HaveOccurred())

		Expect(actionRows(body)).To(Equal([]string{"| `~ update` | 0 | 1 |", "| `~ update` without property changes | 1 | 0 |"}))
		Expect(body).NotTo(ContainSubstring("#### `~ update` (0)"))
	})

	It("drops the row of an action whose resources all change nothing", func() {
		body, err := terraform.RenderReport(report(terraform.Action{
			Action:    "create-replacement",
			Resources: []terraform.Resource{{Name: "pending", Changes: []terraform.Change{pending}}},
		}), terraform.DiffUnified)
		Expect(err).NotTo(HaveOccurred())

		Expect(actionRows(body)).To(Equal([]string{"| `+ create-replacement` without property changes | 1 | 0 |"}))
	})

	DescribeTable("groups every action that changes a resource in place or replaces it",
		func(action string, layout terraform.DiffLayout) {
			body, err := terraform.RenderReport(report(terraform.Action{
				Action: action,
				Resources: []terraform.Resource{
					{Name: "changed", Changes: []terraform.Change{real}},
					{Name: "pending", Changes: []terraform.Change{pending}},
					{Name: "forced", Changes: []terraform.Change{
						{Name: "arn", Action: update, Before: string(computed), After: string(computed), Note: terraform.ForcesReplacement},
					}},
				},
			}), layout)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(MatchRegexp("#### `[^`]+ " + action + "` \\(2\\)\n"))
			Expect(body).To(MatchRegexp("#### `[^`]+ " + action + "` without property changes \\(1\\)\n"))
		},
		Entry(nil, "update", terraform.DiffUnified),
		Entry(nil, "update", terraform.DiffPlan),
		Entry(nil, "update-replacement", terraform.DiffUnified),
		Entry(nil, "replace", terraform.DiffUnified),
		Entry(nil, "replace", terraform.DiffPlan),
		Entry(nil, "create-replacement", terraform.DiffUnified),
		Entry(nil, "delete-replaced", terraform.DiffUnified),
		Entry(nil, "delete-replaced", terraform.DiffPlan),
	)

	DescribeTable("leaves every other action as it is",
		func(action string) {
			body, err := terraform.RenderReport(report(terraform.Action{
				Action:    action,
				Resources: []terraform.Resource{{Name: "pending", Changes: []terraform.Change{pending}}, {Name: "empty"}},
			}), terraform.DiffUnified)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).NotTo(ContainSubstring("without property changes"))
		},
		Entry(nil, "create"),
		Entry(nil, "delete"),
		Entry(nil, "import"),
		Entry(nil, "move"),
		Entry(nil, "read"),
		Entry(nil, "forget"),
	)

	It("counts the updates without property changes next to every update", func() {
		summary := terraform.Summarize(report(
			terraform.Action{Action: "create", Resources: []terraform.Resource{{Name: "pending", Changes: []terraform.Change{pending}}}},
			terraform.Action{Action: "update", Resources: []terraform.Resource{
				{Name: "changed", Changes: []terraform.Change{real}},
				{Name: "pending", Changes: []terraform.Change{pending}},
				{Name: "empty"},
			}},
			terraform.Action{Action: "delete", Resources: []terraform.Resource{{Name: "gone"}}},
			terraform.Action{Action: "update-replacement", Resources: []terraform.Resource{{Name: "pending", Changes: []terraform.Change{pending}}}},
			terraform.Action{Action: "replace", Resources: []terraform.Resource{{Name: "pending", Changes: []terraform.Change{pending}}}},
		))

		Expect(summary).To(Equal(terraform.Summary{Create: 2, Update: 4, Delete: 2, UpdateWithoutPropertyChanges: 3}))
		Expect(terraform.RenderSummary(summary)).To(MatchJSON(`{
  "create": 2,
  "update": 4,
  "delete": 2,
  "update_without_property_changes": 3
}
`))
	})

	It("writes the summary as it always has while every update changes a property", func() {
		summary := terraform.Summarize(report(
			terraform.Action{Action: "update", Resources: []terraform.Resource{{Name: "changed", Changes: []terraform.Change{real}}}},
		))

		Expect(terraform.RenderSummary(summary)).To(MatchJSON(`{
  "create": 0,
  "update": 1,
  "delete": 0
}
`))
	})
})
