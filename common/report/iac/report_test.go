package iac

import (
	"strings"
	"testing"
)

func report(labels Labels) Report {
	return Report{
		Title:  "Example report",
		Labels: labels,
		Metadata: Metadata{
			Target:         "example",
			JobName:        "plan",
			JobUrl:         "https://gitlab.example.test/project/-/jobs/1",
			PipelineId:     "42",
			PipelineUrl:    "https://gitlab.example.test/project/-/pipelines/42",
			CommitShortSha: "01234567",
			ToolVersion:    "1.0.0",
		},
		Actions: []Action{
			{
				Action:    "create",
				Resources: []Resource{{Name: "one"}},
				Outputs:   []Output{{Name: "first"}},
			},
			{
				Action:    "move",
				Resources: []Resource{{Name: "two", PreviousName: "old-two"}},
			},
		},
	}
}

// Both pipes render through this template, so the section skeleton is the contract
// that keeps their reports readable side by side on one merge request.
func TestRenderMergeRequestReportKeepsOneStructureAcrossLabels(t *testing.T) {
	t.Parallel()

	skeleton := func(body string) []string {
		headings := []string{}
		for _, line := range strings.Split(body, "\n") {
			if strings.HasPrefix(line, "#") {
				headings = append(headings, strings.SplitN(line, " ", 2)[0])
			}
		}

		return headings
	}

	terraform, err := RenderMergeRequestReport(report(Labels{Target: "State", Outputs: "Outputs", ToolVersion: "Terraform version"}))
	if err != nil {
		t.Fatalf("render terraform flavored report: %v", err)
	}

	pulumi, err := RenderMergeRequestReport(report(Labels{Target: "Stack", Outputs: "Output properties", ToolVersion: "Pulumi version"}))
	if err != nil {
		t.Fatalf("render pulumi flavored report: %v", err)
	}

	left := strings.Join(skeleton(terraform), ",")
	right := strings.Join(skeleton(pulumi), ",")
	if left != right {
		t.Fatalf("expected both flavors to share a heading skeleton, got %q and %q", left, right)
	}

	if want := "##,###,###,###,####,####,###,####"; left != want {
		t.Fatalf("expected heading skeleton %q, got %q", want, left)
	}
}

func TestRenderMergeRequestReportUsesLabels(t *testing.T) {
	t.Parallel()

	body, err := RenderMergeRequestReport(report(Labels{Target: "Stack", Outputs: "Output properties", ToolVersion: "Pulumi version"}))
	if err != nil {
		t.Fatalf("render report: %v", err)
	}

	for _, expected := range []string{
		"## Example report",
		"| Stack | `example` |",
		"| Pulumi version | `1.0.0` |",
		"| Job | [plan](https://gitlab.example.test/project/-/jobs/1) |",
		"### Output properties",
		"Total planned actions: 2.",
		"- `two` (moved from `old-two`)",
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected rendered report to contain %q, got:\n%s", expected, body)
		}
	}
}

func TestRenderMergeRequestReportWithoutChanges(t *testing.T) {
	t.Parallel()

	body, err := RenderMergeRequestReport(Report{Title: "Example report"})
	if err != nil {
		t.Fatalf("render report: %v", err)
	}

	if !strings.Contains(body, "No changes detected.") {
		t.Fatalf("expected an empty report to say so, got:\n%s", body)
	}

	if strings.Contains(body, "### Metadata") {
		t.Fatalf("expected no metadata section without metadata, got:\n%s", body)
	}
}
