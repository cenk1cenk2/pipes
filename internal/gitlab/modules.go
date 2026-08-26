package gitlab

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// The Terraform module registry, narrowed to the one call the publish pipe makes so
// the upload can be driven without a GitLab to talk to.
type ModuleRegistry interface {
	UploadModule(ctx context.Context, name, system, version string, archive io.Reader) error
}

type moduleRegistry struct {
	apiUrl    string
	projectId string
	token     string
	client    *http.Client
}

var _ ModuleRegistry = (*moduleRegistry)(nil)

func NewModuleRegistry(apiUrl, projectId, token string) ModuleRegistry {
	return &moduleRegistry{
		apiUrl:    apiUrl,
		projectId: projectId,
		token:     token,
		client:    &http.Client{},
	}
}

func (r *moduleRegistry) UploadModule(
	ctx context.Context,
	name, system, version string,
	archive io.Reader,
) error {
	url := fmt.Sprintf(
		"%s/projects/%s/packages/terraform/modules/%s/%s/%s/file",
		r.apiUrl,
		r.projectId,
		name,
		system,
		version,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, archive)
	if err != nil {
		return fmt.Errorf("create GitLab module upload request: %w", err)
	}

	req.Header.Set("Content-Type", "application/tar+gzip")
	req.Header.Set("JOB-TOKEN", r.token)

	res, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("upload module %s@%s: %w", name, version, err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read GitLab module upload response: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"upload module %s@%s: %s: %s",
			name,
			version,
			res.Status,
			strings.TrimSpace(string(body)),
		)
	}

	return nil
}
