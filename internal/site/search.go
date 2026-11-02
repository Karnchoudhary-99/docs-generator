package site

import (
	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

// Entry is one searchable item: a package, a symbol, or a doc page.
// URL is relative to the site root, e.g. "samplepkg.md#func-add".
type Entry struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
}

// BuildSearchIndex collects every searchable item for a package page
// plus the hand-written doc pages, in index-page order.
func BuildSearchIndex(pkg *parser.Package, docs []DocPage) []Entry {
	file := packageFile(pkg)
	entries := []Entry{{Title: "package " + pkg.Name, URL: file, Text: stripMarkdown(pkg.Doc)}}
	for _, fn := range pkg.Functions {
		entries = append(entries, Entry{
			Title: "func " + fn.Name, URL: file + "#" + anchor("func "+fn.Name),
			Text: stripMarkdown(fn.Doc),
		})
	}
	for _, typ := range pkg.Types {
		entries = append(entries, Entry{
			Title: "type " + typ.Name, URL: file + "#" + anchor("type "+typ.Name),
			Text: stripMarkdown(typ.Doc),
		})
		for _, m := range typ.Methods {
			entries = append(entries, Entry{
				Title: "method " + typ.Name + "." + m.Name,
				URL:   file + "#" + anchor("type "+typ.Name),
				Text:  stripMarkdown(m.Doc),
			})
		}
	}
	for _, val := range pkg.Consts {
		entries = append(entries, Entry{
			Title: "const " + val.Name, URL: file + "#" + anchor("const "+val.Name),
			Text: stripMarkdown(val.Doc),
		})
	}
	for _, val := range pkg.Vars {
		entries = append(entries, Entry{
			Title: "var " + val.Name, URL: file + "#" + anchor("var "+val.Name),
			Text: stripMarkdown(val.Doc),
		})
	}
	for _, d := range docs {
		entries = append(entries, Entry{Title: d.Title, URL: d.File, Text: stripMarkdown(d.Body)})
	}
	return entries
}
