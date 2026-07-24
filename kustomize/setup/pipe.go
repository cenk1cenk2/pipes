package setup

import (
	"path/filepath"
	"slices"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Cwd   string `validate:"omitempty,dirpath"`
		Paths []string
	}

	Ctx struct {
		Overlays []string
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

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				ResolveOverlays(tl).Job(),
			)
		})
}

func ResolveOverlays(tl *TaskList) *Task {
	return tl.CreateTask("resolve").
		Set(func(t *Task) error {
			cwd := P.Cwd
			if cwd == "" {
				cwd = "."
			}

			if len(P.Paths) > 0 {
				overlays := make([]string, 0, len(P.Paths))
				for _, path := range P.Paths {
					overlays = append(overlays, filepath.Join(cwd, path))
				}

				C.Overlays = slices.Compact(slices.Sorted(slices.Values(overlays)))
				t.Log.Debugf("Using explicit overlay paths: %v", C.Overlays)

				return nil
			}

			C.Overlays = []string{cwd}
			t.Log.Debugf("Using overlay path: %s", cwd)

			return nil
		})
}
