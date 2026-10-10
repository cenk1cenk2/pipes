package sync

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

const (
	// GitLab refuses a note or a description past a million characters, and the rest of
	// the report has to fit next to the diff.
	reportDiffLimit = 900_000
	// GitLab lists at most this many merge requests in CI_OPEN_MERGE_REQUESTS, so a full
	// list may be missing some.
	openMergeRequestsLimit = 4
)

var fenceRun = regexp.MustCompile("`{3,}")

type reportTarget struct {
	ProjectId       string
	MergeRequestIid int64
}

func render() string {
	body := &strings.Builder{}

	fmt.Fprintf(body, "## Sync into `%s@%s`\n\n", P.Target.Project, P.Target.Branch)

	metadata := []string{}
	if m := P.ReportMetadata; m.JobName != "" {
		metadata = append(metadata, fmt.Sprintf("Job: %s", link(m.JobName, m.JobUrl)))
	}

	if m := P.ReportMetadata; m.PipelineId != "" {
		metadata = append(metadata, fmt.Sprintf("Pipeline: %s", link("#"+m.PipelineId, m.PipelineUrl)))
	}

	if m := P.ReportMetadata; m.CommitShortSha != "" {
		metadata = append(metadata, fmt.Sprintf("Commit: `%s`", m.CommitShortSha))
	}

	if len(metadata) > 0 {
		fmt.Fprintf(body, "%s\n\n", strings.Join(metadata, " | "))
	}

	if C.Diff == "" {
		body.WriteString("No changes.\n")

		return body.String()
	}

	fmt.Fprintf(body, "%d file(s) changed.\n\n", strings.Count("\n"+C.Diff, "\ndiff --git "))

	diff := C.Diff
	truncated := len(diff) > reportDiffLimit
	if truncated {
		cut := strings.LastIndex(diff[:reportDiffLimit], "\n") + 1
		if cut == 0 {
			cut = reportDiffLimit
		}

		diff = diff[:cut]
	}

	// a diff of a markdown file can carry a fence of its own, which would close this one.
	fence := "```"
	for _, run := range fenceRun.FindAllString(diff, -1) {
		if len(run) >= len(fence) {
			fence = strings.Repeat("`", len(run)+1)
		}
	}

	fmt.Fprintf(body, "%sdiff\n%s%s\n", fence, diff, fence)

	if truncated {
		artifacts := "the artifacts of the job"
		if P.ReportMetadata.JobUrl != "" {
			artifacts = link(artifacts, P.ReportMetadata.JobUrl+"/artifacts/browse")
		}

		fmt.Fprintf(body, "\nThe diff is truncated, the full patch `%s` is in %s.\n", P.Patch, artifacts)
	}

	return body.String()
}

func link(text string, url string) string {
	if url == "" {
		return text
	}

	return fmt.Sprintf("[%s](%s)", text, url)
}

// A merge request pipeline names its merge request, while a branch pipeline only gets
// the ones open from its branch, which have to be asked for once the list may be cut.
func reportTargets(ctx context.Context) ([]reportTarget, error) {
	if P.MergeRequestReport.MergeRequestIid != 0 {
		return []reportTarget{{ProjectId: P.MergeRequestReport.ProjectId, MergeRequestIid: P.MergeRequestReport.MergeRequestIid}}, nil
	}

	targets := []reportTarget{}

	for entry := range strings.SplitSeq(P.Project.OpenMergeRequests, ",") {
		if entry = strings.TrimSpace(entry); entry == "" {
			continue
		}

		project, iid, ok := strings.Cut(entry, "!")
		if !ok {
			return nil, fmt.Errorf("Can not parse the open merge request: %s", entry)
		}

		parsed, err := strconv.ParseInt(iid, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Can not parse the open merge request: %s -> %w", entry, err)
		}

		targets = append(targets, reportTarget{ProjectId: project, MergeRequestIid: parsed})
	}

	if len(targets) < openMergeRequestsLimit {
		return targets, nil
	}

	mergeRequests, err := newMergeRequests(gitlab.MergeRequestConfig{
		Token:  P.MergeRequestReport.Token,
		ApiUrl: P.MergeRequestReport.ApiUrl,
	})
	if err != nil {
		return nil, err
	}

	listed, _, err := mergeRequests.ListProjectMergeRequests(
		P.MergeRequestReport.ProjectId,
		&clientgitlab.ListProjectMergeRequestsOptions{
			ListOptions:  clientgitlab.ListOptions{PerPage: 100},
			State:        new("opened"),
			SourceBranch: new(P.Refs.Branch),
		},
		clientgitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("list GitLab merge requests: %w", err)
	}

	targets = []reportTarget{}
	for _, mergeRequest := range listed {
		targets = append(targets, reportTarget{ProjectId: P.MergeRequestReport.ProjectId, MergeRequestIid: mergeRequest.IID})
	}

	return targets, nil
}
