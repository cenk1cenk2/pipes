package preview

import (
	"maps"
	"slices"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/sig"
	"gitlab.kilic.dev/devops/pipes/internal/markdown"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

const (
	secretValue  = terraform.Marker("[secret]")
	unknownValue = terraform.Marker("[unknown]")

	// what Pulumi writes into a plan where a property is only computed on apply.
	computedValuePlaceholder = "04da6b54-80e4-46f7-96ec-b56ff0331ba9"

	assetValue   = terraform.Marker("[asset]")
	archiveValue = terraform.Marker("[archive]")

	// what the inputs of a resource in the stack state carry for the provider's own
	// bookkeeping, never as a property of the resource.
	internalProperty = "__"
)

// The markers are the pipe's own, so the renderer only learns to dim what this pipe
// actually writes.
func init() {
	markdown.RegisterDiffLexer(
		terraform.NoAttributeChanges,
		string(secretValue),
		string(unknownValue),
		string(assetValue),
		string(archiveValue),
	)
}

// The plan carries the inputs a resource is going to be given and no more: a
// deleted property comes as its name alone and an updated one as its new value
// without the old. Old holds the inputs the stack state last gave the resource, which
// is what an update is compared against and what a deletion shows; without it a
// change shows only what the plan carries.
func resourceChanges(goal *apitype.GoalV1, old map[string]any) []terraform.Change {
	var changes []terraform.Change

	// a resource that only goes away carries no goal, and what it goes away with is
	// everything the state last gave it.
	if goal == nil {
		for _, name := range slices.Sorted(maps.Keys(old)) {
			if strings.HasPrefix(name, internalProperty) {
				continue
			}

			changes = append(changes, terraform.Expand(terraform.ChangeDelete, name, masked(old[name])))
		}

		return changes
	}

	diff := goal.InputDiff
	names := sortedUniqueStrings(slices.Concat(
		slices.Collect(maps.Keys(diff.Adds)),
		diff.Deletes,
		slices.Collect(maps.Keys(diff.Updates)),
	))

	for _, name := range names {
		before, known := old[name]

		if value, ok := diff.Adds[name]; ok {
			changes = append(changes, terraform.Expand(terraform.ChangeCreate, name, masked(value)))

			continue
		}

		if value, ok := diff.Updates[name]; ok {
			// two values that read the same once masked, a secret that changed, still
			// went through an update the plan names.
			if known {
				if change, changed := terraform.Compare(name, masked(before), masked(value)); changed {
					changes = append(changes, change)

					continue
				}
			}

			changes = append(changes, terraform.Expand(terraform.ChangeUpdate, name, masked(value)))

			continue
		}

		if known {
			changes = append(changes, terraform.Expand(terraform.ChangeDelete, name, masked(before)))

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
			return assetValue
		case sig.ArchiveSig:
			return archiveValue
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
