package preview

import (
	"maps"
	"slices"

	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/sig"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

const (
	secretValue  = terraform.Marker("[secret]")
	unknownValue = terraform.Marker("[unknown]")

	// what Pulumi writes into a plan where a property is only computed on apply.
	computedValuePlaceholder = "04da6b54-80e4-46f7-96ec-b56ff0331ba9"
)

// The plan carries the inputs a resource is going to be given and no more: a
// deleted property comes as its name alone and an updated one as its new value
// without the old.
func resourceChanges(goal *apitype.GoalV1) []terraform.Change {
	if goal == nil {
		return nil
	}

	diff := goal.InputDiff
	names := sortedUniqueStrings(slices.Concat(
		slices.Collect(maps.Keys(diff.Adds)),
		diff.Deletes,
		slices.Collect(maps.Keys(diff.Updates)),
	))

	changes := make([]terraform.Change, 0, len(names))
	for _, name := range names {
		if value, ok := diff.Adds[name]; ok {
			changes = append(changes, terraform.Expand(terraform.ChangeCreate, name, masked(value)))

			continue
		}

		if value, ok := diff.Updates[name]; ok {
			change := terraform.Expand("", name, masked(value))
			change.Action = terraform.ChangeUpdate
			changes = append(changes, change)

			continue
		}

		changes = append(changes, terraform.Change{Name: name, Action: terraform.ChangeDelete})
	}

	return changes
}

// masked is the value with every secret replaced by its marker, whichever way the
// plan was saved: a secret is an object carrying the secret signature, holding either
// the ciphertext or, when the preview was run with --show-secrets, the plaintext.
func masked(value any) any {
	switch value := value.(type) {
	case string:
		if value == computedValuePlaceholder {
			return unknownValue
		}
	case map[string]any:
		switch value[sig.Key] {
		case sig.Secret:
			return secretValue
		case sig.ResourceReference:
			return value["urn"]
		case sig.AssetSig:
			return terraform.Marker("[asset]")
		case sig.ArchiveSig:
			return terraform.Marker("[archive]")
		}

		out := make(map[string]any, len(value))
		for key, element := range value {
			out[key] = masked(element)
		}

		return out
	case []any:
		out := make([]any, len(value))
		for index, element := range value {
			out[index] = masked(element)
		}

		return out
	}

	return value
}
