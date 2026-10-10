package generate

import (
	"context"
	"fmt"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
)

type (
	Pipe struct {
		Build   string
		Paths   []string
		Command string
	}

	Ctx struct {
		Whoami string
		State  string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			// the stack only lives for the run, so its secrets never need a passphrase.
			setup.C.Env["PULUMI_CONFIG_PASSPHRASE"] = ""

			// the paths are removed before every run, so they can not be allowed to take the
			// working directory or anything outside of it along.
			for _, path := range P.Paths {
				if !filepath.IsLocal(path) || filepath.Clean(path) == "." {
					return fmt.Errorf("Generated path has to be inside the working directory: %s", path)
				}
			}

			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				build(tl).Job(),
				backend(tl).Job(),
				restate(tl).Job(),
				generate(tl).Job(),
			)
		})
}
