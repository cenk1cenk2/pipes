package setup

import (
	"path/filepath"
	"slices"

	"github.com/cenk1cenk2/plumber/v6"
)

func ResolveOverlays(tl *plumber.TaskList) *plumber.Task {
	return tl.CreateTask("resolve").
		Set(func(t *plumber.Task) error {
			cwd := C.Cwd
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
