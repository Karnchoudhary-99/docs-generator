package parser

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"sort"
	"strings"
)

// ParseDir parses the Go package in dir and returns its documentation.
// Test files are ignored. Only exported symbols are reported, matching
// `go doc` behaviour; unexported declarations are dropped by go/doc.
// The directory must contain exactly one package.
func ParseDir(dir string) (*Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("parser: no Go package found in %s", dir)
	}
	if len(pkgs) > 1 {
		names := make([]string, 0, len(pkgs))
		for name := range pkgs {
			names = append(names, name)
		}
		sort.Strings(names)
		return nil, fmt.Errorf("parser: multiple packages in %s: %s", dir, strings.Join(names, ", "))
	}
	var astPkg *ast.Package
	for _, pkg := range pkgs {
		astPkg = pkg
	}
	docPkg := doc.New(astPkg, "", 0)
	out := &Package{
		Name: docPkg.Name,
		Doc:  strings.TrimSpace(docPkg.Doc),
	}
	for _, fn := range docPkg.Funcs {
		out.Functions = append(out.Functions, convertFunc(fset, fn))
	}
	for _, typ := range docPkg.Types {
		out.Types = append(out.Types, convertType(fset, typ))
	}
	for _, val := range docPkg.Consts {
		out.Consts = append(out.Consts, convertValues(fset, val)...)
	}
	for _, val := range docPkg.Vars {
		out.Vars = append(out.Vars, convertValues(fset, val)...)
	}
	// go/doc already sorts, but sort explicitly so output is stable
	// regardless of the go/doc version's internal ordering.
	sort.Slice(out.Functions, func(i, j int) bool { return out.Functions[i].Name < out.Functions[j].Name })
	sort.Slice(out.Types, func(i, j int) bool { return out.Types[i].Name < out.Types[j].Name })
	sort.Slice(out.Consts, func(i, j int) bool { return out.Consts[i].Name < out.Consts[j].Name })
	sort.Slice(out.Vars, func(i, j int) bool { return out.Vars[i].Name < out.Vars[j].Name })
	return out, nil
}

func convertFunc(fset *token.FileSet, fn *doc.Func) Function {
	return Function{
		Name:      fn.Name,
		Doc:       strings.TrimSpace(fn.Doc),
		Signature: funcSignature(fset, fn.Decl),
		Receiver:  receiverName(fset, fn.Decl),
	}
}

func convertType(fset *token.FileSet, typ *doc.Type) Type {
	out := Type{
		Name: typ.Name,
		Doc:  strings.TrimSpace(typ.Doc),
	}
	if spec := findTypeSpec(typ); spec != nil {
		out.Kind = typeKind(spec)
		out.Declaration = "type " + nodeString(fset, spec)
		out.Fields = structFields(fset, spec)
	}
	for _, m := range typ.Methods {
		out.Methods = append(out.Methods, convertFunc(fset, m))
	}
	sort.Slice(out.Methods, func(i, j int) bool { return out.Methods[i].Name < out.Methods[j].Name })
	return out
}

// findTypeSpec returns the TypeSpec for a doc.Type from its GenDecl.
func findTypeSpec(typ *doc.Type) *ast.TypeSpec {
	if typ.Decl == nil {
		return nil
	}
	for _, spec := range typ.Decl.Specs {
		if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == typ.Name {
			return ts
		}
	}
	return nil
}

func convertValues(fset *token.FileSet, val *doc.Value) []Value {
	decl := nodeString(fset, val.Decl)
	out := make([]Value, 0, len(val.Names))
	for _, name := range val.Names {
		out = append(out, Value{
			Name:        name,
			Doc:         strings.TrimSpace(val.Doc),
			Declaration: decl,
		})
	}
	return out
}
