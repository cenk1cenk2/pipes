// Package manifest turns the rendered resources of a base and a head side into one git
// patch with a file per resource, so a merge request shows what a change does to the
// manifests rather than to their sources.
package manifest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	k8syaml "sigs.k8s.io/yaml"
)

const (
	// ClusterScope stands in for the namespace of a resource that has none.
	ClusterScope = "_cluster"
	// CoreGroup stands in for the api group of a resource in the core group.
	CoreGroup = "core"
)

type (
	// Group is one overlay or chart, rendered on both sides. A nil Base is a group the
	// base does not have, so every resource of it is added. A group with BaseErr is
	// left out of the diff and counts nothing, the caller reports the error instead.
	Group struct {
		Name    string
		Base    []*yaml.RNode
		Head    []*yaml.RNode
		BaseErr error
	}

	// Counts is how many resource files of a group the patch adds, changes and removes.
	Counts struct {
		Added   int
		Changed int
		Removed int
	}

	// Result is the patch of every group, and the counts of each group by its name.
	Result struct {
		Patch  string
		Groups map[string]Counts
	}
)

// Diff writes both sides into a scratch repository and diffs them there, which gives
// the same a/ and b/ patch a commit has, also when one side is empty.
func Diff(ctx context.Context, t *Task, groups []Group) (Result, error) {
	result := Result{Groups: map[string]Counts{}}
	dirs := map[string]string{}
	base := map[string][]byte{}
	head := map[string][]byte{}

	for _, group := range groups {
		result.Groups[group.Name] = Counts{}

		if group.BaseErr != nil {
			continue
		}

		dir, err := segment(group.Name)
		if err != nil {
			return Result{}, err
		}

		if other, ok := dirs[dir]; ok {
			return Result{}, fmt.Errorf("Groups write into the same directory: %s, %s -> %s", other, group.Name, dir)
		}
		dirs[dir] = group.Name

		if err := render(base, dir, group.Base); err != nil {
			return Result{}, fmt.Errorf("Can not render the base of the group: %s -> %w", group.Name, err)
		}

		if err := render(head, dir, group.Head); err != nil {
			return Result{}, fmt.Errorf("Can not render the head of the group: %s -> %w", group.Name, err)
		}
	}

	scratch, err := os.MkdirTemp("", "manifest-diff-")
	if err != nil {
		return Result{}, fmt.Errorf("Can not create the scratch directory: %w", err)
	}
	defer os.RemoveAll(scratch)

	repository := filepath.Join(scratch, "repository")
	output := filepath.Join(scratch, "diff.patch")

	// the configuration of the host could prefix, color or sign what git writes here,
	// and the patch has to come out the same wherever the job runs.
	git := func(args ...string) *Command {
		return t.CreateCommand("git", args...).
			SetDir(repository).
			SetLogLevel(LogLevelDebug, LogLevelWarn, LogLevelDebug).
			AppendEnvironment(map[string]string{
				"GIT_CONFIG_GLOBAL":   os.DevNull,
				"GIT_CONFIG_NOSYSTEM": "1",
				"GIT_AUTHOR_NAME":     "gitlab-pipes",
				"GIT_AUTHOR_EMAIL":    "gitlab-pipes@localhost",
				"GIT_AUTHOR_DATE":     "1970-01-01T00:00:00Z",
				"GIT_COMMITTER_NAME":  "gitlab-pipes",
				"GIT_COMMITTER_EMAIL": "gitlab-pipes@localhost",
				"GIT_COMMITTER_DATE":  "1970-01-01T00:00:00Z",
			})
	}

	if err := write(repository, base); err != nil {
		return Result{}, err
	}

	var status string

	for _, command := range []*Command{
		git("init", "--quiet", "--initial-branch", "main"),
		git("add", "-A"),
		// no detached maintenance run, which would outlive the scratch repository.
		git("-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "commit", "--quiet", "--allow-empty", "--message", "base"),
	} {
		if err := command.Run(ctx); err != nil {
			return Result{}, fmt.Errorf("Can not commit the base side: %w", err)
		}
	}

	entries, err := os.ReadDir(repository)
	if err != nil {
		return Result{}, fmt.Errorf("Can not read the base side: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}

		if err := os.RemoveAll(filepath.Join(repository, entry.Name())); err != nil {
			return Result{}, fmt.Errorf("Can not clear the base side: %w", err)
		}
	}

	if err := write(repository, head); err != nil {
		return Result{}, err
	}

	for _, command := range []*Command{
		git("add", "-A"),
		git("diff", "--cached", "--no-renames", "--binary", fmt.Sprintf("--output=%s", output)),
		git("-c", "core.quotePath=false", "diff", "--cached", "--name-status", "--no-renames").
			CaptureStdout(&status),
	} {
		if err := command.Run(ctx); err != nil {
			return Result{}, fmt.Errorf("Can not diff the head side against the base side: %w", err)
		}
	}

	patch, err := os.ReadFile(output)
	if err != nil {
		return Result{}, fmt.Errorf("Can not read the patch: %w", err)
	}
	result.Patch = string(patch)

	for line := range strings.SplitSeq(status, "\n") {
		if line == "" {
			continue
		}

		change, path, ok := strings.Cut(line, "\t")
		if !ok {
			return Result{}, fmt.Errorf("Can not parse the changed file: %s", line)
		}

		if strings.HasPrefix(path, `"`) {
			if path, err = strconv.Unquote(path); err != nil {
				return Result{}, fmt.Errorf("Can not unquote the changed file: %s -> %w", line, err)
			}
		}

		dir, _, _ := strings.Cut(path, "/")
		name, ok := dirs[dir]
		if !ok {
			return Result{}, fmt.Errorf("Changed file belongs to no group: %s", path)
		}

		counts := result.Groups[name]
		switch change {
		case "A":
			counts.Added++
		case "M":
			counts.Changed++
		case "D":
			counts.Removed++
		default:
			return Result{}, fmt.Errorf("Changed file has an unexpected status: %s", line)
		}
		result.Groups[name] = counts
	}

	return result, nil
}

func write(repository string, files map[string][]byte) error {
	if err := os.MkdirAll(repository, 0o755); err != nil {
		return fmt.Errorf("Can not create the scratch repository: %w", err)
	}

	for path, content := range files {
		destination := filepath.Join(repository, filepath.FromSlash(path))

		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("Can not create the directory of the manifest: %s -> %w", path, err)
		}

		if err := os.WriteFile(destination, content, 0o644); err != nil {
			return fmt.Errorf("Can not write the manifest: %s -> %w", path, err)
		}
	}

	return nil
}

// render adds every resource to files under the path that identifies it, sorted-key
// YAML so that the order the renderer emitted fields in never shows.
func render(files map[string][]byte, dir string, resources []*yaml.RNode) error {
	for _, resource := range resources {
		path, err := locate(dir, resource)
		if err != nil {
			return err
		}

		if _, ok := files[path]; ok {
			return fmt.Errorf("Resource is rendered more than once: %s", path)
		}

		fields, err := resource.Map()
		if err != nil {
			return fmt.Errorf("Can not read the resource: %s -> %w", path, err)
		}

		content, err := k8syaml.Marshal(fields)
		if err != nil {
			return fmt.Errorf("Can not serialize the resource: %s -> %w", path, err)
		}

		files[path] = content
	}

	return nil
}

// locate is <group>/<namespace>/<kind>.<api group>/<name>.yaml, so the same resource
// lands on the same path on both sides.
func locate(dir string, resource *yaml.RNode) (string, error) {
	namespace := resource.GetNamespace()
	if namespace == "" {
		namespace = ClusterScope
	}

	group, _, ok := strings.Cut(resource.GetApiVersion(), "/")
	if !ok {
		group = CoreGroup
	}

	segments := []string{dir}
	for _, value := range []string{namespace, resource.GetKind(), group, resource.GetName()} {
		s, err := segment(value)
		if err != nil {
			return "", err
		}

		segments = append(segments, s)
	}

	return fmt.Sprintf("%s/%s/%s.%s/%s.yaml", segments[0], segments[1], segments[2], segments[3], segments[4]), nil
}

// segment keeps a value inside the one directory level it names, a separator in it
// would otherwise nest it or climb out of the scratch repository.
func segment(value string) (string, error) {
	s := strings.NewReplacer("/", "_", `\`, "_").Replace(value)

	switch s {
	case "", ".", "..", ".git":
		return "", fmt.Errorf("Can not use the value as a path segment of a manifest: %q", value)
	}

	return s, nil
}
