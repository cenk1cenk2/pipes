package setup

type (
	EnvironmentConditionJson struct {
		Match       string `json:"match"       validate:"required"`
		Environment string `json:"environment" validate:"required"`
	}

	EnvironmentTemplate struct {
		Environment string
		EnvVars     map[string]string
	}
)
