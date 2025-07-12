package login

import (
	"encoding/json"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Npm struct {
		Login     []NpmLoginJson
		NpmRcFile []string
		NpmRc     string
	}

	Pipe struct {
		Npm
	}
)

var TL = TaskList{}

var P = &Pipe{}
var raw = &struct {
	NpmLogin string
}{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if v := raw.NpmLogin; v != "" {
				if err := json.Unmarshal([]byte(v), &P.Npm.Login); err != nil {
					return fmt.Errorf("Can not unmarshal Npm registry login credentials: %w", err)
				}
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GenerateNpmRc(tl).Job(),
				VerifyNpmLogin(tl).Job(),
			)
		})
}
