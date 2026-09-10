// Package markdown renders the markdown a pipe reports with into the text a pipeline
// log viewer shows, so a report reaches the job log without its reader having to work
// through the markup itself.
package markdown

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/muesli/termenv"
)

const (
	escape = '\x1b'
	reset  = "\x1b[0m"
)

var (
	// the escape sequences and the padding that close a line of a rendered document,
	// which a log viewer shows as noise where a terminal would have swallowed it.
	trailing = regexp.MustCompile(`(?:\x1b\[[0-9;]*m|[ \t])+$`)
	escapes  = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

// The palette stays inside the sixteen colors that every log viewer agrees on, which is
// the same ground the logger of plumber keeps to.
var style = ansi.StyleConfig{
	Document:  ansi.StyleBlock{Margin: new(uint(0))},
	Paragraph: ansi.StyleBlock{},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{BlockSuffix: "\n", Bold: new(true), Color: new("6")},
	},
	// the levels only carry the number sign the renderer prefixes them with, which
	// reads as markup rather than as a heading once the document is not markup anymore.
	H1:   ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: ""}},
	H2:   ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: ""}},
	H3:   ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: ""}},
	H4:   ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: ""}},
	H5:   ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: ""}},
	H6:   ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Prefix: ""}},
	Code: ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Color: new("3")}},
	// a code block sits under the line that opens its folded section. Chroma only takes
	// hex colors and matches them against a table of its own, so these are the entries
	// of its sixteen color table rather than the indexes used above; anything else is
	// dropped rather than approximated.
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{BlockSuffix: "\n"}, Indent: new(uint(2))},
		Chroma: &ansi.Chroma{
			GenericInserted:   ansi.StylePrimitive{Color: new("#007f00")},
			GenericDeleted:    ansi.StylePrimitive{Color: new("#7f0000")},
			GenericSubheading: ansi.StylePrimitive{Color: new("#7f7fe0")},
			GenericStrong:     ansi.StylePrimitive{Color: new("#ff0000"), Bold: new(true)},
			Keyword:           ansi.StylePrimitive{Color: new("#00007f")},
			KeywordType:       ansi.StylePrimitive{Color: new("#007f7f")},
			Error:             ansi.StylePrimitive{Color: new("#7f007f")},
			Comment:           ansi.StylePrimitive{Color: new("#555555")},
		},
	},
	// the line that opens a folded section reaches the terminal as the text of its own
	// html, which carries no markup left to style: bold is what separates one resource
	// from the block of attributes under it.
	HTMLBlock: ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: new(true)}},
	Link:      ansi.StylePrimitive{Color: new("4"), Underline: new(true)},
	LinkText:  ansi.StylePrimitive{Color: new("4")},
	List:      ansi.StyleList{LevelIndent: 2},
	Item:      ansi.StylePrimitive{BlockPrefix: "- "},
	Emph:      ansi.StylePrimitive{Italic: new(true)},
	Strong:    ansi.StylePrimitive{Bold: new(true)},
}

// Render turns a markdown document into the styled text of a terminal.
func Render(body string) (string, error) {
	profile := colorProfile()

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(style),
		// the document is read in a log viewer that never reflows it, so it keeps its
		// natural width instead of being wrapped into a column.
		glamour.WithWordWrap(0),
		// the footnotes a table would otherwise collect are truncated against the word
		// wrap, which leaves every one of them without its url once the wrap is off.
		glamour.WithInlineTableLinks(true),
		glamour.WithColorProfile(profile),
		glamour.WithChromaFormatter("terminal16"),
	)
	if err != nil {
		return "", fmt.Errorf("create markdown renderer: %w", err)
	}

	rendered, err := renderer.Render(body)
	if err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}

	// a profile without colors still leaves the attributes that carry none behind, and
	// an environment that asked for no styling wants none of it.
	if profile == termenv.Ascii {
		rendered = escapes.ReplaceAllString(rendered, "")
	}

	lines := strings.Split(strings.Trim(rendered, "\n"), "\n")
	for index, line := range lines {
		line = trailing.ReplaceAllString(line, "")

		// the padding a line closes with buries its reset, which the trim takes along
		// with it and would leave the line bleeding its color into the ones after it.
		if strings.ContainsRune(line, escape) {
			line += reset
		}

		lines[index] = line
	}

	// a block that closes the document leaves its padding on a line of its own, which
	// only shows once the padding is gone.
	for len(lines) > 0 && escapes.ReplaceAllString(lines[len(lines)-1], "") == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n"), nil
}

// Log writes a markdown document to the pipeline log, one record per line, which is how
// the output of a command reaches it as well.
func Log(log *slog.Logger, body string) error {
	rendered, err := Render(body)
	if err != nil {
		return err
	}

	for line := range strings.SplitSeq(rendered, "\n") {
		log.Info(line)
	}

	return nil
}

// Mirrors the logger of plumber, which forces the profile instead of detecting it
// because the output of a pipe ends up in a pipeline log viewer that renders the escape
// sequences although it is never a terminal.
func colorProfile() termenv.Profile {
	if force := os.Getenv("CLICOLOR_FORCE"); force != "" && force != "0" {
		return termenv.ANSI
	}

	if os.Getenv("NO_COLOR") != "" {
		return termenv.Ascii
	}

	return termenv.ANSI
}
