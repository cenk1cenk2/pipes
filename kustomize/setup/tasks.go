package setup

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	. "github.com/cenk1cenk2/plumber/v6"
	"sigs.k8s.io/kustomize/api/provenance"
)

func KustomizeVersion(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.Log.Infof("Kustomize library version: %s", provenance.GetProvenance().Version)

			return nil
		})
}

func DiscoverOverlays(tl *TaskList) *Task {
	return tl.CreateTask("discover").
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

				t.Log.Debugf("Using explicit overlay paths: %s", strings.Join(C.Overlays, ", "))

				return nil
			}

			fs := os.DirFS(cwd)

			matches := []string{}

			for _, pattern := range P.DiscoveryPattern {
				match, err := glob.Glob(fs, pattern)

				if err != nil {
					return err
				}

				matches = append(matches, match...)
			}

			if len(matches) == 0 {
				t.Log.Warnf(
					"Can not discover any Kustomize overlays with the given patterns: %s",
					strings.Join(P.DiscoveryPattern, ", "),
				)

				return nil
			}

			overlays := []string{}
			for _, match := range matches {
				overlays = append(overlays, filepath.Join(cwd, filepath.Dir(match)))
			}

			overlays = slices.Compact(slices.Sorted(slices.Values(overlays)))

			t.Log.Debugf("Discovered Kustomize overlays: %s", strings.Join(overlays, ", "))

			C.Overlays = overlays

			return nil
		})
}
