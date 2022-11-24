package manifest

type Ctx struct {
	ManifestedImages map[string][]string
	Matches          []string
}
