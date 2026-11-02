package site

import (
	"testing"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

func TestSummarizeCountsSymbolsAndFiles(t *testing.T) {
	pkg := &parser.Package{
		Name:      "samplepkg",
		Functions: []parser.Function{{Name: "Add"}, {Name: "Sub"}},
		Types: []parser.Type{
			{Name: "Counter", Methods: []parser.Function{{Name: "Inc"}, {Name: "Value"}}},
		},
		Consts: []parser.Value{{Name: "DefaultLimit"}},
		Vars:   []parser.Value{{Name: "Version"}},
	}
	docs := []DocPage{{File: "guide.md", Title: "Guide"}}
	s := Summarize(pkg, docs)
	if s.Functions != 2 || s.Types != 1 || s.Methods != 2 || s.Consts != 1 || s.Vars != 1 || s.Guides != 1 {
		t.Fatalf("summary = %+v, want 2 funcs, 1 type, 2 methods, 1 const, 1 var, 1 guide", s)
	}
	want := []string{"index.md", "samplepkg.md", "guide.md", "search.json"}
	if len(s.Files) != len(want) {
		t.Fatalf("files = %v, want %v", s.Files, want)
	}
	for i, f := range want {
		if s.Files[i] != f {
			t.Fatalf("files = %v, want %v", s.Files, want)
		}
	}
}
