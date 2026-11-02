package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

// LoadDocs reads every .md file in dir (non-recursive, sorted) as a
// DocPage. A page's title is its first level-1 heading; a file without
// one falls back to its file name. An empty directory yields no pages.
func LoadDocs(dir string) ([]DocPage, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var docs []DocPage
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		docs = append(docs, parseDocPage(e.Name(), string(raw)))
	}
	sort.Slice(docs, func(i, j int) bool { return docs[i].File < docs[j].File })
	return docs, nil
}

// parseDocPage splits raw Markdown into a title (first "# " line)
// and the remaining body.
func parseDocPage(file, raw string) DocPage {
	page := DocPage{File: file, Title: strings.TrimSuffix(file, ".md")}
	var body []string
	for _, line := range strings.Split(raw, "\n") {
		if page.Title == strings.TrimSuffix(file, ".md") && strings.HasPrefix(line, "# ") {
			page.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			continue
		}
		body = append(body, line)
	}
	page.Body = strings.Trim(strings.Join(body, "\n"), "\n")
	return page
}

// Generate writes a complete site for pkg into outDir: index.md, one
// page per package, search.json, and verbatim copies of the doc pages.
// outDir is created if missing; existing site files are overwritten.
func Generate(pkg *parser.Package, docs []DocPage, theme Theme, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	write := func(name, content string) error {
		return os.WriteFile(filepath.Join(outDir, name), []byte(content), 0o644)
	}
	if err := write("index.md", RenderIndex(pkg, docs)); err != nil {
		return err
	}
	if err := write(packageFile(pkg), RenderPackage(pkg, theme)); err != nil {
		return err
	}
	for _, d := range docs {
		if err := write(d.File, "# "+d.Title+"\n\n"+d.Body+"\n"); err != nil {
			return err
		}
	}
	search, err := json.MarshalIndent(BuildSearchIndex(pkg, docs), "", "  ")
	if err != nil {
		return err
	}
	if err := write("search.json", string(search)+"\n"); err != nil {
		return err
	}
	return nil
}

// GenerateSite is the one-call path used by the CLI: parse srcDir,
// optionally load docsDir, then generate into outDir.
func GenerateSite(srcDir, docsDir, outDir string, theme Theme) error {
	pkg, err := parser.ParseDir(srcDir)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", srcDir, err)
	}
	var docs []DocPage
	if docsDir != "" {
		docs, err = LoadDocs(docsDir)
		if err != nil {
			return fmt.Errorf("loading docs from %s: %w", docsDir, err)
		}
	}
	return Generate(pkg, docs, theme, outDir)
}
