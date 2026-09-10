package markdown

import (
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// The glyph a line of a plan diff opens with, and the token it is colored through.
// Compound glyphs come first, so one is never read as the single character it starts
// with.
var diffGlyphs = []struct {
	glyph string
	token chroma.TokenType
}{
	{"-/+", chroma.GenericStrong},
	{"<=", chroma.Keyword},
	{"+", chroma.GenericInserted},
	{"-", chroma.GenericDeleted},
	{"~", chroma.GenericSubheading},
	{">", chroma.KeywordType},
	{".", chroma.Comment},
	{"?", chroma.Error},
}

// What a plan names as the reason a resource cannot be changed in place.
const forcesReplacement = "# forces replacement"

/*
RegisterDiffLexer teaches the renderer the diff blocks a pipe writes, taking the text
it stands in for a value it will not show so that it is dimmed wherever it turns up.

Chroma reads a diff block with its own lexer, which only knows unified diffs: a line
opening with neither a plus nor a minus is context there and stays unstyled, and so is
every line opening with whitespace. That covers every changed attribute and every
nested line of a plan, which is most of one. The lexer built here colors a line by the
glyph it carries instead, and claims the alias of the one it replaces so that the block
stays a diff block for the merge request note, where GitLab highlights it itself.
*/
func RegisterDiffLexer(dimmed ...string) {
	quoted := make([]string, 0, len(dimmed))
	for _, value := range dimmed {
		quoted = append(quoted, regexp.QuoteMeta(value))
	}
	dim := strings.Join(quoted, "|")

	// a line is read to its end in a state of its own, so the glyph rules only ever
	// match where the position already is the start of a line.
	line := func(token chroma.TokenType) []chroma.Rule {
		rules := []chroma.Rule{{Pattern: `\n`, Type: chroma.Text, Mutator: chroma.Pop(1)}}
		if dim != "" {
			rules = append(rules, chroma.Rule{Pattern: dim, Type: chroma.Comment})
		}

		return append(rules,
			chroma.Rule{Pattern: regexp.QuoteMeta(forcesReplacement), Type: chroma.GenericStrong},
			chroma.Rule{Pattern: `[^\n]`, Type: token},
		)
	}

	rules := chroma.Rules{"root": {}}
	for _, entry := range diffGlyphs {
		state := "glyph" + entry.glyph

		rules["root"] = append(rules["root"], chroma.Rule{
			Pattern: regexp.QuoteMeta(entry.glyph) + ` `,
			Type:    entry.token,
			Mutator: chroma.Push(state),
		})
		rules[state] = line(entry.token)
	}

	rules["root"] = append(rules["root"], chroma.Rule{Pattern: `\n`, Type: chroma.Text})
	if dim != "" {
		rules["root"] = append(rules["root"], chroma.Rule{Pattern: dim, Type: chroma.Comment})
	}

	// what is left carries no glyph of its own: the nested lines of a change the plan
	// shows one side of, which read as the line above them.
	rules["root"] = append(rules["root"],
		chroma.Rule{Pattern: `[ \t]+`, Type: chroma.Text, Mutator: chroma.Push("plain")},
		chroma.Rule{Pattern: `[^\n]`, Type: chroma.Text},
	)
	rules["plain"] = line(chroma.Text)

	// the name is the one chroma gives its own, which a lookup resolves before it ever
	// reaches the aliases: a lexer registered under any other name would be shadowed by
	// the one it means to replace.
	lexers.Register(chroma.MustNewLexer(
		&chroma.Config{
			Name:     "Diff",
			Aliases:  []string{"diff", "udiff"},
			EnsureNL: true,
		},
		func() chroma.Rules { return rules },
	))
}
