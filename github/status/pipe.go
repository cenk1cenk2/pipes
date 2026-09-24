package status

import (
	"context"
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"

	"gitlab.kilic.dev/devops/pipes/github/client"
)

type (
	Status struct {
		Token       string
		Project     string `validate:"required"`
		State       string `validate:"oneof=pending success failure error"`
		Sha         string `validate:"required"`
		TargetUrl   string `validate:"omitempty,url"`
		Context     string `validate:"required"`
		Description string
	}

	Pipe struct {
		App client.Config
		Status
	}

	// the client lives in the context, since it is dialled only once the flags
	// carrying its address have been parsed.
	Ctx struct {
		Client     client.ApplicationClientAdapter
		Repository string
		Token      string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			// an empty variable still sets its flag, which would hide the value the
			// runner exports under the name the flag falls back to.
			for field, fallback := range map[*string]string{
				&P.Status.Token:     "GH_TOKEN",
				&P.Status.State:     "GITHUB_STATUS_REPORT",
				&P.Status.Sha:       "CI_COMMIT_SHA",
				&P.Status.TargetUrl: "CI_PIPELINE_URL",
			} {
				if *field == "" {
					*field = os.Getenv(fallback)
				}
			}

			p.AppendSecrets(P.Status.Token)

			if err := p.Validate(P); err != nil {
				return err
			}

			if P.Status.Token == "" && !P.App.Enabled() {
				return fmt.Errorf("Either a GitHub token or the GitHub App credentials are required to post a status.")
			}

			owner, repository, ok := strings.Cut(P.Status.Project, "/")

			if !ok || owner == "" || repository == "" || strings.Contains(repository, "/") {
				return fmt.Errorf("Project has to be in the form of owner/repository: %s", P.Status.Project)
			}

			C.Repository = repository
			C.Client = client.NewApplicationClient(P.App.ApiUrl, p.Cli.Name)

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				authenticate(tl).Job(),
				post(tl).Job(),
			)
		})
}
