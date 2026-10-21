package install

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// layers is the depth of <host>/<namespace>/<type>/<version> in the plugin cache and in
// .terraform/providers, which init fills with links into the cache.
const layers = 4

// pruneCache removes the provider versions of the cache that the providers directory
// of the root does not list, and returns them with the bytes freed. A root without any
// selected provider prunes nothing, since an init that did not run must not empty the cache.
func pruneCache(cache string, providers string) ([]string, int64, error) {
	selected, err := listVersions(providers)
	if err != nil {
		return nil, 0, err
	}

	if len(selected) == 0 {
		return nil, 0, nil
	}

	cached, err := listVersions(cache)
	if err != nil {
		return nil, 0, err
	}

	var removed []string
	var freed int64

	for version := range cached {
		if _, ok := selected[version]; ok {
			continue
		}

		dir := filepath.Join(cache, version)

		freed += size(dir)

		if err := os.RemoveAll(dir); err != nil {
			return removed, freed, fmt.Errorf("Cannot remove the provider version: %s -> %w", dir, err)
		}

		removed = append(removed, version)

		// os.Remove refuses a directory with content, so this stops at the first parent in use.
		for parent := filepath.Dir(dir); parent != cache; parent = filepath.Dir(parent) {
			if os.Remove(parent) != nil {
				break
			}
		}
	}

	return removed, freed, nil
}

// listVersions returns the <host>/<namespace>/<type>/<version> entries below root as a
// set keyed by that relative path. A missing root is an empty set.
func listVersions(root string) (map[string]struct{}, error) {
	versions := map[string]struct{}{}

	var walk func(dir string, depth int) error

	walk = func(dir string, depth int) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}

			return fmt.Errorf("Cannot read the directory: %s -> %w", dir, err)
		}

		for _, entry := range entries {
			path := filepath.Join(dir, entry.Name())

			// stat follows the links init leaves behind.
			if info, err := os.Stat(path); err != nil || !info.IsDir() {
				continue
			}

			if depth == layers-1 {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}

				versions[rel] = struct{}{}

				continue
			}

			if err := walk(path, depth+1); err != nil {
				return err
			}
		}

		return nil
	}

	if err := walk(root, 0); err != nil {
		return nil, err
	}

	return versions, nil
}

func size(dir string) int64 {
	var total int64

	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		if info, err := d.Info(); err == nil {
			total += info.Size()
		}

		return nil
	})

	return total
}
