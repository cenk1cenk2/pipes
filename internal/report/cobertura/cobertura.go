// Package cobertura writes the GitLab artifacts:reports:coverage_report XML for test
// runners that have no Cobertura writer of their own, through the adapter the pipe
// builds for its runner's coverage output.
package cobertura

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aereal/gobertura"
	. "github.com/cenk1cenk2/plumber/v7"
	"golang.org/x/mod/modfile"
)

// ConverterAdapter converts the coverage output a test run left at input into a
// Cobertura report at output. dir is the directory the run happened in, which the
// paths of the report are resolved against.
type ConverterAdapter interface {
	Convert(input string, output string, dir string) error
}

// GoProfileConverter converts a Go coverage profile. The profile names files by import
// path, while GitLab matches them against the repository, so dir is the module the
// profile was written for: it is the source of the report and every file under its
// module path is named relative to it.
type GoProfileConverter struct{}

var _ ConverterAdapter = GoProfileConverter{}

func (GoProfileConverter) Convert(input string, output string, dir string) error {
	source, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("resolve module directory %s: %w", dir, err)
	}

	mod, err := os.ReadFile(filepath.Join(source, "go.mod"))
	if err != nil {
		return fmt.Errorf("read module of the coverage profile: %w", err)
	}

	in, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("open coverage profile %s: %w", input, err)
	}
	defer in.Close()

	coverage, err := gobertura.Parse(in, time.Now())
	if err != nil {
		return fmt.Errorf("parse coverage profile %s: %w", input, err)
	}

	module := modfile.ModulePath(mod) + "/"
	for i := range coverage.Packages {
		for j := range coverage.Packages[i].Classes {
			class := &coverage.Packages[i].Classes[j]
			class.Filename = strings.TrimPrefix(class.Filename, module)
		}
	}
	coverage.Sources.Sources = []string{source}

	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create cobertura report %s: %w", output, err)
	}
	defer out.Close()

	if err := gobertura.Write(out, coverage); err != nil {
		return fmt.Errorf("write cobertura report %s: %w", output, err)
	}

	return nil
}

// Everything a Cobertura report needs from the task of the pipe that converts it.
type Source struct {
	Adapter ConverterAdapter
	Input   string
	Output  string
	// the directories the tests ran in, each holding the coverage output of its own module.
	Dirs []string
}

// Convert converts the coverage output in every directory of the source through its
// adapter, logging to the task it runs in. A directory the tests never reached under
// fail-fast leaves nothing to convert.
func (s Source) Convert(t *Task) error {
	if len(s.Dirs) > 1 {
		for _, path := range []string{s.Input, s.Output} {
			if !filepath.IsLocal(path) {
				return fmt.Errorf("Cobertura report paths have to be inside every directory, otherwise the directories overwrite each other's report: %s", path)
			}
		}
	}

	var errs []error

	for _, dir := range s.Dirs {
		input := resolve(dir, s.Input)
		if _, err := os.Stat(input); errors.Is(err, os.ErrNotExist) {
			t.Log.Warn(fmt.Sprintf("Skipping the Cobertura report because there is no coverage output: %s", input))

			continue
		}

		output := resolve(dir, s.Output)
		if err := s.Adapter.Convert(input, output, dir); err != nil {
			errs = append(errs, fmt.Errorf("Can not write the Cobertura report of %s: %w", dir, err))

			continue
		}

		t.Log.Info(fmt.Sprintf("Wrote Cobertura report: %s", output))
	}

	return errors.Join(errs...)
}

func resolve(dir string, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(dir, path)
}
