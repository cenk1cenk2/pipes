package gitlab

import (
	"testing"

	validator "github.com/go-playground/validator/v10"
)

func TestMergeRequestReportConfigAllowsEnabledReportOutsideMergeRequest(t *testing.T) {
	t.Parallel()

	err := validator.New().Struct(MergeRequestReportConfig{
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("expected enabled non-MR report config to validate: %v", err)
	}
}

func TestResolveReportIdentifier(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		override       string
		jobName        string
		discriminators []string
		expected       string
	}{
		"explicit override wins":        {override: "custom", jobName: "plan", discriminators: []string{"cache"}, expected: "custom"},
		"job name alone when no target": {jobName: "tf-plan", expected: "tf-plan"},
		"target disambiguates the job":  {jobName: "plan", discriminators: []string{"cache"}, expected: "plan:cache"},
		"several targets join":          {jobName: "plan", discriminators: []string{"cache", ".deploy/sun"}, expected: "plan:cache:.deploy/sun"},
		"empty parts drop out":          {jobName: "", discriminators: []string{"", "cache"}, expected: "cache"},
		"marker terminator is stripped": {jobName: "plan", discriminators: []string{"ca-->che"}, expected: "plan:cache"},
		"non ascii is stripped":         {jobName: "plan", discriminators: []string{"caché"}, expected: "plan:cach"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := ResolveReportIdentifier(test.override, test.jobName, test.discriminators...)
			if got != test.expected {
				t.Fatalf("expected identifier %q, got %q", test.expected, got)
			}
		})
	}
}

func TestMergeRequestReportConfigRejectsInvalidMergeRequestId(t *testing.T) {
	t.Parallel()

	err := validator.New().Struct(MergeRequestReportConfig{
		MergeRequestId: -1,
	})
	if err == nil {
		t.Fatal("expected negative merge request id to fail validation")
	}
}

func TestMergeRequestReportConfigRequiresGitlabContextWithMergeRequestId(t *testing.T) {
	t.Parallel()

	err := validator.New().Struct(MergeRequestReportConfig{
		Enabled:        true,
		MergeRequestId: 1,
	})
	if err == nil {
		t.Fatal("expected missing GitLab context with merge request id to fail validation")
	}
}
