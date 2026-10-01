package preview

import (
	"bytes"
	"context"
	json "encoding/json/v2"
	"fmt"
	"os/exec"
)

// The stack state as `pulumi stack export` writes it, as far as the report reads it:
// the inputs each resource was last given. A preview plan only carries the values a
// resource is going to be given, so these are what its diff compares them against.
type stackState struct {
	Deployment struct {
		Resources []struct {
			URN    string         `json:"urn"`
			Inputs map[string]any `json:"inputs"`
			// a resource that is pending deletion after a replacement shares its urn
			// with the one that replaced it.
			Delete bool `json:"delete"`
		} `json:"resources"`
	} `json:"deployment"`
}

// exportStackState reads the state of the stack the preview ran against, keeping it
// in memory: it holds every value the stack carries, so it never reaches the disk
// where an artifact could pick it up.
func exportStackState(ctx context.Context, cwd string, stack string) ([]byte, error) {
	args := []string{"stack", "export"}
	if stack != "" {
		args = append(args, "--stack", stack)
	}

	cmd := exec.CommandContext(ctx, "pulumi", args...)
	cmd.Dir = cwd

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("export Pulumi stack state: %w: %s", err, bytes.TrimSpace(stderr.Bytes()))
	}

	return out, nil
}

func parseStackState(data []byte) (map[string]map[string]any, error) {
	var state stackState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parse Pulumi stack state: %w", err)
	}

	inputs := make(map[string]map[string]any, len(state.Deployment.Resources))
	for _, resource := range state.Deployment.Resources {
		if resource.Delete {
			continue
		}

		inputs[resource.URN] = resource.Inputs
	}

	return inputs, nil
}
