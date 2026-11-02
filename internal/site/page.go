package site

import (
	"fmt"
	"strings"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

// DocPage is one hand-written Markdown page to include in the site.
// Title is taken from its first level-1 heading; Body is the page
// content without that heading line.
type DocPage struct {
	File  string // file name within the docs directory, e.g. "guide.md"
	Title string
	Body  string
}

// packageFile is the output file name for a package page.
func packageFile(pkg *parser.Package) string {
	return pkg.Name + ".md"
}

// RenderPackage renders one package as a Markdown page under theme.
func RenderPackage(pkg *parser.Package, theme Theme) string {
	var b strings.Builder
	b.WriteString(heading(1, "package "+pkg.Name))
	if pkg.Doc != "" {
		b.WriteString(theme.docText(pkg.Doc) + "\n\n")
	}
	b.WriteString(heading(2, "Index"))
	for _, fn := range pkg.Functions {
		b.WriteString("- " + link("func "+fn.Name, "#"+anchor("func "+fn.Name)) + "\n")
	}
	for _, typ := range pkg.Types {
		b.WriteString("- " + link("type "+typ.Name, "#"+anchor("type "+typ.Name)) + "\n")
	}
	for _, val := range pkg.Consts {
		b.WriteString("- " + link("const "+val.Name, "#"+anchor("const "+val.Name)) + "\n")
	}
	for _, val := range pkg.Vars {
		b.WriteString("- " + link("var "+val.Name, "#"+anchor("var "+val.Name)) + "\n")
	}
	b.WriteString("\n")

	if len(pkg.Functions) > 0 {
		b.WriteString(heading(2, "Functions"))
		for _, fn := range pkg.Functions {
			writeFunction(&b, fn, theme)
		}
	}
	if len(pkg.Types) > 0 {
		b.WriteString(heading(2, "Types"))
		for _, typ := range pkg.Types {
			writeType(&b, typ, theme)
		}
	}
	if len(pkg.Consts) > 0 {
		b.WriteString(heading(2, "Constants"))
		for _, val := range pkg.Consts {
			writeValue(&b, "const", val, theme)
		}
	}
	if len(pkg.Vars) > 0 {
		b.WriteString(heading(2, "Variables"))
		for _, val := range pkg.Vars {
			writeValue(&b, "var", val, theme)
		}
	}
	return b.String()
}

func writeFunction(b *strings.Builder, fn parser.Function, theme Theme) {
	// Parser signatures omit the func keyword ("Add(a, b int) int"),
	// so rendering adds it back for a complete declaration.
	if theme.SectionPerSymbol {
		b.WriteString(heading(3, "func "+fn.Name))
		b.WriteString(codeBlock("go", "func "+fn.Signature))
		if fn.Doc != "" {
			b.WriteString(theme.docText(fn.Doc) + "\n\n")
		}
		return
	}
	line := "- `func " + fn.Signature + "`"
	if d := theme.docText(fn.Doc); d != "" {
		line += " — " + d
	}
	b.WriteString(line + "\n")
}

func writeType(b *strings.Builder, typ parser.Type, theme Theme) {
	if theme.SectionPerSymbol {
		b.WriteString(heading(3, "type "+typ.Name))
		b.WriteString(codeBlock("go", typ.Declaration))
		if typ.Doc != "" {
			b.WriteString(theme.docText(typ.Doc) + "\n\n")
		}
		if len(typ.Fields) > 0 {
			b.WriteString(heading(4, "Fields"))
			b.WriteString("| Name | Type | Doc |\n| --- | --- | --- |\n")
			for _, f := range typ.Fields {
				b.WriteString(fmt.Sprintf("| %s | `%s` | %s |\n", f.Name, f.Type, firstLine(f.Doc)))
			}
			b.WriteString("\n")
		}
		for _, m := range typ.Methods {
			b.WriteString("- method `func " + m.Signature + "`")
			if d := theme.docText(m.Doc); d != "" {
				b.WriteString(" — " + d)
			}
			b.WriteString("\n")
		}
		if len(typ.Methods) > 0 {
			b.WriteString("\n")
		}
		return
	}
	// Declarations can span lines (struct/interface bodies), which
	// would break a one-line bullet, so compact lists the type name.
	line := "- type `" + typ.Name + "`"
	if d := theme.docText(typ.Doc); d != "" {
		line += " — " + d
	}
	b.WriteString(line + "\n")
}

func writeValue(b *strings.Builder, kind string, val parser.Value, theme Theme) {
	if theme.SectionPerSymbol {
		b.WriteString(heading(3, kind+" "+val.Name))
		b.WriteString(codeBlock("go", val.Declaration))
		if val.Doc != "" {
			b.WriteString(theme.docText(val.Doc) + "\n\n")
		}
		return
	}
	line := "- " + kind + " `" + val.Declaration + "`"
	if d := theme.docText(val.Doc); d != "" {
		line += " — " + d
	}
	b.WriteString(line + "\n")
}

// RenderIndex renders the site index: package pages first, then any
// hand-written doc pages, plus a pointer to the search index.
func RenderIndex(pkg *parser.Package, docs []DocPage) string {
	var b strings.Builder
	b.WriteString(heading(1, "Documentation"))
	b.WriteString(heading(2, "Packages"))
	b.WriteString("- " + link("package "+pkg.Name, packageFile(pkg)) + "\n")
	if len(docs) > 0 {
		b.WriteString("\n" + heading(2, "Guides"))
		for _, d := range docs {
			b.WriteString("- " + link(d.Title, d.File) + "\n")
		}
	}
	b.WriteString("\nSearch the whole site with the generated `search.json` index.\n")
	return b.String()
}
