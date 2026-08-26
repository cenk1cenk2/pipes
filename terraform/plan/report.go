package plan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
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

func parseTerraformShowPlan(output []byte, metadata iac.Metadata) (iac.Report, error) {
	plan, err := decodeTerraformShowPlan(output)
	if err != nil {
		return iac.Report{}, err
	}

	metadata.ToolVersion = plan.TerraformVersion
	metadata.PlanTime = plan.Timestamp
	metadata.PlanSchema = plan.FormatVersion

	resources := map[string][]iac.Resource{}
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

		resource := iac.Resource{
			Name: change.Address,
		}
		if moved {
			resource.PreviousName = change.PreviousAddress
		}

		resources[action] = append(resources[action], resource)
	}

	outputs := map[string][]iac.Output{}
	for name, change := range plan.OutputChanges {
		if change == nil {
			continue
		}

		action := terraformChangeAction(change.Actions)
		if action == "no-op" {
			continue
		}

		outputs[action] = append(outputs[action], iac.Output{Name: name})
	}

	return iac.Report{
		Title: "Terraform plan report",
		Labels: iac.Labels{
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

	decoder := json.NewDecoder(bytes.NewReader(output))
	if err := decoder.Decode(&plan); err != nil {
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
	resources map[string][]iac.Resource,
	outputs map[string][]iac.Output,
) []iac.Action {
	names := slices.Concat(slices.Collect(maps.Keys(resources)), slices.Collect(maps.Keys(outputs)))
	slices.Sort(names)

	actions := []iac.Action{}
	for _, name := range slices.Compact(names) {
		action := iac.Action{
			Action:    name,
			Resources: slices.Clone(resources[name]),
			Outputs:   slices.Clone(outputs[name]),
		}

		slices.SortFunc(action.Resources, func(left, right iac.Resource) int {
			return strings.Compare(left.Name, right.Name)
		})
		slices.SortFunc(action.Outputs, func(left, right iac.Output) int {
			return strings.Compare(left.Name, right.Name)
		})

		actions = append(actions, action)
	}

	iac.SortActions(actions, mergeRequestReportActionOrder)

	return actions
}
