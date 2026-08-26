package preview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
)

type (
	actionAccumulator struct {
		Action      string
		Resources   map[string]iac.Resource
		OutputNames map[string][]string
	}
)

var pulumiPlanActionOrder = []string{
	"same",
	"read",
	"create",
	"import",
	"update",
	"replace",
	"create-replacement",
	"update-replacement",
	"read-replacement",
	"delete-replaced",
	"delete",
	"discard",
	"remove-pending-replace",
}

func parsePulumiPlanReport(data []byte, metadata iac.Metadata) (iac.Report, error) {
	plan, planVersion, err := parsePulumiPlan(data)
	if err != nil {
		return iac.Report{}, err
	}

	metadata.ToolVersion = plan.Manifest.Version
	metadata.PlanTime = formatManifestTime(plan.Manifest.Time)
	if planVersion != 0 {
		metadata.PlanSchema = fmt.Sprintf("%d", planVersion)
	}

	accumulators := map[string]*actionAccumulator{}

	for urn, resource := range plan.ResourcePlans {
		steps := resourceActions(resource)
		steps = slices.DeleteFunc(steps, func(action string) bool {
			return action == "same"
		})
		if len(steps) == 0 {
			continue
		}

		summary := resourceSummary(string(urn), resource)
		outputNames := resourceOutputNames(resource)

		for _, action := range steps {
			accumulator := accumulators[action]
			if accumulator == nil {
				accumulator = &actionAccumulator{
					Action:      action,
					Resources:   map[string]iac.Resource{},
					OutputNames: map[string][]string{},
				}
				accumulators[action] = accumulator
			}

			accumulator.Resources[summary.Id] = summary

			if len(outputNames) > 0 {
				accumulator.OutputNames[summary.Id] = outputNames
			}
		}
	}

	return iac.Report{
		Title: "Pulumi preview report",
		Labels: iac.Labels{
			Target:      "Stack",
			Outputs:     "Output properties",
			ToolVersion: "Pulumi version",
		},
		Metadata: metadata,
		Actions:  buildPulumiPlanActions(accumulators),
	}, nil
}

func parsePulumiPlan(data []byte) (apitype.DeploymentPlanV1, int, error) {
	var versioned apitype.VersionedDeploymentPlan
	if err := json.Unmarshal(data, &versioned); err != nil {
		return apitype.DeploymentPlanV1{}, 0, fmt.Errorf("parse Pulumi plan JSON: %w", err)
	}

	plan := bytes.TrimSpace(versioned.Plan)
	if len(plan) > 0 && !bytes.Equal(plan, []byte("null")) {
		deployment, err := parsePulumiDeploymentPlan(plan)

		return deployment, versioned.Version, err
	}

	deployment, err := parsePulumiDeploymentPlan(data)

	return deployment, 0, err
}

func parsePulumiDeploymentPlan(data []byte) (apitype.DeploymentPlanV1, error) {
	var plan apitype.DeploymentPlanV1
	if err := json.Unmarshal(data, &plan); err != nil {
		return apitype.DeploymentPlanV1{}, fmt.Errorf("parse Pulumi plan: %w", err)
	}

	return plan, nil
}

func formatManifestTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Format(time.RFC3339)
}

func resourceActions(resource apitype.ResourcePlanV1) []string {
	actions := make([]string, 0, len(resource.Steps))

	for _, action := range resource.Steps {
		action := strings.TrimSpace(string(action))
		if action == "" {
			continue
		}

		if slices.Contains(actions, action) {
			continue
		}

		actions = append(actions, action)
	}

	return actions
}

func resourceSummary(urn string, plan apitype.ResourcePlanV1) iac.Resource {
	resourceType := ""
	name := ""

	if plan.Goal != nil {
		resourceType = string(plan.Goal.Type)
		name = plan.Goal.Name
	}

	if resourceType == "" || name == "" {
		urnType, urnName := splitPulumiUrn(urn)
		if resourceType == "" {
			resourceType = urnType
		}

		if name == "" {
			name = urnName
		}
	}

	return iac.Resource{
		Name: pulumiResourceName(resourceType, name),
		Id:   urn,
	}
}

func pulumiResourceName(resourceType string, name string) string {
	switch {
	case resourceType != "" && name != "":
		return fmt.Sprintf("%s/%s", resourceType, name)
	case name != "":
		return name
	case resourceType != "":
		return resourceType
	}

	return "unknown resource"
}

func resourceOutputNames(resource apitype.ResourcePlanV1) []string {
	if resource.Goal == nil {
		return nil
	}

	diff := resource.Goal.OutputDiff

	return sortedUniqueStrings(slices.Concat(
		slices.Collect(maps.Keys(diff.Adds)),
		diff.Deletes,
		slices.Collect(maps.Keys(diff.Updates)),
	))
}

func buildPulumiPlanActions(accumulators map[string]*actionAccumulator) []iac.Action {
	actions := make([]iac.Action, 0, len(accumulators))

	for _, accumulator := range accumulators {
		resources := slices.Collect(maps.Values(accumulator.Resources))
		slices.SortFunc(resources, func(left, right iac.Resource) int {
			return strings.Compare(left.Id, right.Id)
		})

		outputs := make([]iac.Output, 0, len(accumulator.OutputNames))
		for _, urn := range slices.Sorted(maps.Keys(accumulator.OutputNames)) {
			outputs = append(outputs, iac.Output{
				Name:   accumulator.Resources[urn].Name,
				Fields: sortedUniqueStrings(accumulator.OutputNames[urn]),
			})
		}

		actions = append(actions, iac.Action{
			Action:    accumulator.Action,
			Resources: resources,
			Outputs:   outputs,
		})
	}

	iac.SortActions(actions, pulumiPlanActionOrder)

	return actions
}

func sortedUniqueStrings(values []string) []string {
	values = slices.Clone(values)
	for index, value := range values {
		values[index] = strings.TrimSpace(value)
	}
	values = slices.DeleteFunc(values, func(value string) bool {
		return value == ""
	})
	slices.Sort(values)

	return slices.Compact(values)
}

func splitPulumiUrn(urn string) (string, string) {
	parts := strings.Split(urn, "::")
	if len(parts) < 2 {
		return "", ""
	}

	return parts[len(parts)-2], parts[len(parts)-1]
}
