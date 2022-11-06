package setup

type (
	EnvironmentConditionJson struct {
		Condition   string `json:"condition"   validate:"required"`
		Environment string `json:"environment" validate:"required"`
	}
)
