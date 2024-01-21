package publish

type Ctx struct {
	Tags     []string
	Packages []PublishablePackage
}
