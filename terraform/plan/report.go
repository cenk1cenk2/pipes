package plan

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
)

//go:embed assets/mr-report.md.gotmpl
var mergeRequestReportTemplate string

type (
	mergeRequestReport struct {
		Metadata       []mergeRequestReportMetadata
		Summary        []mergeRequestReportSummary
		ResourceGroups []mergeRequestReportGroup
		OutputGroups   []mergeRequestReportGroup
	}

	mergeRequestReportMetadata struct {
		Name  string
		Value string
	}

	mergeRequestReportSummary struct {
		Action        string
		ResourceCount int
		OutputCount   int
	}

	mergeRequestReportGroup struct {
		Action string
		Items  []string
	}
)

var mergeRequestReportActionOrder = []string{
	"create",
	"update",
	"delete",
	"replace",
	"read",
	"import",
	"forget",
	"unknown",
}

func parseTerraformShowPlan(output []byte) (*mergeRequestReport, error) {
	var plan tfjson.Plan

	decoder := json.NewDecoder(bytes.NewReader(output))
	if err := decoder.Decode(&plan); err != nil {
		return nil, fmt.Errorf("parse terraform show -json output: %w", err)
	}
	if err := plan.Validate(); err != nil {
		return nil, fmt.Errorf("validate terraform show -json output: %w", err)
	}

	resources := map[string][]string{}
	for _, change := range plan.ResourceChanges {
		if change == nil || change.Change == nil {
			continue
		}

		action := terraformChangeAction(change.Change.Actions)
		if action == "no-op" {
			continue
		}

		resources[action] = append(resources[action], change.Address)
	}

	outputs := map[string][]string{}
	for name, change := range plan.OutputChanges {
		if change == nil {
			continue
		}

		action := terraformChangeAction(change.Actions)
		if action == "no-op" {
			continue
		}

		outputs[action] = append(outputs[action], name)
	}

	return &mergeRequestReport{
		Summary:        mergeRequestReportSummaryItems(resources, outputs),
		ResourceGroups: mergeRequestReportGroups(resources),
		OutputGroups:   mergeRequestReportGroups(outputs),
	}, nil
}

func terraformChangeAction(actions tfjson.Actions) string {
	if len(actions) == 0 {
		return "unknown"
	}

	if actions.Replace() {
		return "replace"
	}

	names := make([]string, 0, len(actions))
	for _, action := range actions {
		names = append(names, string(action))
	}

	return strings.Join(names, "+")
}

func mergeRequestReportSummaryItems(
	resources map[string][]string,
	outputs map[string][]string,
) []mergeRequestReportSummary {
	summary := []mergeRequestReportSummary{}

	for _, action := range mergeRequestReportActions(resources, outputs) {
		summary = append(summary, mergeRequestReportSummary{
			Action:        action,
			ResourceCount: len(resources[action]),
			OutputCount:   len(outputs[action]),
		})
	}

	return summary
}

func mergeRequestReportGroups(items map[string][]string) []mergeRequestReportGroup {
	groups := []mergeRequestReportGroup{}

	for _, action := range mergeRequestReportActions(items) {
		actionItems := slices.Clone(items[action])
		slices.Sort(actionItems)

		groups = append(groups, mergeRequestReportGroup{
			Action: action,
			Items:  actionItems,
		})
	}

	return groups
}

func mergeRequestReportActions(groups ...map[string][]string) []string {
	seen := map[string]struct{}{}

	for _, group := range groups {
		for action, items := range group {
			if len(items) == 0 {
				continue
			}

			seen[action] = struct{}{}
		}
	}

	actions := []string{}
	for _, action := range mergeRequestReportActionOrder {
		if _, ok := seen[action]; ok {
			actions = append(actions, action)
			delete(seen, action)
		}
	}

	remaining := make([]string, 0, len(seen))
	for action := range seen {
		remaining = append(remaining, action)
	}

	slices.Sort(remaining)

	return append(actions, remaining...)
}

func renderMergeRequestReport(report *mergeRequestReport) (string, error) {
	tmpl, err := template.New("mr-report.md.gotmpl").
		Parse(mergeRequestReportTemplate)
	if err != nil {
		return "", fmt.Errorf("parse GitLab merge request report template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, report); err != nil {
		return "", fmt.Errorf("render GitLab merge request report template: %w", err)
	}

	return body.String(), nil
}
