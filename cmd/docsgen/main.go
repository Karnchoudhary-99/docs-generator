// Command docsgen generates documentation from source comments and Markdown.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "parse":
		err = runParse(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "docsgen: unknown command %q\n", os.Args[1])
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "docsgen: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: docsgen parse <dir> [--json]")
}

// runParse parses the Go package in dir and prints a summary, or JSON
// with --json.
func runParse(args []string) error {
	var dir string
	var asJSON bool
	for _, arg := range args {
		switch {
		case arg == "--json":
			asJSON = true
		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown flag %q", arg)
		case dir == "":
			dir = arg
		default:
			return fmt.Errorf("unexpected argument %q", arg)
		}
	}
	if dir == "" {
		return fmt.Errorf("parse requires a directory argument")
	}
	pkg, err := parser.ParseDir(dir)
	if err != nil {
		return err
	}
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(pkg)
	}
	printPackage(pkg)
	return nil
}

func printPackage(pkg *parser.Package) {
	fmt.Printf("package %s\n", pkg.Name)
	if pkg.Doc != "" {
		fmt.Printf("%s\n", firstLine(pkg.Doc))
	}
	for _, fn := range pkg.Functions {
		fmt.Printf("\nfunc %s\n", fn.Signature)
		printDoc(fn.Doc)
	}
	for _, typ := range pkg.Types {
		fmt.Printf("\ntype %s\n", typ.Declaration)
		printDoc(typ.Doc)
		for _, m := range typ.Methods {
			fmt.Printf("  method %s\n", m.Signature)
		}
	}
	for _, val := range pkg.Consts {
		fmt.Printf("\nconst %s\n", val.Declaration)
		printDoc(val.Doc)
	}
	for _, val := range pkg.Vars {
		fmt.Printf("\nvar %s\n", val.Declaration)
		printDoc(val.Doc)
	}
}

func printDoc(doc string) {
	if doc != "" {
		fmt.Printf("  %s\n", firstLine(doc))
	}
}

func firstLine(s string) string {
	if i := strings.Index(s, "\n"); i >= 0 {
		return s[:i]
	}
	return s
}
