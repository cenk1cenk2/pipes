package pipe

type (
	TagsSanitizeJson struct {
		Match    string `json:"match"    validate:"required"`
		Template string `json:"template" validate:"required"`
	}

	TagsTemplateJson struct {
		Match    string `json:"match"    validate:"required"`
		Template string `json:"template" validate:"required"`
	}
)
