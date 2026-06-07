package plan

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"maps"
	"slices"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

//go:embed assets/report.html.gotmpl
var reportTemplate string

type (
	terraformReport struct {
		Metadata       []reportMetadata
		Summary        []reportSummary
		ResourceGroups []reportGroup
		OutputGroups   []reportGroup
	}

	reportMetadata struct {
		Name  string
		Value string
	}

	reportSummary struct {
		Action        string
		ResourceCount int
		OutputCount   int
	}

	reportGroup struct {
		Action string
		Items  []string
	}

	terraformSummary struct {
		Create int `json:"create"`
		Update int `json:"update"`
		Delete int `json:"delete"`
	}
)

var reportActionOrder = []string{
	"create",
	"update",
	"delete",
	"replace",
	"read",
	"import",
	"forget",
	"unknown",
}

func parseTerraformShowPlan(output []byte) (*tfjson.Plan, error) {
	var plan tfjson.Plan

	decoder := json.NewDecoder(bytes.NewReader(output))
	if err := decoder.Decode(&plan); err != nil {
		return nil, fmt.Errorf("parse terraform show -json output: %w", err)
	}
	if err := plan.Validate(); err != nil {
		return nil, fmt.Errorf("validate terraform show -json output: %w", err)
	}

	return &plan, nil
}

func buildTerraformReport(plan *tfjson.Plan) *terraformReport {
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

	return &terraformReport{
		Summary:        reportSummaryItems(resources, outputs),
		ResourceGroups: reportGroups(resources),
		OutputGroups:   reportGroups(outputs),
	}
}

func summarizeTerraformPlan(plan *tfjson.Plan) terraformSummary {
	report := terraformSummary{}

	for _, change := range plan.ResourceChanges {
		if change == nil || change.Change == nil {
			continue
		}

		for _, action := range change.Change.Actions {
			switch action {
			case tfjson.ActionCreate:
				report.Create++
			case tfjson.ActionUpdate:
				report.Update++
			case tfjson.ActionDelete:
				report.Delete++
			}
		}
	}

	return report
}

func renderSummary(report terraformSummary) ([]byte, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("render Terraform summary report: %w", err)
	}

	return append(data, '\n'), nil
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

func reportSummaryItems(
	resources map[string][]string,
	outputs map[string][]string,
) []reportSummary {
	summary := []reportSummary{}

	for _, action := range reportActions(resources, outputs) {
		summary = append(summary, reportSummary{
			Action:        action,
			ResourceCount: len(resources[action]),
			OutputCount:   len(outputs[action]),
		})
	}

	return summary
}

func reportGroups(items map[string][]string) []reportGroup {
	groups := []reportGroup{}

	for _, action := range reportActions(items) {
		groups = append(groups, reportGroup{
			Action: action,
			Items:  slices.Sorted(slices.Values(items[action])),
		})
	}

	return groups
}

func reportActions(groups ...map[string][]string) []string {
	actions := []string{}
	for _, group := range groups {
		actions = slices.Concat(actions, slices.Collect(maps.Keys(group)))
	}
	slices.Sort(actions)
	actions = slices.Compact(actions)

	actionRank := func(action string) int {
		if rank := slices.Index(reportActionOrder, action); rank >= 0 {
			return rank
		}

		return len(reportActionOrder)
	}

	slices.SortFunc(actions, func(left, right string) int {
		if rank := cmp.Compare(actionRank(left), actionRank(right)); rank != 0 {
			return rank
		}

		return cmp.Compare(left, right)
	})

	return actions
}

func renderReport(report *terraformReport) ([]byte, error) {
	tmpl, err := template.New("report.html.gotmpl").
		Parse(reportTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse Terraform report template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, report); err != nil {
		return nil, fmt.Errorf("render Terraform report template: %w", err)
	}

	return body.Bytes(), nil
}
