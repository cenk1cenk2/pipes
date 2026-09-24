package token

import (
	"context"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"

	"gitlab.kilic.dev/devops/pipes/github/client"
)

type (
	Token struct {
		Repositories           []string
		Permissions            map[string]string
		Variable               string
		GitCredentialsVariable string
		File                   string `validate:"required,filepath"`
	}

	Pipe struct {
		App client.Config
		Token
	}

	// the client lives in the context, since it is dialled only once the flags
	// carrying its address have been parsed.
	Ctx struct {
		Client client.ApplicationClientAdapter
		Token  string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			if !P.App.Enabled() {
				return fmt.Errorf("GitHub App id, installation id and private key are required to mint a token.")
			}

			if P.Token.Variable == "" && P.Token.GitCredentialsVariable == "" {
				return fmt.Errorf("Token variable or git credentials variable is required, otherwise nothing would be written.")
			}

			C.Client = client.NewApplicationClient(P.App.ApiUrl, p.Cli.Name)

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				mint(tl).Job(),
				write(tl).Job(),
			)
		})
}
