package update

import (
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
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
		Token  string
		Readme map[string]ParsedReadme
		Hub    hub.Client
	}

	// Deps dials the service only once the flags carrying its address have been
	// parsed, so the pipe carries the way to reach one rather than a connection.
	Deps struct {
		Hub hub.ClientFactory
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := icli.Validated(p, P); err != nil {
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
			C.Hub = deps.Hub(P.DockerHub.Address, p.Cli.Name)

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

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
