package plan

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/json"
	"fmt"
	"maps"
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

	terraformSummary struct {
		Create int `json:"create"`
		Update int `json:"update"`
		Delete int `json:"delete"`
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

func summarizeTerraformShowPlan(output []byte) (terraformSummary, error) {
	var plan tfjson.Plan

	decoder := json.NewDecoder(bytes.NewReader(output))
	if err := decoder.Decode(&plan); err != nil {
		return terraformSummary{}, fmt.Errorf("parse terraform show -json output: %w", err)
	}
	if err := plan.Validate(); err != nil {
		return terraformSummary{}, fmt.Errorf("validate terraform show -json output: %w", err)
	}

	summary := terraformSummary{}
	for _, change := range plan.ResourceChanges {
		if change == nil || change.Change == nil {
			continue
		}

		for _, action := range change.Change.Actions {
			switch action {
			case tfjson.ActionCreate:
				summary.Create++
			case tfjson.ActionUpdate:
				summary.Update++
			case tfjson.ActionDelete:
				summary.Delete++
			}
		}
	}

	return summary, nil
}

func renderSummary(summary terraformSummary) ([]byte, error) {
	body, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("render terraform summary: %w", err)
	}

	return append(body, '\n'), nil
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
		groups = append(groups, mergeRequestReportGroup{
			Action: action,
			Items:  slices.Sorted(slices.Values(items[action])),
		})
	}

	return groups
}

func mergeRequestReportActions(groups ...map[string][]string) []string {
	actions := []string{}
	for _, group := range groups {
		actions = slices.Concat(actions, slices.Collect(maps.Keys(group)))
	}
	slices.Sort(actions)
	actions = slices.Compact(actions)

	actionRank := func(action string) int {
		if rank := slices.Index(mergeRequestReportActionOrder, action); rank >= 0 {
			return rank
		}

		return len(mergeRequestReportActionOrder)
	}

	slices.SortFunc(actions, func(left, right string) int {
		if rank := cmp.Compare(actionRank(left), actionRank(right)); rank != 0 {
			return rank
		}

		return cmp.Compare(left, right)
	})

	return actions
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
