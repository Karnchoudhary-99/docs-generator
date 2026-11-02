// Package parser extracts documentation from Go source comments.
//
// It parses a directory of Go files with the standard library's
// go/parser and go/doc packages, so it understands exactly the same
// documentation rules as `go doc`: only exported symbols are reported,
// and doc comments are the comments attached to each declaration.
package parser

// Package is the parsed documentation of one Go package.
type Package struct {
	Name      string     `json:"name"`
	Doc       string     `json:"doc,omitempty"`
	Functions []Function `json:"functions,omitempty"`
	Types     []Type     `json:"types,omitempty"`
	Consts    []Value    `json:"consts,omitempty"`
	Vars      []Value    `json:"vars,omitempty"`
}

// Function is a package-level function or a method on a type.
// For a method, Receiver holds the receiver type name (for example
// "Counter"); it is empty for package-level functions.
type Function struct {
	Name      string `json:"name"`
	Doc       string `json:"doc,omitempty"`
	Signature string `json:"signature"`
	Receiver  string `json:"receiver,omitempty"`
}

// Type is a named type declaration.
type Type struct {
	Name        string     `json:"name"`
	Doc         string     `json:"doc,omitempty"`
	Kind        string     `json:"kind"` // "struct", "interface", "alias", or "other"
	Declaration string     `json:"declaration"`
	Fields      []Field    `json:"fields,omitempty"`
	Methods     []Function `json:"methods,omitempty"`
}

// Field is one field of a struct type.
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Doc  string `json:"doc,omitempty"`
}

// Value is a single constant or variable name from a const/var block.
type Value struct {
	Name        string `json:"name"`
	Doc         string `json:"doc,omitempty"`
	Declaration string `json:"declaration"`
}
