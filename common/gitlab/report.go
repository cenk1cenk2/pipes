package gitlab

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	clientgitlab "gitlab.com/gitlab-org/api/client-go"
)

const (
	CATEGORY_GITLAB_MERGE_REQUEST_REPORT = "Gitlab Merge Request Report"
)

type MergeRequestReportConfig struct {
	Enabled        bool
	Token          string `validate:"required_if=Enabled true"`
	ApiUrl         string `validate:"required_if=Enabled true"`
	ProjectId      string `validate:"required_if=Enabled true"`
	MergeRequestId int64  `validate:"required_if=Enabled true,omitempty,gt=0"`
	Identifier     string `validate:"required_if=Enabled true,omitempty,printascii,excludes=-->"`
}

type MergeRequestReportResult struct {
	NoteId int64
}

func NewMergeRequestReportFlags(config *MergeRequestReportConfig) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.enabled",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITLAB_MR_REPORT_ENABLED"),
			),
			Usage:       "Enable GitLab merge request report note on the given merge request.",
			Required:    false,
			Value:       true,
			Destination: &config.Enabled,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.token",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GL_PIPES_TOKEN"),
			),
			Usage:       "GitLab API token for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &config.Token,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.api-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_API_V4_URL"),
			),
			Usage:       "GitLab API URL for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &config.ApiUrl,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.project-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PROJECT_ID"),
			),
			Usage:       "GitLab project id for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &config.ProjectId,
		},

		&cli.Int64Flag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.merge-request-iid",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_MERGE_REQUEST_IID"),
			),
			Usage:       "GitLab merge request iid for merge request report notes.",
			Required:    false,
			Value:       0,
			Destination: &config.MergeRequestId,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.identifier",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITLAB_MR_REPORT_IDENTIFIER"),
				cli.EnvVar("CI_JOB_NAME"),
			),
			Usage:       "Hidden marker identifier for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &config.Identifier,
		},
	}
}

func UpsertMergeRequestReport(
	ctx context.Context,
	config MergeRequestReportConfig,
	body string,
) (*MergeRequestReportResult, error) {
	marker := fmt.Sprintf("<!-- %s%s -->", "gitlab-pipes:mr-report:", config.Identifier)
	if !strings.Contains(body, marker) {
		body = fmt.Sprintf("%s\n\n%s", strings.TrimRight(body, "\n"), marker)
	}

	client, err := clientgitlab.NewClient(
		config.Token,
		clientgitlab.WithBaseURL(config.ApiUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab client: %w", err)
	}

	var note *clientgitlab.Note
	page := int64(1)

	for {
		existing, response, err := client.Notes.ListMergeRequestNotes(
			config.ProjectId,
			config.MergeRequestId,
			&clientgitlab.ListMergeRequestNotesOptions{
				ListOptions: clientgitlab.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			},
			clientgitlab.WithContext(ctx),
		)
		if err != nil {
			return nil, fmt.Errorf("list GitLab merge request notes: %w", err)
		}

		for _, n := range existing {
			if strings.Contains(n.Body, marker) {
				note = n
				break
			}
		}

		if note != nil || response == nil || response.NextPage == 0 {
			break
		}

		page = response.NextPage
	}

	if note != nil {
		updated, _, err := client.Notes.UpdateMergeRequestNote(
			config.ProjectId,
			config.MergeRequestId,
			note.ID,
			&clientgitlab.UpdateMergeRequestNoteOptions{
				Body: clientgitlab.Ptr(body),
			},
			clientgitlab.WithContext(ctx),
		)
		if err != nil {
			return nil, fmt.Errorf("update GitLab merge request report note: %w", err)
		}

		return &MergeRequestReportResult{
			NoteId: updated.ID,
		}, nil
	}

	created, _, err := client.Notes.CreateMergeRequestNote(
		config.ProjectId,
		config.MergeRequestId,
		&clientgitlab.CreateMergeRequestNoteOptions{
			Body: clientgitlab.Ptr(body),
		},
		clientgitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab merge request report note: %w", err)
	}

	return &MergeRequestReportResult{
		NoteId: created.ID,
	}, nil
}
