package pipe

type Ctx struct {
	Tags                       []string
	UseExistingBuildXInstance  bool
	SanitizedRegularExpression map[string]string
}
