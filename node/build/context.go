package build

type Ctx struct {
	EnvironmentVariables []string
	SelectedEnvironment  string
	FallbackEnvironment  string
}
