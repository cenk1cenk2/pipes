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
