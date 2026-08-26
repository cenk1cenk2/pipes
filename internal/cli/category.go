package cli

// Categories that more than one pipe registers flags under. A category that
// belongs to a single feature lives with that feature instead.
//
// The values are what the generated documentation prints as section headings, so
// they stay as they are until a step deliberately rewrites the docs.
const (
	CATEGORY_GIT    = "GIT"
	CATEGORY_CI     = "Gitlab Pipeline"
	CATEGORY_GITLAB = "Gitlab"
)
