package site

import "fmt"

// Theme controls how much detail a generated page carries. The classic
// theme gives every symbol its own section with the full doc comment;
// the compact theme lists symbols as one-liners with only the first
// line of each doc comment, for large packages.
type Theme struct {
	Name             string
	FullDocs         bool // include the whole doc comment, not just its first line
	SectionPerSymbol bool // one section per symbol instead of a bullet list
}

// Classic is the default, detailed theme.
var Classic = Theme{Name: "classic", FullDocs: true, SectionPerSymbol: true}

// Compact is the terse theme for large packages.
var Compact = Theme{Name: "compact", FullDocs: false, SectionPerSymbol: false}

// ParseTheme resolves a theme name, or returns an error naming the
// valid themes so a typo fails loudly instead of silently defaulting.
func ParseTheme(name string) (Theme, error) {
	switch name {
	case "", "classic":
		return Classic, nil
	case "compact":
		return Compact, nil
	default:
		return Theme{}, fmt.Errorf("unknown theme %q (valid themes: classic, compact)", name)
	}
}

// docText returns the doc text a theme shows: full text for classic,
// first line only for compact.
func (t Theme) docText(doc string) string {
	if t.FullDocs {
		return doc
	}
	return firstLine(doc)
}
