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

			dirs := make([]string, 0, len(matches))
			for _, match := range matches {
				dirs = append(dirs, filepath.Dir(match))
			}

			dirs = slices.Compact(slices.Sorted(slices.Values(dirs)))

			if P.DiscoveryStrategy == DISCOVERY_STRATEGY_ROOTS {
				dirs = rootOverlays(dirs)
			}

			overlays := make([]string, 0, len(dirs))
			for _, dir := range dirs {
				overlays = append(overlays, filepath.Join(cwd, dir))
			}

			t.Log.Debugf("Discovered Kustomize overlays: %s", strings.Join(overlays, ", "))

			C.Overlays = overlays

			return nil
		})
}

// rootOverlays keeps only overlay directories that are not nested under another
// discovered overlay. ArgoCD points at the top-level overlay (e.g.
// .deploy/<cluster>) and pulls nested kustomizations in through resources:, so
// rendering each nested one standalone is redundant and can fail on relative
// paths that only resolve from the parent context.
func rootOverlays(dirs []string) []string {
	set := make(map[string]struct{}, len(dirs))
	for _, dir := range dirs {
		set[dir] = struct{}{}
	}

	roots := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		nested := false
		for parent := filepath.Dir(dir); parent != "." && parent != string(filepath.Separator); parent = filepath.Dir(parent) {
			if _, ok := set[parent]; ok {
				nested = true

				break
			}
		}

		if !nested {
			roots = append(roots, dir)
		}
	}

	return roots
}
