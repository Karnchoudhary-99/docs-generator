package parser

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

// nodeString renders an AST node back to Go source text, trimming the
// trailing newline printer adds. It is used for signatures and type
// expressions, never for whole function bodies.
func nodeString(fset *token.FileSet, node any) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}

// funcSignature renders "Name(params) results" for a FuncDecl, prefixed
// with "(Receiver) " when the function is a method. The body is never
// included: only the FuncType is printed.
func funcSignature(fset *token.FileSet, decl *ast.FuncDecl) string {
	var b strings.Builder
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		b.WriteString("(")
		b.WriteString(nodeString(fset, decl.Recv.List[0].Type))
		b.WriteString(") ")
	}
	b.WriteString(decl.Name.Name)
	b.WriteString(nodeString(fset, decl.Type))
	return b.String()
}

// receiverName returns the bare type name of a method receiver, dropping
// any pointer star and any package qualifier.
func receiverName(fset *token.FileSet, decl *ast.FuncDecl) string {
	if decl.Recv == nil || len(decl.Recv.List) == 0 {
		return ""
	}
	typ := nodeString(fset, decl.Recv.List[0].Type)
	typ = strings.TrimPrefix(typ, "*")
	if i := strings.LastIndex(typ, "."); i >= 0 {
		typ = typ[i+1:]
	}
	return typ
}

// structFields extracts the fields of a struct type, or nil for
// non-struct types. Anonymous fields keep their type name as the name.
func structFields(fset *token.FileSet, spec *ast.TypeSpec) []Field {
	st, ok := spec.Type.(*ast.StructType)
	if !ok || st.Fields == nil {
		return nil
	}
	var fields []Field
	for _, f := range st.Fields.List {
		typ := nodeString(fset, f.Type)
		if len(f.Names) == 0 {
			fields = append(fields, Field{Name: typ, Type: typ, Doc: docText(f.Doc)})
			continue
		}
		for _, name := range f.Names {
			fields = append(fields, Field{Name: name.Name, Type: typ, Doc: docText(f.Doc)})
		}
	}
	return fields
}

// typeKind classifies a type declaration: struct, interface, alias
// (a declaration with `=`), or anything else.
func typeKind(spec *ast.TypeSpec) string {
	switch spec.Type.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	default:
		if spec.Assign.IsValid() {
			return "alias"
		}
		return "other"
	}
}

// docText returns the comment text of a comment group, or "" when the
// group is nil. go/doc already strips comment markers; this helper is
// for raw AST comment groups such as struct fields.
func docText(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	return strings.TrimSpace(group.Text())
}
