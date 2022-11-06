package pipe

type (
	TagsSanitizeJson struct {
		Condition string `json:"condition" validate:"required"`
		Template  string `json:"template"  validate:"required"`
	}
)
