package pipe

import (
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	DockerHub struct {
		Username string
		Password string
		Address  string
	}

	Readme struct {
		Repository  string
		File        string
		Description string
		Matrix      []ReadmeMatrixJson
	}

	Pipe struct {
		DockerHub
		Readme
	}

	Ctx struct {
		Token       string
		Readme      map[string]ParsedReadme
		ReadmeFiles map[string][]byte
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			if len(P.Readme.Description) > 100 {
				return fmt.Errorf(
					"Readme short description can only be 100 characters long while you have: %d",
					len(P.Readme.Description),
				)
			}

			if P.Readme.Repository == "" && len(P.Readme.Matrix) == 0 {
				return fmt.Errorf("You have to either provide a target via Repository or multiple targets through the Matrix.")
			}

			C.Readme = make(map[string]ParsedReadme)

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				LoginToDockerHubRegistry(tl).Job(),
				DiscoverJobs(tl).Job(),
				UpdateDockerReadme(tl).Job(),
			)
		})
}
