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

	// git quotes an unusual name differently in the header than in the list of changed
	// paths, so such a file can not be matched and stays with what the branch lacks.
	pending := []string{}
	generated := []string{}
	for _, file := range files {
		if C.Pending == nil || strings.HasPrefix(file, "diff --git \"") || C.Pending[filename(file)] || C.Pending[previous(file)] {
			pending = append(pending, file)
		} else {
			generated = append(generated, file)
		}
	}

	truncated := false
	if len(generated) == 0 {
		_, truncated = sections(body, files, reportDiffLimit)
	} else {
		branch := "the branch"
		if P.Refs.Branch != "" {
			branch = fmt.Sprintf("`%s`", P.Refs.Branch)
		}

		// what the branch still lacks is what the next sync brings, so it leads and gets
		// the limit first, while what the branch already carries only gets the rest.
		limit := reportDiffLimit
		fmt.Fprintf(body, "### Not yet generated on %s\n\n", branch)
		if len(pending) == 0 {
			fmt.Fprintf(body, "No changes beyond what %s already carries.\n", branch)
		} else {
			limit, truncated = sections(body, pending, limit)
		}

		fmt.Fprintf(body, "\n### Already generated on %s\n\n", branch)
		_, cut := sections(body, generated, limit)
		truncated = truncated || cut
	}

	if truncated {
		artifacts := "the artifacts of the job"
		if P.ReportMetadata.JobUrl != "" {
			artifacts = link(artifacts, P.ReportMetadata.JobUrl+"/artifacts/browse")
		}

		fmt.Fprintf(body, "\nThe diff is truncated, the full patch `%s` is in %s.\n", P.Patch, artifacts)
	}

	return body.String()
}

// Renders the files as one section each within the limit, and returns what is left of
// it together with whether any file did not fit. A file cut at a line ends the diff,
// so nothing after it can make it look complete.
func sections(body *strings.Builder, files []string, limit int) (int, bool) {
	fmt.Fprintf(body, "%d file(s) changed.\n", len(files))

	// the name and the fences of every file count against the limit as well, which a
	// diff of many small files would otherwise outgrow.
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
			limit = 0

			break
		}
	}

	if more := len(files) - shown; more > 0 {
		fmt.Fprintf(body, "\n%d more file(s) are not shown.\n", more)
	}

	return limit, cut || shown < len(files)
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

// The name a file carries before the change, which a rename leaves behind.
func previous(diff string) string {
	header, _, _ := strings.Cut(diff, "\n")
	header = strings.TrimPrefix(header, "diff --git a/")
	if i := strings.LastIndex(header, " b/"); i >= 0 {
		return header[:i]
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
