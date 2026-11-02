package site

import (
	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

// Summary describes what a site build produced, for the CLI's build
// report. Files lists every file Generate writes, in write order.
type Summary struct {
	Package   string
	Functions int
	Types     int
	Methods   int
	Consts    int
	Vars      int
	Guides    int
	Files     []string
}

// Summarize counts the symbols in pkg and the output files a build
// for pkg plus docs would write.
func Summarize(pkg *parser.Package, docs []DocPage) Summary {
	s := Summary{Package: pkg.Name, Guides: len(docs)}
	s.Functions = len(pkg.Functions)
	s.Types = len(pkg.Types)
	s.Consts = len(pkg.Consts)
	s.Vars = len(pkg.Vars)
	for _, typ := range pkg.Types {
		s.Methods += len(typ.Methods)
	}
	s.Files = append(s.Files, "index.md", packageFile(pkg))
	for _, d := range docs {
		s.Files = append(s.Files, d.File)
	}
	s.Files = append(s.Files, "search.json")
	return s
}
