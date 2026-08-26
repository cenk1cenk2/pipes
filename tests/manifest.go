package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

// Manifest is pipes.yaml: the list of pipes this repository publishes, and the
// only place that list is written down by hand.
type Manifest struct {
	Pipes []ManifestEntry `yaml:"pipes"`
}

type ManifestEntry struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
	// Readme is the path the pipeline uploads to the registry, relative to the
	// repository root.
	Readme      string `yaml:"readme"`
	Description string `yaml:"description"`
}

// ReadmeMatrix is the shape the update-docker-hub-readme pipe is handed through
// README_MATRIX. The manifest has to agree with it, since the pipeline is what
// actually publishes the descriptions.
type ReadmeMatrix []ReadmeMatrixEntry

type ReadmeMatrixEntry struct {
	Repository  string `json:"repository"`
	File        string `json:"file"`
	Description string `json:"description"`
}

// pipeline is the part of .gitlab-ci.yml the manifest is checked against.
type pipeline struct {
	UpdateDockerHubReadme struct {
		Variables struct {
			ReadmeMatrix string `yaml:"README_MATRIX"`
		} `yaml:"variables"`
	} `yaml:"update-docker-hub-readme"`
}

// Root is the repository root. The specs run from the module directory, so
// everything they read is addressed from one level up.
func Root() string {
	root, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}

	return root
}

func ReadManifest() (Manifest, error) {
	var manifest Manifest

	contents, err := os.ReadFile(filepath.Join(Root(), "pipes.yaml"))
	if err != nil {
		return manifest, err
	}

	if err := yaml.Unmarshal(contents, &manifest); err != nil {
		return manifest, fmt.Errorf("Can not unmarshal pipes.yaml: %w", err)
	}

	return manifest, nil
}

func ReadReadmeMatrix() (ReadmeMatrix, error) {
	var (
		ci     pipeline
		matrix ReadmeMatrix
	)

	contents, err := os.ReadFile(filepath.Join(Root(), ".gitlab-ci.yml"))
	if err != nil {
		return matrix, err
	}

	if err := yaml.Unmarshal(contents, &ci); err != nil {
		return matrix, fmt.Errorf("Can not unmarshal .gitlab-ci.yml: %w", err)
	}

	raw := ci.UpdateDockerHubReadme.Variables.ReadmeMatrix

	if raw == "" {
		return matrix, fmt.Errorf("README_MATRIX is not set on the update-docker-hub-readme job")
	}

	if err := json.Unmarshal([]byte(raw), &matrix); err != nil {
		return matrix, fmt.Errorf("Can not unmarshal README_MATRIX: %w", err)
	}

	return matrix, nil
}

// ModuleDirs are the directories of the workspace that hold a Go module, which
// is where a pipe that nothing has registered yet would show up.
func ModuleDirs() ([]string, error) {
	entries, err := os.ReadDir(Root())
	if err != nil {
		return nil, err
	}

	dirs := []string{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if _, err := os.Stat(filepath.Join(Root(), entry.Name(), "go.mod")); err != nil {
			continue
		}

		dirs = append(dirs, entry.Name())
	}

	return dirs, nil
}
