package patch

import (
	"regexp"
	"strings"
)

var fenceRun = regexp.MustCompile("`{3,}")

type (
	Section struct {
		Count int
		Files []File
		More  int
	}

	File struct {
		Name  string
		Diff  string
		Fence string
	}
)

// Cuts a patch into the diffs of the files it changes, each with its own header.
func Split(diff string) []string {
	files := []string{}
	for diff != "" {
		end := len(diff)
		if i := strings.Index(diff, "\ndiff --git "); i >= 0 {
			end = i + 1
		}

		files = append(files, diff[:end])
		diff = diff[end:]
	}

	return files
}

// Takes the files into one section each within the limit, and returns what is left of
// it together with whether any file did not fit. A file cut at a line ends the diff,
// so nothing after it can make it look complete.
func Sections(files []string, limit int) (Section, int, bool) {
	section := Section{Count: len(files)}

	// the name and the fences of every file count against the limit as well, which a
	// diff of many small files would otherwise outgrow.
	cut := false
	for _, file := range files {
		name := Filename(file)
		diff := Hunks(file)
		marker := Fence(diff)
		overhead := len(name) + 2*len(marker) + len("\n``\n\ndiff\n\n")
		if overhead >= limit {
			break
		}

		if len(diff) > limit-overhead {
			end := strings.LastIndex(diff[:limit-overhead], "\n") + 1
			if end == 0 {
				break
			}

			diff = diff[:end]
			cut = true
		}

		section.Files = append(section.Files, File{Name: name, Diff: diff, Fence: marker})
		limit -= overhead + len(diff)

		if cut {
			limit = 0

			break
		}
	}

	section.More = len(files) - len(section.Files)

	return section, limit, cut || section.More > 0
}

// The name a file carries after the change, which is the one a rename ends up with.
func Filename(diff string) string {
	header, _, _ := strings.Cut(diff, "\n")
	header = strings.TrimPrefix(header, "diff --git ")
	if i := strings.LastIndex(header, " b/"); i >= 0 {
		return header[i+len(" b/"):]
	}

	return header
}

// The name a file carries before the change, which a rename leaves behind.
func Previous(diff string) string {
	header, _, _ := strings.Cut(diff, "\n")
	header = strings.TrimPrefix(header, "diff --git a/")
	if i := strings.LastIndex(header, " b/"); i >= 0 {
		return header[:i]
	}

	return header
}

// The name above the section already says which file it is, so only what git marks
// beyond that stays: the hunks, and the mode, rename and binary lines of the header.
func Hunks(diff string) string {
	body := &strings.Builder{}
	header := true
	for i, line := range strings.SplitAfter(diff, "\n") {
		if i == 0 {
			continue
		}

		if strings.HasPrefix(line, "@@") || line == "GIT binary patch\n" {
			header = false
		}

		if header && (strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ")) {
			continue
		}

		body.WriteString(line)
	}

	return body.String()
}

// A diff of a markdown file can carry a fence of its own, which would close this one.
func Fence(diff string) string {
	fence := "```"
	for _, run := range fenceRun.FindAllString(diff, -1) {
		if len(run) >= len(fence) {
			fence = strings.Repeat("`", len(run)+1)
		}
	}

	return fence
}
