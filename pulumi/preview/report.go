package preview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
)

//go:embed assets/mr-report.md.gotmpl
var mrReportTemplate string

type (
	MergeRequestReportMetadata struct {
		Stack          string
		JobName        string
		JobUrl         string
		PipelineId     string
		PipelineUrl    string
		CommitSha      string
		CommitShortSha string
	}

	PulumiPlanReport struct {
		Metadata     MergeRequestReportMetadata
		PlanVersion  int
		Manifest     PulumiPlanManifest
		Actions      []PulumiPlanAction
		TotalActions int
	}

	PulumiPlanManifest struct {
		Time    string
		Version string
	}

	PulumiPlanAction struct {
		Action    string
		Count     int
		Resources []PulumiPlanResource
		Outputs   []PulumiPlanOutput
	}

	PulumiPlanResource struct {
		Urn  string
		Type string
		Name string
	}

	PulumiPlanOutput struct {
		Resource PulumiPlanResource
		Names    []string
	}

	actionAccumulator struct {
		Action       string
		Count        int
		Resources    map[string]PulumiPlanResource
		OutputNames  map[string]map[string]struct{}
		OutputLookup map[string]PulumiPlanResource
	}
)

func (r PulumiPlanReport) HasMetadata() bool {
	return r.Metadata.Stack != "" ||
		r.Metadata.JobName != "" ||
		r.Metadata.JobUrl != "" ||
		r.Metadata.PipelineId != "" ||
		r.Metadata.PipelineUrl != "" ||
		r.Metadata.CommitSha != "" ||
		r.Metadata.CommitShortSha != "" ||
		r.PlanVersion != 0 ||
		r.Manifest.Version != "" ||
		r.Manifest.Time != ""
}

func parsePulumiPlanReport(data []byte, metadata MergeRequestReportMetadata) (*PulumiPlanReport, error) {
	plan, planVersion, err := parsePulumiPlan(data)
	if err != nil {
		return nil, err
	}

	accumulators := map[string]*actionAccumulator{}

	for urn, resource := range plan.ResourcePlans {
		steps := resourceActions(resource)
		if len(steps) == 0 {
			continue
		}

		summary := resourceSummary(string(urn), resource)
		outputNames := resourceOutputNames(resource)

		for _, action := range steps {
			accumulator := getActionAccumulator(accumulators, action)
			accumulator.Count++
			accumulator.Resources[summary.Urn] = summary

			if len(outputNames) > 0 {
				accumulator.OutputLookup[summary.Urn] = summary

				for _, outputName := range outputNames {
					if accumulator.OutputNames[summary.Urn] == nil {
						accumulator.OutputNames[summary.Urn] = map[string]struct{}{}
					}

					accumulator.OutputNames[summary.Urn][outputName] = struct{}{}
				}
			}
		}
	}

	report := &PulumiPlanReport{
		Metadata:    metadata,
		PlanVersion: planVersion,
		Manifest: PulumiPlanManifest{
			Time:    formatManifestTime(plan.Manifest.Time),
			Version: plan.Manifest.Version,
		},
		Actions: buildPulumiPlanActions(accumulators),
	}

	for _, action := range report.Actions {
		report.TotalActions += action.Count
	}

	return report, nil
}

func renderMergeRequestReport(report *PulumiPlanReport) (string, error) {
	tmpl, err := template.New("mr-report.md.gotmpl").
		Parse(mrReportTemplate)
	if err != nil {
		return "", fmt.Errorf("parse Pulumi merge request report template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, report); err != nil {
		return "", fmt.Errorf("render Pulumi merge request report: %w", err)
	}

	return strings.TrimRight(buf.String(), "\n") + "\n", nil
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
	seen := map[string]struct{}{}
	actions := make([]string, 0, len(resource.Steps))

	for _, action := range resource.Steps {
		action := strings.TrimSpace(string(action))
		if action == "" {
			continue
		}

		if _, ok := seen[action]; ok {
			continue
		}

		seen[action] = struct{}{}
		actions = append(actions, action)
	}

	return actions
}

func resourceSummary(urn string, plan apitype.ResourcePlanV1) PulumiPlanResource {
	summary := PulumiPlanResource{
		Urn: urn,
	}

	if plan.Goal != nil {
		summary.Type = string(plan.Goal.Type)
		summary.Name = plan.Goal.Name
	}

	if summary.Type == "" || summary.Name == "" {
		urnType, urnName := splitPulumiUrn(summary.Urn)
		if summary.Type == "" {
			summary.Type = urnType
		}

		if summary.Name == "" {
			summary.Name = urnName
		}
	}

	return summary
}

func resourceOutputNames(resource apitype.ResourcePlanV1) []string {
	if resource.Goal == nil {
		return nil
	}

	diff := resource.Goal.OutputDiff
	names := make([]string, 0, len(diff.Adds)+len(diff.Deletes)+len(diff.Updates))
	names = append(names, mapKeys(diff.Adds)...)
	names = append(names, diff.Deletes...)
	names = append(names, mapKeys(diff.Updates)...)

	return uniqueSortedStrings(names)
}

func getActionAccumulator(accumulators map[string]*actionAccumulator, action string) *actionAccumulator {
	accumulator, ok := accumulators[action]
	if ok {
		return accumulator
	}

	accumulator = &actionAccumulator{
		Action:       action,
		Resources:    map[string]PulumiPlanResource{},
		OutputNames:  map[string]map[string]struct{}{},
		OutputLookup: map[string]PulumiPlanResource{},
	}
	accumulators[action] = accumulator

	return accumulator
}

func buildPulumiPlanActions(accumulators map[string]*actionAccumulator) []PulumiPlanAction {
	actions := make([]PulumiPlanAction, 0, len(accumulators))

	for _, accumulator := range accumulators {
		action := PulumiPlanAction{
			Action:    accumulator.Action,
			Count:     accumulator.Count,
			Resources: sortedResources(accumulator.Resources),
			Outputs:   sortedOutputs(accumulator.OutputNames, accumulator.OutputLookup),
		}

		actions = append(actions, action)
	}

	sort.Slice(actions, func(i, j int) bool {
		left := actionSortRank(actions[i].Action)
		right := actionSortRank(actions[j].Action)
		if left == right {
			return actions[i].Action < actions[j].Action
		}

		return left < right
	})

	return actions
}

func sortedResources(resources map[string]PulumiPlanResource) []PulumiPlanResource {
	keys := make([]string, 0, len(resources))
	for key := range resources {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	sorted := make([]PulumiPlanResource, 0, len(keys))
	for _, key := range keys {
		sorted = append(sorted, resources[key])
	}

	return sorted
}

func sortedOutputs(
	outputNames map[string]map[string]struct{},
	outputLookup map[string]PulumiPlanResource,
) []PulumiPlanOutput {
	keys := make([]string, 0, len(outputNames))
	for key := range outputNames {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	outputs := make([]PulumiPlanOutput, 0, len(keys))
	for _, key := range keys {
		outputs = append(outputs, PulumiPlanOutput{
			Resource: outputLookup[key],
			Names:    sortedStringSet(outputNames[key]),
		})
	}

	return outputs
}

func sortedStringSet(values map[string]struct{}) []string {
	sorted := make([]string, 0, len(values))
	for value := range values {
		sorted = append(sorted, value)
	}

	sort.Strings(sorted)

	return sorted
}

func mapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}

	return keys
}

func uniqueSortedStrings(values []string) []string {
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			seen[value] = struct{}{}
		}
	}

	return sortedStringSet(seen)
}

func splitPulumiUrn(urn string) (string, string) {
	parts := strings.Split(urn, "::")
	if len(parts) < 2 {
		return "", ""
	}

	return parts[len(parts)-2], parts[len(parts)-1]
}

func actionSortRank(action string) int {
	ranks := map[string]int{
		"same":                   10,
		"read":                   20,
		"create":                 30,
		"import":                 40,
		"update":                 50,
		"replace":                60,
		"create-replacement":     70,
		"update-replacement":     80,
		"read-replacement":       90,
		"delete-replaced":        100,
		"delete":                 110,
		"discard":                120,
		"remove-pending-replace": 130,
	}

	if rank, ok := ranks[action]; ok {
		return rank
	}

	return 1000
}
