package iac

import (
	"bytes"
	"cmp"
	_ "embed"
	"fmt"
	"slices"
	"strings"
	"text/template"
)

//go:embed assets/mr-report.md.gotmpl
var mergeRequestReportTemplate string

// The plan report shared by the Terraform and Pulumi pipes. Other pipes report on
// entirely different things and are expected to bring their own model and template
// rather than bend this one.
type (
	Report struct {
		Title    string
		Labels   Labels
		Metadata Metadata
		Actions  []Action
	}

	// Names the tool-specific concepts the template renders, so both pipes produce
	// one document structure while still calling things what they are.
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

func RenderMergeRequestReport(report Report) (string, error) {
	tmpl, err := template.New("mr-report.md.gotmpl").
		Parse(mergeRequestReportTemplate)
	if err != nil {
		return "", fmt.Errorf("parse merge request report template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, report); err != nil {
		return "", fmt.Errorf("render merge request report template: %w", err)
	}

	return strings.TrimRight(body.String(), "\n") + "\n", nil
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
