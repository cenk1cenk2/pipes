package plan

import (
	json "encoding/json/v2"
	"fmt"
	"maps"
	"slices"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

var mergeRequestReportActionOrder = []string{
	"create",
	"update",
	"delete",
	"replace",
	"move",
	"read",
	"import",
	"forget",
	"unknown",
}

func parseTerraformShowPlan(output []byte, metadata terraform.Metadata) (terraform.Report, error) {
	plan, err := decodeTerraformShowPlan(output)
	if err != nil {
		return terraform.Report{}, err
	}

	metadata.ToolVersion = plan.TerraformVersion
	metadata.PlanTime = plan.Timestamp
	metadata.PlanSchema = plan.FormatVersion

	resources := map[string][]terraform.Resource{}
	for _, change := range plan.ResourceChanges {
		if change == nil || change.Change == nil {
			continue
		}

		action := terraformChangeAction(change.Change.Actions)
		moved := change.PreviousAddress != "" && change.PreviousAddress != change.Address

		// A resource that only moved carries no-op actions, which would otherwise drop
		// it from the report entirely.
		if action == "no-op" {
			if !moved {
				continue
			}

			action = "move"
		}

		resource := terraform.Resource{
			Name: change.Address,
		}
		if moved {
			resource.PreviousName = change.PreviousAddress
		}

		resources[action] = append(resources[action], resource)
	}

	outputs := map[string][]terraform.Output{}
	for name, change := range plan.OutputChanges {
		if change == nil {
			continue
		}

		action := terraformChangeAction(change.Actions)
		if action == "no-op" {
			continue
		}

		outputs[action] = append(outputs[action], terraform.Output{Name: name})
	}

	return terraform.Report{
		Title: "Terraform plan report",
		Labels: terraform.Labels{
			Target:      "State",
			Outputs:     "Outputs",
			ToolVersion: "Terraform version",
		},
		Metadata: metadata,
		Actions:  mergeRequestReportActions(resources, outputs),
	}, nil
}

func decodeTerraformShowPlan(output []byte) (tfjson.Plan, error) {
	var plan tfjson.Plan

	if err := json.Unmarshal(output, &plan, json.RejectUnknownMembers(true)); err != nil {
		return tfjson.Plan{}, fmt.Errorf("parse terraform show -json output: %w", err)
	}
	if err := plan.Validate(); err != nil {
		return tfjson.Plan{}, fmt.Errorf("validate terraform show -json output: %w", err)
	}

	return plan, nil
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

func mergeRequestReportActions(
	resources map[string][]terraform.Resource,
	outputs map[string][]terraform.Output,
) []terraform.Action {
	names := slices.Concat(slices.Collect(maps.Keys(resources)), slices.Collect(maps.Keys(outputs)))
	slices.Sort(names)

	actions := []terraform.Action{}
	for _, name := range slices.Compact(names) {
		action := terraform.Action{
			Action:    name,
			Resources: slices.Clone(resources[name]),
			Outputs:   slices.Clone(outputs[name]),
		}

		slices.SortFunc(action.Resources, func(left, right terraform.Resource) int {
			return strings.Compare(left.Name, right.Name)
		})
		slices.SortFunc(action.Outputs, func(left, right terraform.Output) int {
			return strings.Compare(left.Name, right.Name)
		})

		actions = append(actions, action)
	}

	terraform.SortActions(actions, mergeRequestReportActionOrder)

	return actions
}
