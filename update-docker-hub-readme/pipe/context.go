package pipe

type Ctx struct {
	Token       string
	Readme      map[string]ParsedReadme
	ReadmeFiles map[string][]byte
}
