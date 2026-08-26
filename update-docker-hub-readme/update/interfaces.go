package update

type (
	ReadmeMatrixJson struct {
		Repository  string `json:"repository"`
		File        string `json:"file"`
		Description string `json:"description,omitempty"`
	}
)

type (
	ParsedReadme struct {
		File        string
		Description string
	}
)
