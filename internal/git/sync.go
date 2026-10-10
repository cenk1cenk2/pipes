package git

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"

	"gitlab.kilic.dev/devops/pipes/internal/paths"
)

// Path maps a source the job produced onto a destination in the target.
//
// Include and Exclude globs are anchored at the source and match the whole relative
// path, so `*.log` only matches at the top level and `cache/*.log` one level down;
// there is no `**`, so `src/gen/*.go` works and `**/*.go` does not. An entry is
// synced when it falls under an include and matches no exclude; an empty Include
// takes everything. A directory matching an exclude skips its tree, one matching an
// include takes its whole tree, and one that only leads towards an include is entered
// for what matches beneath it. The source itself is never matched.
type Path struct {
	Source      string
	Destination string
	Include     []string
	Exclude     []string
}

// Sync replaces every destination under root with its source. An artifact can not
// carry a deletion, so a destination is rebuilt from scratch instead of merged into.
// Every source is checked before anything is removed, since a source that failed to
// generate would otherwise wipe its destination; allowEmpty lifts that guard.
func Sync(t *Task, root string, entries []Path, allowEmpty bool) error {
	destinations := []string{}

	for _, p := range entries {
		if err := paths.Destination(root, p.Destination); err != nil {
			return err
		}

		destinations = append(destinations, p.Destination)

		for _, pattern := range p.Include {
			if _, err := filepath.Match(pattern, ""); err != nil {
				return fmt.Errorf("Can not process include pattern: %s -> %w", pattern, err)
			}
		}

		for _, pattern := range p.Exclude {
			if _, err := filepath.Match(pattern, ""); err != nil {
				return fmt.Errorf("Can not process exclude pattern: %s -> %w", pattern, err)
			}
		}

		if allowEmpty {
			continue
		}

		if err := ensureSource(p); err != nil {
			return err
		}
	}

	if err := paths.Disjoint(destinations); err != nil {
		return err
	}

	for _, p := range entries {
		destination := filepath.Join(root, p.Destination)

		if err := os.RemoveAll(destination); err != nil {
			return fmt.Errorf("Can not clear the destination: %s -> %w", destination, err)
		}

		if _, err := os.Lstat(p.Source); os.IsNotExist(err) {
			t.Log.Warn(fmt.Sprintf("Source does not exist, destination is left empty: %s -> %s", p.Source, p.Destination))

			continue
		}

		if err := copyTree(p, destination); err != nil {
			return fmt.Errorf("Can not sync the source: %s -> %s -> %w", p.Source, p.Destination, err)
		}

		t.Log.Info(fmt.Sprintf("Synced: %s -> %s", p.Source, p.Destination))
	}

	return nil
}

// A source with nothing but directories, or nothing the excludes let through, would
// still sync to an empty destination, so it takes a file to count as not empty.
func ensureSource(p Path) error {
	if _, err := os.Lstat(p.Source); err != nil {
		return fmt.Errorf("Source does not exist: %s -> %w", p.Source, err)
	}

	found := false

	if err := walk(p, func(_ string, _ string, d fs.DirEntry) error {
		if d.IsDir() {
			return nil
		}

		found = true

		return fs.SkipAll
	}); err != nil {
		return fmt.Errorf("Can not read the source: %s -> %w", p.Source, err)
	}

	if !found {
		return fmt.Errorf("Source is empty: %s", p.Source)
	}

	return nil
}

func copyTree(p Path, destination string) error {
	return walk(p, func(path string, relative string, d fs.DirEntry) error {
		target := filepath.Join(destination, relative)

		info, err := d.Info()
		if err != nil {
			return err
		}

		switch {
		case d.IsDir():
			return os.MkdirAll(target, info.Mode().Perm())
		case d.Type()&fs.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}

			return os.Symlink(link, target)
		default:
			return copyFile(path, target, info.Mode().Perm())
		}
	})
}

func walk(p Path, fn func(path string, relative string, d fs.DirEntry) error) error {
	return filepath.WalkDir(p.Source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relative, err := filepath.Rel(p.Source, path)
		if err != nil {
			return err
		}

		if relative == "." {
			return fn(path, relative, d)
		}

		for _, pattern := range p.Exclude {
			matched, err := filepath.Match(pattern, relative)
			if err != nil {
				return fmt.Errorf("Can not process exclude pattern: %s -> %w", pattern, err)
			}

			if !matched {
				continue
			}

			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		matched, err := included(p.Include, relative, d.IsDir())
		if err != nil {
			return err
		}

		if matched {
			return fn(path, relative, d)
		}

		if d.IsDir() {
			return filepath.SkipDir
		}

		return nil
	})
}

// A `*` never crosses a separator, so a pattern matches a path exactly when it has as
// many segments and each one matches, which lets a directory be held against just the
// leading segments of a pattern to tell whether anything beneath it can still match.
func included(include []string, relative string, dir bool) (bool, error) {
	if len(include) == 0 {
		return true, nil
	}

	segments := strings.Split(filepath.ToSlash(relative), "/")

	for _, pattern := range include {
		parts := strings.Split(pattern, "/")

		matched := 0
		for matched < len(parts) && matched < len(segments) {
			ok, err := filepath.Match(parts[matched], segments[matched])
			if err != nil {
				return false, fmt.Errorf("Can not process include pattern: %s -> %w", pattern, err)
			}

			if !ok {
				break
			}

			matched++
		}

		if matched == len(parts) || dir && matched == len(segments) {
			return true, nil
		}
	}

	return false, nil
}

func copyFile(source string, destination string, mode fs.FileMode) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}

	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		return errors.Join(err, out.Close())
	}

	return out.Close()
}
