package iac

import (
	"encoding/json"
	"slices"
)

// The counts a plan summary file carries, in the shape both pipes have always
// written them.
type Summary struct {
	Create int `json:"create"`
	Update int `json:"update"`
	Delete int `json:"delete"`
}

// A replacement costs a create and a delete, but a tool that also reports the two
// halves as their own actions would then have them counted twice. Those halves win
// whenever the plan carries them.
func Summarize(report Report) Summary {
	summary := Summary{}
	hasReplacementDetails := slices.ContainsFunc(report.Actions, func(action Action) bool {
		return action.Action == "create-replacement" || action.Action == "delete-replaced"
	})

	for _, action := range report.Actions {
		count := len(action.Resources)

		switch action.Action {
		case "create", "create-replacement":
			summary.Create += count
		case "update", "update-replacement":
			summary.Update += count
		case "delete", "delete-replaced":
			summary.Delete += count
		case "replace":
			if !hasReplacementDetails {
				summary.Create += count
				summary.Delete += count
			}
		}
	}

	return summary
}

func RenderSummary(summary Summary) string {
	// Three ints behind a fixed set of tags leave the encoder nothing to fail on.
	body, _ := json.MarshalIndent(summary, "", "  ")

	return string(body) + "\n"
}
