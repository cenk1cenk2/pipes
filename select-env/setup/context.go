package setup

type Ctx struct {
	References  []string
	Environment string
	EnvVars     map[string]string
}
