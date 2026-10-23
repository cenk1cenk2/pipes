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

		resource := terraform.Resource{
			Name:    change.Address,
			Changes: resourceChanges(change.Change),
		}
		if moved {
			resource.PreviousName = change.PreviousAddress
		}

		details := []string{}
		if change.DeposedKey != "" {
			details = append(details, "deposed object "+change.DeposedKey)
		}

		switch {
		case change.Change.Actions.DestroyBeforeCreate():
			details = append(details, "destroy before create")
		case change.Change.Actions.CreateBeforeDestroy():
			details = append(details, "create before destroy")
		}

		if reason := actionReason(change); reason != "" {
			details = append(details, reason)
		}

		// an import and a move are operations of their own on top of the change the
		// actions carry, the way Terraform counts an import next to the update it
		// brings; a resource that only moved or is only imported carries no-op actions.
		importing := change.Change.Importing
		if importing != nil {
			if importing.ID != "" {
				details = append(details, "imported from "+importing.ID)
			} else {
				details = append(details, "imported by identity "+terraform.FormatValue(importing.Identity))
			}
		}

		resource.Detail = strings.Join(details, ", ")

		if importing != nil {
			resources["import"] = append(resources["import"], resource)
		}
		if moved {
			resources["move"] = append(resources["move"], resource)
		}

		if action != "no-op" {
			resources[action] = append(resources[action], resource)
		}
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

// actionReason is the note Terraform prints under a resource it replaces or destroys
// for a reason the attribute changes do not show.
func actionReason(change *tfjson.ResourceChange) string {
	switch change.ActionReason {
	case tfjson.ActionReasonReplaceBecauseTainted:
		return "tainted"
	case tfjson.ActionReasonReplaceByRequest:
		return "replaced as requested"
	case tfjson.ActionReasonReplaceByTriggers:
		return "replaced due to changes in replace_triggered_by"
	case tfjson.ActionReasonDeleteBecauseNoResourceConfig:
		return fmt.Sprintf("because %s.%s is not in configuration", change.Type, change.Name)
	case tfjson.ActionReasonDeleteBecauseNoMoveTarget:
		return fmt.Sprintf("because %s was moved to %s, which is not in configuration", change.PreviousAddress, change.Address)
	case tfjson.ActionReasonDeleteBecauseNoModule:
		return fmt.Sprintf("because %s is not in configuration", change.ModuleAddress)
	case tfjson.ActionReasonDeleteBecauseWrongRepetition:
		switch change.Index.(type) {
		case nil:
			return "because resource uses count or for_each"
		case float64:
			return "because resource does not use count"
		case string:
			return "because resource does not use for_each"
		}
	case tfjson.ActionReasonDeleteBecauseCountIndex:
		return fmt.Sprintf("because index [%s] is out of range for count", terraform.FormatValue(change.Index))
	case tfjson.ActionReasonDeleteBecauseEachKey:
		return fmt.Sprintf("because key [%s] is not in for_each map", terraform.FormatValue(change.Index))
	}

	return ""
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
