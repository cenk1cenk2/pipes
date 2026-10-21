package install

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	"golang.org/x/mod/module"
)

// Only one of the vendor tasks fires, since the gates below are the two halves
// of the workspace condition the setup resolved.
func vendor(tl *TaskList) *Task {
	return tl.CreateTask("vendor").
		SetJobWrapper(func(_ Job, t *Task) Job {
			return JobParallel(
				vendorModule(tl).Job(),
				vendorWorkspace(tl).Job(),
			)
		})
}

func vendorModule(tl *TaskList) *Task {
	return tl.CreateTask("vendor", "module").
		ShouldDisable(func(_ *Task) bool {
			return setup.C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"mod",
				"vendor",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Vendoring: in %s", setup.C.Cwd))

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func vendorWorkspace(tl *TaskList) *Task {
	return tl.CreateTask("vendor", "workspace").
		ShouldDisable(func(_ *Task) bool {
			return !setup.C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"work",
				"vendor",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Vendoring workspace: in %s", setup.C.Cwd))

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func verify(tl *TaskList) *Task {
	return tl.CreateTask("verify").
		ShouldDisable(func(t *Task) bool {
			return !P.Verify
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"mod",
				"verify",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Verifying modules: in %s", setup.C.Cwd))

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

// the download cache keeps the zip of every module version go ever fetched, so the
// cached directory only stays the size of the project when the zips of the versions
// the build list left behind are dropped. The go.mod and info files stay, since
// resolving the module graph reads them for versions it does not select.
func prune(tl *TaskList) *Task {
	return tl.CreateTask("prune").
		ShouldDisable(func(_ *Task) bool {
			return setup.C.Env["GOMODCACHE"] == ""
		}).
		Set(func(_ context.Context, t *Task) error {
			var output string

			// the vendor directory refuses to compute all, and readonly is the one
			// override a workspace accepts as well.
			t.CreateCommand(
				"go",
				"list",
				"-mod=readonly",
				"-m",
				"-json",
				"all",
			).
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				SetDir(setup.C.Cwd).
				CaptureStdout(&output).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					keep, err := buildList(output)
					if err != nil {
						return err
					}

					download := filepath.Join(setup.C.Env["GOMODCACHE"], "cache", "download")

					files, err := stale(download, keep)
					if err != nil {
						return err
					}

					var size int64

					for _, file := range files {
						info, err := os.Stat(file)
						if err != nil {
							return fmt.Errorf("Cannot stat the cached module file: %s -> %w", file, err)
						}

						if err := os.Remove(file); err != nil {
							return fmt.Errorf("Cannot remove the cached module file: %s -> %w", file, err)
						}

						size += info.Size()
					}

					t.Log.Info(
						fmt.Sprintf("Pruned module versions outside the build list: %d files, %d bytes in %s", len(files), size, download),
					)

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

// a replaced module is downloaded under its replacement, so both count as used.
func buildList(output string) (map[module.Version]bool, error) {
	type listed struct {
		Path    string
		Version string
		Replace *listed
	}

	keep := map[module.Version]bool{}

	decoder := json.NewDecoder(strings.NewReader(output))

	for {
		var m listed

		if err := decoder.Decode(&m); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Cannot parse the build list of the modules: %w", err)
		}

		keep[module.Version{Path: m.Path, Version: m.Version}] = true

		if m.Replace != nil && m.Replace.Version != "" {
			keep[module.Version{Path: m.Replace.Path, Version: m.Replace.Version}] = true
		}
	}

	if len(keep) == 0 {
		return nil, fmt.Errorf("Can not resolve the build list of the modules.")
	}

	return keep, nil
}

// the checksum database and the toolchains go switches to are no modules of the
// build list, so they are never candidates.
func stale(download string, keep map[module.Version]bool) ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(download, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			if file == download && errors.Is(err, fs.ErrNotExist) {
				return fs.SkipAll
			}

			return err
		}

		if !entry.IsDir() || entry.Name() != "@v" {
			return nil
		}

		relative, err := filepath.Rel(download, filepath.Dir(file))
		if err != nil {
			return err
		}

		escaped := filepath.ToSlash(relative)

		if escaped == "sumdb" || strings.HasPrefix(escaped, "sumdb/") || escaped == "golang.org/toolchain" {
			return fs.SkipDir
		}

		p, err := module.UnescapePath(escaped)
		if err != nil {
			return fs.SkipDir
		}

		versions, err := os.ReadDir(file)
		if err != nil {
			return err
		}

		for _, version := range versions {
			ext := path.Ext(version.Name())

			if version.IsDir() || (ext != ".zip" && ext != ".ziphash" && ext != ".lock") {
				continue
			}

			v, err := module.UnescapeVersion(strings.TrimSuffix(version.Name(), ext))
			if err != nil || keep[module.Version{Path: p, Version: v}] {
				continue
			}

			files = append(files, filepath.Join(file, version.Name()))
		}

		return fs.SkipDir
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot walk the module download cache: %s -> %w", download, err)
	}

	return files, nil
}
