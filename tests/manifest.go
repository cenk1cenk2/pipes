// Package tests holds the checks that are about the repository and not about any
// one pipe. Nothing here imports a pipe, which lets a pipe keep its command tree
// to itself.
package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

// ReadmeMatrixEntry is one repository of DOCKER_HUB_README_MATRIX, which is what
// the update-docker-hub-readme job actually publishes to Docker Hub.
type ReadmeMatrixEntry struct {
	Repository  string `json:"repository"`
	File        string `json:"file"`
	Description string `json:"description"`
}

// pipeline is the part of .gitlab-ci.yml the specs read.
type pipeline struct {
	UpdateDockerHubReadme struct {
		Variables struct {
			ReadmeMatrix string `yaml:"DOCKER_HUB_README_MATRIX"`
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

func ReadReadmeMatrix() ([]ReadmeMatrixEntry, error) {
	var ci pipeline

	contents, err := os.ReadFile(filepath.Join(Root(), ".gitlab-ci.yml"))
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(contents, &ci); err != nil {
		return nil, fmt.Errorf("Can not unmarshal .gitlab-ci.yml: %w", err)
	}

	raw := ci.UpdateDockerHubReadme.Variables.ReadmeMatrix

	if raw == "" {
		return nil, fmt.Errorf("DOCKER_HUB_README_MATRIX is not set on the update-docker-hub-readme job")
	}

	var matrix []ReadmeMatrixEntry

	if err := json.Unmarshal([]byte(raw), &matrix); err != nil {
		return nil, fmt.Errorf("Can not unmarshal DOCKER_HUB_README_MATRIX: %w", err)
	}

	return matrix, nil
}
