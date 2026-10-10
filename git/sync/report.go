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

// GitLab refuses a note or a description past a million characters, and the rest of
// the report has to fit next to the diff; the specs lower it to cut realistic patches.
var reportDiffLimit = 900_000

const (
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

	files := []string{}
	for diff := C.Diff; diff != ""; {
		end := len(diff)
		if i := strings.Index(diff, "\ndiff --git "); i >= 0 {
			end = i + 1
		}

		files = append(files, diff[:end])
		diff = diff[end:]
	}

	fmt.Fprintf(body, "%d file(s) changed.\n", len(files))

	// the name and the fences of every file count against the limit as well, which a
	// diff of many small files would otherwise outgrow.
	limit := reportDiffLimit
	shown := 0
	cut := false
	for _, file := range files {
		name := filename(file)
		diff := hunks(file)
		marker := fence(diff)
		overhead := len(name) + 2*len(marker) + len("\n``\n\ndiff\n\n")
		if overhead >= limit {
			break
		}

		if len(diff) > limit-overhead {
			end := strings.LastIndex(diff[:limit-overhead], "\n") + 1
			if end == 0 {
				break
			}

			diff = diff[:end]
			cut = true
		}

		fmt.Fprintf(body, "\n`%s`\n\n%sdiff\n%s%s\n", name, marker, diff, marker)
		limit -= overhead + len(diff)
		shown++

		if cut {
			break
		}
	}

	if more := len(files) - shown; more > 0 {
		fmt.Fprintf(body, "\n%d more file(s) are not shown.\n", more)
	}

	if cut || shown < len(files) {
		artifacts := "the artifacts of the job"
		if P.ReportMetadata.JobUrl != "" {
			artifacts = link(artifacts, P.ReportMetadata.JobUrl+"/artifacts/browse")
		}

		fmt.Fprintf(body, "\nThe diff is truncated, the full patch `%s` is in %s.\n", P.Patch, artifacts)
	}

	return body.String()
}

// The name a file carries after the change, which is the one a rename ends up with.
func filename(diff string) string {
	header, _, _ := strings.Cut(diff, "\n")
	header = strings.TrimPrefix(header, "diff --git ")
	if i := strings.LastIndex(header, " b/"); i >= 0 {
		return header[i+len(" b/"):]
	}

	return header
}

// The name above the section already says which file it is, so only what git marks
// beyond that stays: the hunks, and the mode, rename and binary lines of the header.
func hunks(diff string) string {
	body := &strings.Builder{}
	header := true
	for i, line := range strings.SplitAfter(diff, "\n") {
		if i == 0 {
			continue
		}

		if strings.HasPrefix(line, "@@") || line == "GIT binary patch\n" {
			header = false
		}

		if header && (strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ")) {
			continue
		}

		body.WriteString(line)
	}

	return body.String()
}

// A diff of a markdown file can carry a fence of its own, which would close this one.
func fence(diff string) string {
	fence := "```"
	for _, run := range fenceRun.FindAllString(diff, -1) {
		if len(run) >= len(fence) {
			fence = strings.Repeat("`", len(run)+1)
		}
	}

	return fence
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
