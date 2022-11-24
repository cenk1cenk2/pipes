package manifest

type (
	DockerManifestMatrixJson struct {
		Target string   `json:"target,omitempty"`
		Images []string `json:"images"           validate:"required"`
	}
)
