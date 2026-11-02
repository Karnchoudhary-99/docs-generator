package parser_test

import (
	"strings"
	"testing"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

func mustParse(t *testing.T) *parser.Package {
	t.Helper()
	pkg, err := parser.ParseDir("testdata/samplepkg")
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}
	return pkg
}

func TestParseDirPackageAndFunctions(t *testing.T) {
	pkg := mustParse(t)
	if pkg.Name != "samplepkg" {
		t.Fatalf("Name = %q, want samplepkg", pkg.Name)
	}
	if !strings.Contains(pkg.Doc, "fixture for the docs-generator parser tests") {
		t.Fatalf("Doc = %q, want fixture description", pkg.Doc)
	}
	if len(pkg.Functions) != 1 {
		t.Fatalf("Functions = %d, want 1 (unexported helper must be dropped)", len(pkg.Functions))
	}
	fn := pkg.Functions[0]
	if fn.Name != "Add" || fn.Doc != "Add returns the sum of a and b." {
		t.Fatalf("Function = %+v, want Add with its doc", fn)
	}
	if fn.Signature != "Add(a, b int) int" {
		t.Fatalf("Signature = %q, want Add(a, b int) int", fn.Signature)
	}
	if fn.Receiver != "" {
		t.Fatalf("Receiver = %q, want empty for a package function", fn.Receiver)
	}
}

func TestParseDirTypesAndMethods(t *testing.T) {
	pkg := mustParse(t)
	byName := map[string]parser.Type{}
	for _, typ := range pkg.Types {
		byName[typ.Name] = typ
	}
	counter, ok := byName["Counter"]
	if !ok {
		t.Fatalf("Types = %v, want Counter", pkg.Types)
	}
	if counter.Kind != "struct" {
		t.Fatalf("Counter.Kind = %q, want struct", counter.Kind)
	}
	// go/doc drops unexported struct fields, so only Count survives.
	if len(counter.Fields) != 1 {
		t.Fatalf("Counter.Fields = %+v, want only exported Count", counter.Fields)
	}
	if counter.Fields[0].Name != "Count" || counter.Fields[0].Type != "int" {
		t.Fatalf("Fields[0] = %+v, want Count int", counter.Fields[0])
	}
	methods := map[string]parser.Function{}
	for _, m := range counter.Methods {
		methods[m.Name] = m
	}
	inc, ok := methods["Inc"]
	if !ok {
		t.Fatalf("Counter.Methods = %v, want Inc", counter.Methods)
	}
	if inc.Receiver != "Counter" || !strings.HasPrefix(inc.Signature, "(*Counter) Inc()") {
		t.Fatalf("Inc = %+v, want (*Counter) Inc() receiver", inc)
	}
	if _, ok := methods["Value"]; !ok {
		t.Fatalf("Counter.Methods = %v, want Value", counter.Methods)
	}
	if greeter, ok := byName["Greeter"]; !ok || greeter.Kind != "interface" {
		t.Fatalf("Greeter = %+v, want interface kind", greeter)
	}
	if list, ok := byName["StringList"]; !ok || list.Kind != "other" {
		t.Fatalf("StringList = %+v, want other kind", list)
	}
}

func TestParseDirConstsAndVars(t *testing.T) {
	pkg := mustParse(t)
	if len(pkg.Consts) != 1 || pkg.Consts[0].Name != "DefaultLimit" {
		t.Fatalf("Consts = %+v, want DefaultLimit", pkg.Consts)
	}
	if !strings.Contains(pkg.Consts[0].Declaration, "DefaultLimit = 10") {
		t.Fatalf("Const declaration = %q", pkg.Consts[0].Declaration)
	}
	if len(pkg.Vars) != 1 || pkg.Vars[0].Name != "Version" {
		t.Fatalf("Vars = %+v, want Version", pkg.Vars)
	}
}

func TestParseDirErrors(t *testing.T) {
	if _, err := parser.ParseDir("testdata/does-not-exist"); err == nil {
		t.Fatal("ParseDir on missing dir: want error, got nil")
	}
	if _, err := parser.ParseDir(t.TempDir()); err == nil {
		t.Fatal("ParseDir on empty dir: want error, got nil")
	}
}
