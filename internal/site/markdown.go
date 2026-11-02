// Package site turns parsed Go packages and hand-written Markdown pages
// into a small static documentation site: one Markdown page per package,
// an index page, and a JSON search index for client-side search.
package site

import (
	"strings"
)

// heading renders a Markdown heading of the given level (1-6).
func heading(level int, text string) string {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	return strings.Repeat("#", level) + " " + text + "\n"
}

// codeBlock renders a fenced code block. The fence is lengthened past
// any backtick run inside the code so the block cannot break early.
func codeBlock(lang, code string) string {
	fence := "```"
	for strings.Contains(code, fence) {
		fence += "`"
	}
	return fence + lang + "\n" + strings.TrimSuffix(code, "\n") + "\n" + fence + "\n"
}

// link renders a Markdown link, escaping brackets in the label.
func link(label, target string) string {
	label = strings.ReplaceAll(label, "[", "\\[")
	label = strings.ReplaceAll(label, "]", "\\]")
	return "[" + label + "](" + target + ")"
}

// anchor converts a heading or symbol name into a GitHub-style anchor
// slug: lowercased, runs of non-alphanumerics collapsed to one dash.
func anchor(s string) string {
	var b strings.Builder
	lastDash := true // trims leading dashes
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// stripMarkdown removes the most common Markdown markers so search
// text is plain words rather than formatting syntax.
func stripMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"#", "", "*", "", "_", "", "`", "", ">", "",
		"[", "", "]", "", "(", " ", ")", " ",
	)
	fields := strings.Fields(replacer.Replace(s))
	return strings.Join(fields, " ")
}

// firstLine returns the first line of s, or "" for empty input.
func firstLine(s string) string {
	if i := strings.Index(s, "\n"); i >= 0 {
		return s[:i]
	}
	return s
}
