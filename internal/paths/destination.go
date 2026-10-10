package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Destination checks a path a pipe is about to replace completely under root. The
// path has to stay inside root, may not be root itself and may not reach into a .git
// directory, where the content written would replace the repository or run as a hook.
// A symlinked parent is resolved, since the removal would otherwise land outside root.
func Destination(root string, path string) error {
	if !filepath.IsLocal(path) || filepath.Clean(path) == "." {
		return fmt.Errorf("Destination has to be a path inside the target: %s", path)
	}

	if slices.ContainsFunc(strings.Split(filepath.ToSlash(filepath.Clean(path)), "/"), func(component string) bool {
		return strings.EqualFold(component, ".git")
	}) {
		return fmt.Errorf("Destination can not be inside a .git directory: %s", path)
	}

	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("Can not resolve the target: %s -> %w", root, err)
	}

	parent := filepath.Dir(filepath.Join(root, path))
	for {
		if _, err := os.Lstat(parent); err == nil {
			break
		}

		parent = filepath.Dir(parent)
	}

	resolvedParent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return fmt.Errorf("Can not resolve the destination: %s -> %w", path, err)
	}

	relative, err := filepath.Rel(resolvedRoot, resolvedParent)
	if err != nil || !filepath.IsLocal(relative) {
		return fmt.Errorf("Destination resolves outside the target: %s -> %s", path, resolvedParent)
	}

	return nil
}

// Disjoint fails when one destination contains another, since replacing the outer one
// would wipe what was just written into the inner one.
func Disjoint(paths []string) error {
	for i, a := range paths {
		for _, b := range paths[i+1:] {
			if contains(a, b) || contains(b, a) {
				return fmt.Errorf("Destinations overlap: %s -> %s", a, b)
			}
		}
	}

	return nil
}

func contains(parent string, child string) bool {
	relative, err := filepath.Rel(filepath.Clean(parent), filepath.Clean(child))

	return err == nil && filepath.IsLocal(relative)
}
