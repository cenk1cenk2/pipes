package preview

import (
	"bytes"
	json "encoding/json/v2"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

type (
	actionAccumulator struct {
		Action      string
		Resources   map[string]terraform.Resource
		OutputNames map[string][]string
	}
)

var pulumiPlanActionOrder = []string{
	"same",
	"read",
	"create",
	"import",
	"import-replacement",
	"refresh",
	"update",
	"replace",
	"create-replacement",
	"update-replacement",
	"read-replacement",
	"delete-replaced",
	"delete",
	"discard",
	"discard-replaced",
	"remove-pending-replace",
}

// state holds the inputs the stack state last gave each resource by urn, and may be
// nil when the state could not be read.
func parsePulumiPlanReport(data []byte, state map[string]map[string]any, metadata terraform.Metadata) (terraform.Report, error) {
	plan, planVersion, err := parsePulumiPlan(data)
	if err != nil {
		return terraform.Report{}, err
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

		summary := resourceSummary(string(urn), resource, state[string(urn)])

		// the steps of a replacement come in the order they run, which is how a
		// replacement that deletes first tells itself apart.
		if created, deleted := slices.Index(steps, "create-replacement"), slices.Index(steps, "delete-replaced"); created >= 0 && deleted >= 0 {
			summary.Detail = "create before delete"
			if deleted < created {
				summary.Detail = "delete before create"
			}
		}

		outputNames := resourceOutputNames(resource)

		for _, action := range steps {
			accumulator := accumulators[action]
			if accumulator == nil {
				accumulator = &actionAccumulator{
					Action:      action,
					Resources:   map[string]terraform.Resource{},
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

	return terraform.Report{
		Title: "Pulumi preview report",
		Labels: terraform.Labels{
			Target:      "Stack",
			Outputs:     "Output properties",
			ToolVersion: "Pulumi version",
		},
		Metadata: metadata,
		Actions:  buildPulumiPlanActions(accumulators),
	}, nil
}

func parsePulumiPlan(data []byte) (apitype.DeploymentPlanV1, int, error) {
	// a decode failure means the payload is an unwrapped plan rather than a version
	// envelope, so it falls through to the unversioned path instead of erroring.
	var versioned apitype.VersionedDeploymentPlan
	if err := json.Unmarshal(data, &versioned, json.RejectUnknownMembers(true)); err == nil {
		plan := bytes.TrimSpace(versioned.Plan)
		if len(plan) > 0 && !bytes.Equal(plan, []byte("null")) {
			deployment, err := parsePulumiDeploymentPlan(plan)

			return deployment, versioned.Version, err
		}
	}

	deployment, err := parsePulumiDeploymentPlan(data)

	return deployment, 0, err
}

func parsePulumiDeploymentPlan(data []byte) (apitype.DeploymentPlanV1, error) {
	var plan apitype.DeploymentPlanV1
	if err := json.Unmarshal(data, &plan, json.RejectUnknownMembers(true)); err != nil {
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

func resourceSummary(urn string, plan apitype.ResourcePlanV1, old map[string]any) terraform.Resource {
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

	return terraform.Resource{
		Name:    pulumiResourceName(resourceType, name),
		Id:      urn,
		Changes: resourceChanges(plan.Goal, old),
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

func buildPulumiPlanActions(accumulators map[string]*actionAccumulator) []terraform.Action {
	actions := make([]terraform.Action, 0, len(accumulators))

	for _, accumulator := range accumulators {
		resources := slices.Collect(maps.Values(accumulator.Resources))
		slices.SortFunc(resources, func(left, right terraform.Resource) int {
			return strings.Compare(left.Id, right.Id)
		})

		outputs := make([]terraform.Output, 0, len(accumulator.OutputNames))
		for _, urn := range slices.Sorted(maps.Keys(accumulator.OutputNames)) {
			outputs = append(outputs, terraform.Output{
				Name:   accumulator.Resources[urn].Name,
				Fields: sortedUniqueStrings(accumulator.OutputNames[urn]),
			})
		}

		actions = append(actions, terraform.Action{
			Action:    accumulator.Action,
			Resources: resources,
			Outputs:   outputs,
		})
	}

	terraform.SortActions(actions, pulumiPlanActionOrder)

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
