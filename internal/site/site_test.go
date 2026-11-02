package site

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

func fixturePkg(t *testing.T) *parser.Package {
	t.Helper()
	pkg, err := parser.ParseDir("../parser/testdata/samplepkg")
	if err != nil {
		t.Fatalf("ParseDir fixture: %v", err)
	}
	return pkg
}

func TestRenderPackageClassic(t *testing.T) {
	page := RenderPackage(fixturePkg(t), Classic)
	for _, want := range []string{
		"# package samplepkg",
		"### func Add",
		"```go\nfunc Add(a, b int) int\n```",
		"returns the sum of a and b",
		"### type Counter",
		"| Count | `int` |",
		"method `func (*Counter) Inc()`",
		"### const DefaultLimit",
		"### var Version",
	} {
		if !strings.Contains(page, want) {
			t.Errorf("classic page missing %q", want)
		}
	}
	if strings.Contains(page, "helper") {
		t.Errorf("classic page leaks unexported symbol")
	}
}

func TestRenderPackageCompact(t *testing.T) {
	page := RenderPackage(fixturePkg(t), Compact)
	if strings.Contains(page, "### func Add") {
		t.Errorf("compact page should not use per-symbol sections")
	}
	if !strings.Contains(page, "- `func Add(a, b int) int` — Add returns the sum of a and b.") {
		t.Errorf("compact page missing one-line function entry:\n%s", page)
	}
}

func TestRenderIndexLinks(t *testing.T) {
	docs := []DocPage{{File: "guide.md", Title: "User Guide", Body: "How to use it."}}
	index := RenderIndex(fixturePkg(t), docs)
	for _, want := range []string{"[package samplepkg](samplepkg.md)", "[User Guide](guide.md)", "search.json"} {
		if !strings.Contains(index, want) {
			t.Errorf("index missing %q", want)
		}
	}
}

func TestParseTheme(t *testing.T) {
	if th, err := ParseTheme(""); err != nil || th.Name != "classic" {
		t.Errorf("empty theme should default to classic, got %v, %v", th, err)
	}
	if th, err := ParseTheme("compact"); err != nil || th.Name != "compact" {
		t.Errorf("compact theme: %v, %v", th, err)
	}
	if _, err := ParseTheme("neon"); err == nil || !strings.Contains(err.Error(), "classic, compact") {
		t.Errorf("unknown theme should name valid themes, got %v", err)
	}
}

func TestAnchor(t *testing.T) {
	cases := map[string]string{"func Add": "func-add", "type Counter": "type-counter", "HTTP Server!": "http-server"}
	for in, want := range cases {
		if got := anchor(in); got != want {
			t.Errorf("anchor(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSearchIndex(t *testing.T) {
	entries := BuildSearchIndex(fixturePkg(t), []DocPage{{File: "guide.md", Title: "User Guide", Body: "How to **use** it."}})
	byTitle := map[string]Entry{}
	for _, e := range entries {
		byTitle[e.Title] = e
	}
	fn, ok := byTitle["func Add"]
	if !ok || fn.URL != "samplepkg.md#func-add" || !strings.Contains(fn.Text, "sum of a and b") {
		t.Errorf("func Add entry wrong: %+v (present=%v)", fn, ok)
	}
	if _, ok := byTitle["method Counter.Inc"]; !ok {
		t.Errorf("search index missing method entry")
	}
	guide, ok := byTitle["User Guide"]
	if !ok || strings.Contains(guide.Text, "**") {
		t.Errorf("doc entry should have stripped markdown: %+v (present=%v)", guide, ok)
	}
}

func TestLoadDocs(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "b.md"), []byte("# Second\nBody two.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("No heading here.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("ignored"), 0o644); err != nil {
		t.Fatal(err)
	}
	docs, err := LoadDocs(dir)
	if err != nil {
		t.Fatalf("LoadDocs: %v", err)
	}
	if len(docs) != 2 || docs[0].File != "a.md" || docs[0].Title != "a" || docs[1].Title != "Second" {
		t.Fatalf("unexpected docs: %+v", docs)
	}
	if docs[1].Body != "Body two." {
		t.Errorf("body = %q, want %q", docs[1].Body, "Body two.")
	}
}

func TestGenerateWritesSite(t *testing.T) {
	out := t.TempDir()
	docs := []DocPage{{File: "guide.md", Title: "User Guide", Body: "How to use it."}}
	if err := Generate(fixturePkg(t), docs, Classic, out); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	for _, name := range []string{"index.md", "samplepkg.md", "guide.md", "search.json"} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Errorf("missing site file %s: %v", name, err)
		}
	}
	raw, err := os.ReadFile(filepath.Join(out, "search.json"))
	if err != nil {
		t.Fatal(err)
	}
	var entries []Entry
	if err := json.Unmarshal(raw, &entries); err != nil {
		t.Fatalf("search.json is not valid JSON: %v", err)
	}
	if len(entries) < 5 {
		t.Errorf("search.json has only %d entries", len(entries))
	}
	guide, err := os.ReadFile(filepath.Join(out, "guide.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(guide), "# User Guide\n") {
		t.Errorf("copied doc page lost its heading: %q", string(guide))
	}
}
