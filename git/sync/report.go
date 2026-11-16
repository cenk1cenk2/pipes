package sync

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/patch"
)

var (
	//go:embed assets/header.md.gotmpl
	headerTemplate string
	//go:embed assets/report.md.gotmpl
	reportTemplate string
	//go:embed assets/close.md.gotmpl
	closeTemplate string
	//go:embed assets/merged.md.gotmpl
	mergedTemplate string
)

// GitLab refuses a note or a description past a million characters, and the rest of
// the report has to fit next to the diff; the specs lower it to cut realistic patches.
var reportDiffLimit = 900_000

const (
	// GitLab lists at most this many merge requests in CI_OPEN_MERGE_REQUESTS, so a full
	// list may be missing some.
	openMergeRequestsLimit = 4
)

type reportTarget struct {
	ProjectId       string
	MergeRequestIid int64
}

type (
	// What the report and the close templates render: the pipe, with the changes split
	// into the sections the limit lets through.
	reportView struct {
		*Pipe
		Changed   bool
		Split     bool
		Files     patch.Section
		Pending   patch.Section
		Generated patch.Section
		Truncated bool
		// the merge request a publish opened or updated, which only the note on the
		// merged merge requests names.
		MergeRequest *gitlab.MergeRequestResult
	}
)

func render() (string, error) {
	view := reportView{Pipe: P, Changed: C.Diff != ""}

	files := patch.Split(C.Diff)

	// git quotes an unusual name differently in the header than in the list of changed
	// paths, so such a file can not be matched and stays with what the branch lacks.
	pending := []string{}
	generated := []string{}
	for _, file := range files {
		if C.Pending == nil || strings.HasPrefix(file, "diff --git \"") || C.Pending[patch.Filename(file)] || C.Pending[patch.Previous(file)] {
			pending = append(pending, file)
		} else {
			generated = append(generated, file)
		}
	}

	if len(generated) == 0 {
		view.Files, _, view.Truncated = patch.Sections(files, reportDiffLimit)
	} else {
		// what the branch still lacks is what the next sync brings, so it leads and gets
		// the limit first, while what the branch already carries only gets the rest.
		view.Split = true

		var limit int
		var cut bool
		view.Pending, limit, view.Truncated = patch.Sections(pending, reportDiffLimit)
		view.Generated, _, cut = patch.Sections(generated, limit)
		view.Truncated = view.Truncated || cut
	}

	return template(reportTemplate, view)
}

// The comment left on the merge request a sync closes once the target already matches.
func renderClose() (string, error) {
	return template(closeTemplate, reportView{Pipe: P})
}

// The note left on the merge requests merged into the commit a publish ran for.
func renderMerged() (string, error) {
	return template(mergedTemplate, reportView{Pipe: P, MergeRequest: C.MergeRequest})
}

func template(body string, view reportView) (string, error) {
	return InlineTemplate(headerTemplate+body, view, map[string]any{"link": link})
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
