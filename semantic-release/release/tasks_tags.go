package release

import (
	"context"
	json "encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

// The tag is already pushed once the release exits, so nothing here may fail the job;
// a report that can not be written whole is not written at all.
func tags(tl *TaskList) *Task {
	return tl.CreateTask("tags").
		ShouldDisable(func(t *Task) bool {
			return P.SemanticRelease.DryRun
		}).
		Set(func(_ context.Context, t *Task) error {
			env, err := tagsEnv()

			if err != nil {
				t.Log.Warn(fmt.Sprintf("Skipping the dotenv report of the released tags: %s", err))

				return nil
			}

			if len(env) == 0 {
				t.Log.Info("Nothing was released, skipping the dotenv report.")

				return nil
			}

			keys := slices.Sorted(maps.Keys(env))
			lines := make([]string, 0, len(keys))

			for _, key := range keys {
				lines = append(lines, fmt.Sprintf("%s=%s\n", key, env[key]))
			}

			if err := os.WriteFile(P.SemanticRelease.EnvFile, []byte(strings.Join(lines, "")), 0o600); err != nil {
				t.Log.Warn(fmt.Sprintf("Skipping the dotenv report of the released tags: %s", err))

				return nil
			}

			t.Log.Info(fmt.Sprintf("Released tags written to dotenv file: %s > %s", P.SemanticRelease.EnvFile, strings.Join(keys, ", ")))

			return nil
		})
}

func tagsEnv() (map[string]string, error) {
	files, err := tagsFiles()

	if err != nil {
		return nil, err
	}

	env := map[string]string{}

	for _, file := range files {
		content, err := os.ReadFile(file)

		if err != nil {
			return nil, err
		}

		value := strings.TrimSpace(string(content))

		if value == "" || strings.ContainsAny(value, "\r\n\"'`#") {
			return nil, fmt.Errorf("Tags file can only hold a single tag: %s", file)
		}

		key := P.SemanticRelease.EnvVariable

		if P.SemanticRelease.Workspace {
			name, err := packageName(filepath.Dir(file))

			if err != nil {
				return nil, err
			}

			key = fmt.Sprintf("%s_%s", P.SemanticRelease.EnvVariable, name)
		}

		if _, ok := env[key]; ok {
			return nil, fmt.Errorf("More than one tags file maps to the same variable: %s", key)
		}

		env[key] = value
	}

	return env, nil
}

// The release only writes a tags file for what it released, which in a workspace is
// what tells the released packages apart from the rest.
func tagsFiles() ([]string, error) {
	if filepath.IsAbs(P.SemanticRelease.TagsFile) {
		if _, err := os.Stat(P.SemanticRelease.TagsFile); errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		return []string{P.SemanticRelease.TagsFile}, nil
	}

	name := filepath.Clean(P.SemanticRelease.TagsFile)
	files := []string{}

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}

			return nil
		}

		if path != name && !strings.HasSuffix(path, string(filepath.Separator)+name) {
			return nil
		}

		files = append(files, path)

		return nil
	})

	return files, err
}

// The package name keys the tag of a workspace package, reduced to what GitLab takes
// in a variable name.
func packageName(dir string) (string, error) {
	content, err := os.ReadFile(filepath.Join(dir, "package.json"))

	if err != nil {
		return "", err
	}

	pkg := struct {
		Name string `json:"name"`
	}{}

	if err := json.Unmarshal(content, &pkg); err != nil {
		return "", fmt.Errorf("Can not unmarshal package.json: %s -> %w", dir, err)
	}

	if pkg.Name == "" {
		return "", fmt.Errorf("Package does not have a name: %s", dir)
	}

	return regexp.MustCompile(`[^A-Za-z0-9]`).
		ReplaceAllString(strings.ReplaceAll(strings.ToUpper(pkg.Name), "@", ""), "_"), nil
}
