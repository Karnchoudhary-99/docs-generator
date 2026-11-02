// Command docsgen generates documentation from source comments and Markdown.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
	"github.com/Karnchoudhary-99/docs-generator/internal/site"
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
	case "site":
		err = runSite(os.Args[2:])
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
	fmt.Fprintln(os.Stderr, "       docsgen site <srcdir> <outdir> [--theme classic|compact] [--docs <dir>]")
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

// runSite generates a Markdown site for the Go package in srcdir
// into outdir, optionally including hand-written pages from docsDir.
func runSite(args []string) error {
	var srcDir, outDir, docsDir, themeName string
	positional := 0
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--theme":
			if i+1 >= len(args) {
				return fmt.Errorf("--theme requires a value")
			}
			i++
			themeName = args[i]
		case strings.HasPrefix(arg, "--theme="):
			themeName = strings.TrimPrefix(arg, "--theme=")
		case arg == "--docs":
			if i+1 >= len(args) {
				return fmt.Errorf("--docs requires a value")
			}
			i++
			docsDir = args[i]
		case strings.HasPrefix(arg, "--docs="):
			docsDir = strings.TrimPrefix(arg, "--docs=")
		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown flag %q", arg)
		case positional == 0:
			srcDir = arg
			positional++
		case positional == 1:
			outDir = arg
			positional++
		default:
			return fmt.Errorf("unexpected argument %q", arg)
		}
	}
	if srcDir == "" || outDir == "" {
		return fmt.Errorf("site requires <srcdir> and <outdir> arguments")
	}
	theme, err := site.ParseTheme(themeName)
	if err != nil {
		return err
	}
	if err := site.GenerateSite(srcDir, docsDir, outDir, theme); err != nil {
		return err
	}
	fmt.Printf("site written to %s\n", outDir)
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
		fmt.Printf("\n%s\n", typ.Declaration)
		printDoc(typ.Doc)
		for _, m := range typ.Methods {
			fmt.Printf("  method %s\n", m.Signature)
		}
	}
	for _, val := range pkg.Consts {
		fmt.Printf("\n%s\n", val.Declaration)
		printDoc(val.Doc)
	}
	for _, val := range pkg.Vars {
		fmt.Printf("\n%s\n", val.Declaration)
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
