// Package terraform writes the GitLab artifacts:reports:terraform JSON, the job log
// copy of the plan and the merge request note that goes with them. Pulumi previews
// report through it too, since GitLab has no report kind of their own.
package terraform

import (
	"bytes"
	"cmp"
	_ "embed"
	"fmt"
	"slices"
	"strings"
	"text/template"
)

//go:embed assets/report.md.gotmpl
var reportTemplate string

// The plan report shared by the Terraform and Pulumi pipes; a pipe reporting on
// something else brings its own model and template.
type (
	Report struct {
		Title    string
		Labels   Labels
		Metadata Metadata
		Actions  []Action
	}

	// names the tool-specific concepts the template renders, so both pipes produce one document structure.
	Labels struct {
		Target      string
		Outputs     string
		ToolVersion string
	}

	Metadata struct {
		Target         string
		Cwd            string
		JobName        string
		JobUrl         string
		PipelineId     string
		PipelineUrl    string
		CommitSha      string
		CommitShortSha string
		ToolVersion    string
		PlanTime       string
		PlanSchema     string
	}

	Action struct {
		Action    string
		Resources []Resource
		Outputs   []Output
	}

	Resource struct {
		Name         string
		Id           string
		PreviousName string
		Changes      []Change
	}

	Output struct {
		Name   string
		Fields []string
	}
)

func (m Metadata) Any() bool {
	return m != Metadata{}
}

func (r Report) Total() int {
	total := 0
	for _, action := range r.Actions {
		total += len(action.Resources)
	}

	return total
}

func (r Report) HasResources() bool {
	return slices.ContainsFunc(r.Actions, func(action Action) bool {
		return len(action.Resources) > 0
	})
}

func (r Report) HasOutputs() bool {
	return slices.ContainsFunc(r.Actions, func(action Action) bool {
		return len(action.Outputs) > 0
	})
}

// What the template renders: the report with the attribute diff of every resource
// written out in full.
type (
	reportView struct {
		Report
		Actions []actionView
	}

	actionView struct {
		Action    string
		Resources []resourceView
		Outputs   []Output
	}

	resourceView struct {
		Resource
		Glyph string
		Diff  string
		Fence string
	}
)

func RenderReport(report Report) (string, error) {
	tmpl, err := template.New("report.md.gotmpl").
		Parse(reportTemplate)
	if err != nil {
		return "", fmt.Errorf("parse report template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, newReportView(report)); err != nil {
		return "", fmt.Errorf("render report template: %w", err)
	}

	return strings.TrimRight(body.String(), "\n") + "\n", nil
}

func newReportView(report Report) reportView {
	view := reportView{Report: report}

	for _, action := range report.Actions {
		resources := make([]resourceView, 0, len(action.Resources))
		for _, resource := range action.Resources {
			diff := RenderChanges(resource.Changes)
			if diff == "" {
				diff = NoAttributeChanges
			}

			resources = append(resources, resourceView{
				Resource: resource,
				Glyph:    Glyph(action.Action),
				Diff:     diff,
				Fence:    fence(diff),
			})
		}

		view.Actions = append(view.Actions, actionView{
			Action:    action.Action,
			Resources: resources,
			Outputs:   action.Outputs,
		})
	}

	return view
}

// A value can carry a run of backticks of its own, which would close the code block
// early and let the rest of the diff render as markup.
func fence(body string) string {
	longest := 0
	run := 0
	for _, char := range body {
		if char == '`' {
			run++
			longest = max(longest, run)

			continue
		}

		run = 0
	}

	return strings.Repeat("`", max(3, longest+1))
}

// Orders actions by the tool's own step vocabulary, with anything unrecognized
// sorted alphabetically after the known ones.
func SortActions(actions []Action, order []string) {
	rank := func(action string) int {
		if index := slices.Index(order, action); index >= 0 {
			return index
		}

		return len(order)
	}

	slices.SortFunc(actions, func(left, right Action) int {
		if compared := cmp.Compare(rank(left.Action), rank(right.Action)); compared != 0 {
			return compared
		}

		return cmp.Compare(left.Action, right.Action)
	})
}
