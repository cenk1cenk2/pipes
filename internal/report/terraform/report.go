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
		// qualifies the action the resource is listed under, such as the order a
		// replacement runs in or what an import reads the resource from.
		Detail  string
		Changes []Change
	}

	Output struct {
		Name   string
		Fields []string
	}
)

// The actions that change a resource in place or replace it, whose resources are told
// apart by whether the plan changes a property of theirs at all.
var propertyChangeActions = []string{
	"update",
	"update-replacement",
	"replace",
	"create-replacement",
	"delete-replaced",
}

var unknowns = map[string]bool{}

// RegisterUnknowns names the markers a pipe writes in place of a value that is only
// known after apply. A secret is not one of them: a plan only pairs a secret with
// itself when its value changed.
func RegisterUnknowns(markers ...Marker) {
	for _, marker := range markers {
		unknowns[string(marker)] = true
	}
}

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

// HasPropertyChanges says whether the plan changes a property of the resource, which
// a change that only moves between values not known yet does not tell.
func (r Resource) HasPropertyChanges() bool {
	return slices.ContainsFunc(r.Changes, Change.changesProperty)
}

func (c Change) changesProperty() bool {
	if c.Note == ForcesReplacement {
		return true
	}

	if len(c.Children) > 0 {
		return slices.ContainsFunc(c.Children, Change.changesProperty)
	}

	pending := func(value string) bool {
		return value == "" || unknowns[value]
	}

	return !pending(c.Before) || !pending(c.After) || (!unknowns[c.Before] && !unknowns[c.After])
}

// What the template renders: the report with the attribute diff of every resource
// written out in full.
type (
	reportView struct {
		Report
		Actions []actionView
	}

	actionView struct {
		Action                 string
		Glyph                  string
		WithoutPropertyChanges bool
		Resources              []resourceView
		Outputs                []Output
	}

	resourceView struct {
		Resource
		Diff  string
		Fence string
	}
)

func RenderReport(report Report, layout DiffLayout) (string, error) {
	tmpl, err := template.New("report.md.gotmpl").
		Funcs(template.FuncMap{"cell": cell}).
		Parse(reportTemplate)
	if err != nil {
		return "", fmt.Errorf("parse report template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, newReportView(report, layout)); err != nil {
		return "", fmt.Errorf("render report template: %w", err)
	}

	return strings.TrimRight(body.String(), "\n") + "\n", nil
}

func newReportView(report Report, layout DiffLayout) reportView {
	view := reportView{Report: report}

	for _, action := range report.Actions {
		resources, unchanged := action.Resources, []Resource(nil)
		if slices.Contains(propertyChangeActions, action.Action) {
			resources = slices.DeleteFunc(slices.Clone(action.Resources), func(resource Resource) bool {
				return !resource.HasPropertyChanges()
			})
			unchanged = slices.DeleteFunc(slices.Clone(action.Resources), Resource.HasPropertyChanges)
		}

		if len(resources) > 0 || len(action.Outputs) > 0 || len(unchanged) == 0 {
			view.Actions = append(view.Actions, actionView{
				Action:    action.Action,
				Glyph:     Glyph(action.Action),
				Resources: newResourceViews(resources, layout),
				Outputs:   action.Outputs,
			})
		}

		if len(unchanged) > 0 {
			view.Actions = append(view.Actions, actionView{
				Action:                 action.Action,
				Glyph:                  Glyph(action.Action),
				WithoutPropertyChanges: true,
				Resources:              newResourceViews(unchanged, layout),
			})
		}
	}

	return view
}

func newResourceViews(resources []Resource, layout DiffLayout) []resourceView {
	views := make([]resourceView, 0, len(resources))
	for _, resource := range resources {
		diff := RenderChanges(resource.Changes, layout)
		if diff == "" {
			diff = NoAttributeChanges
		}

		views = append(views, resourceView{
			Resource: resource,
			Diff:     diff,
			Fence:    fence(diff),
		})
	}

	return views
}

// A pipe in a value would end its table cell early and shift the rest of the row; an
// escaped one stays part of the value, inside a code span as well.
func cell(value string) string {
	return strings.ReplaceAll(value, "|", `\|`)
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
