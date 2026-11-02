// Command docsgen generates documentation from source comments and Markdown.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Karnchoudhary-99/docs-generator/internal/parser"
	"github.com/Karnchoudhary-99/docs-generator/internal/serve"
	"github.com/Karnchoudhary-99/docs-generator/internal/site"
	"github.com/Karnchoudhary-99/docs-generator/internal/watch"
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
	case "generate":
		err = runGenerate(os.Args[2:])
	case "serve":
		err = runServe(os.Args[2:])
	case "watch":
		err = runWatch(os.Args[2:])
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
	fmt.Fprintln(os.Stderr, "       docsgen generate <srcdir> <outdir> [--theme classic|compact] [--docs <dir>] [--dry-run]")
	fmt.Fprintln(os.Stderr, "       docsgen serve <dir> [--addr 127.0.0.1] [--port 8080]")
	fmt.Fprintln(os.Stderr, "       docsgen watch <srcdir> <outdir> [--theme classic|compact] [--docs <dir>] [--interval 1s] [--port 0]")
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

// parseGenerateArgs parses the flags shared by generate and watch:
// --theme, --docs, plus positional <srcdir> <outdir>. Any other flag
// is returned in rest for the caller to reject.
func parseGenerateArgs(args []string) (srcDir, outDir, docsDir, themeName string, rest []string, err error) {
	positional := 0
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--theme":
			if i+1 >= len(args) {
				return "", "", "", "", nil, fmt.Errorf("--theme requires a value")
			}
			i++
			themeName = args[i]
		case strings.HasPrefix(arg, "--theme="):
			themeName = strings.TrimPrefix(arg, "--theme=")
		case arg == "--docs":
			if i+1 >= len(args) {
				return "", "", "", "", nil, fmt.Errorf("--docs requires a value")
			}
			i++
			docsDir = args[i]
		case strings.HasPrefix(arg, "--docs="):
			docsDir = strings.TrimPrefix(arg, "--docs=")
		case strings.HasPrefix(arg, "-"):
			rest = append(rest, arg)
		case positional == 0:
			srcDir = arg
			positional++
		case positional == 1:
			outDir = arg
			positional++
		default:
			return "", "", "", "", nil, fmt.Errorf("unexpected argument %q", arg)
		}
	}
	if srcDir == "" || outDir == "" {
		return "", "", "", "", nil, fmt.Errorf("requires <srcdir> and <outdir> arguments")
	}
	return srcDir, outDir, docsDir, themeName, rest, nil
}

// loadForGenerate parses srcDir and optionally loads docsDir, for
// generate and watch.
func loadForGenerate(srcDir, docsDir string) (*parser.Package, []site.DocPage, error) {
	pkg, err := parser.ParseDir(srcDir)
	if err != nil {
		return nil, nil, err
	}
	var docs []site.DocPage
	if docsDir != "" {
		docs, err = site.LoadDocs(docsDir)
		if err != nil {
			return nil, nil, fmt.Errorf("loading docs from %s: %w", docsDir, err)
		}
	}
	return pkg, docs, nil
}

func printSummary(pkg *parser.Package, docs []site.DocPage, outDir string) {
	s := site.Summarize(pkg, docs)
	fmt.Printf("package %s: %d functions, %d types (%d methods), %d consts, %d vars, %d guides\n",
		s.Package, s.Functions, s.Types, s.Methods, s.Consts, s.Vars, s.Guides)
	fmt.Printf("files written to %s:\n", outDir)
	for _, f := range s.Files {
		fmt.Printf("  %s\n", f)
	}
}

// runGenerate builds a site like site, then prints a build summary.
// With --dry-run it parses and reports without writing anything.
func runGenerate(args []string) error {
	var dryRun bool
	var plain []string
	for _, arg := range args {
		if arg == "--dry-run" {
			dryRun = true
			continue
		}
		plain = append(plain, arg)
	}
	srcDir, outDir, docsDir, themeName, rest, err := parseGenerateArgs(plain)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("unknown flag %q", rest[0])
	}
	theme, err := site.ParseTheme(themeName)
	if err != nil {
		return err
	}
	pkg, docs, err := loadForGenerate(srcDir, docsDir)
	if err != nil {
		return err
	}
	if dryRun {
		printSummary(pkg, docs, outDir)
		fmt.Println("(dry run: nothing written)")
		return nil
	}
	if err := site.Generate(pkg, docs, theme, outDir); err != nil {
		return err
	}
	printSummary(pkg, docs, outDir)
	return nil
}

// runServe previews the generated site in dir on --addr:--port until
// interrupted, shutting down gracefully on SIGINT/SIGTERM.
func runServe(args []string) error {
	var dir, addr string
	port := 8080
	addr = "127.0.0.1"
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--port":
			if i+1 >= len(args) {
				return fmt.Errorf("--port requires a value")
			}
			i++
			p, err := strconv.Atoi(args[i])
			if err != nil {
				return fmt.Errorf("invalid --port %q", args[i])
			}
			port = p
		case strings.HasPrefix(arg, "--port="):
			p, err := strconv.Atoi(strings.TrimPrefix(arg, "--port="))
			if err != nil {
				return fmt.Errorf("invalid --port %q", strings.TrimPrefix(arg, "--port="))
			}
			port = p
		case arg == "--addr":
			if i+1 >= len(args) {
				return fmt.Errorf("--addr requires a value")
			}
			i++
			addr = args[i]
		case strings.HasPrefix(arg, "--addr="):
			addr = strings.TrimPrefix(arg, "--addr=")
		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown flag %q", arg)
		case dir == "":
			dir = arg
		default:
			return fmt.Errorf("unexpected argument %q", arg)
		}
	}
	if dir == "" {
		return fmt.Errorf("serve requires a directory argument")
	}
	if _, err := os.Stat(dir); err != nil {
		return err
	}
	listen := fmt.Sprintf("%s:%d", addr, port)
	srv := serve.NewServer(listen, dir, os.Stdout)
	fmt.Printf("serving %s at http://%s\n", dir, listen)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdown)
	}
}

// runWatch builds the site once, then rebuilds it whenever a .go or
// .md file under srcdir (or docsdir) changes. A failed rebuild is
// reported and watching continues. With --port nonzero it also serves
// outdir for preview while rebuilding.
func runWatch(args []string) error {
	var plain []string
	interval := time.Second
	port := 0
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--port":
			if i+1 >= len(args) {
				return fmt.Errorf("--port requires a value")
			}
			i++
			p, err := strconv.Atoi(args[i])
			if err != nil {
				return fmt.Errorf("invalid --port %q", args[i])
			}
			port = p
		case strings.HasPrefix(arg, "--port="):
			p, err := strconv.Atoi(strings.TrimPrefix(arg, "--port="))
			if err != nil {
				return fmt.Errorf("invalid --port %q", strings.TrimPrefix(arg, "--port="))
			}
			port = p
		case arg == "--interval":
			if i+1 >= len(args) {
				return fmt.Errorf("--interval requires a value")
			}
			i++
			d, err := time.ParseDuration(args[i])
			if err != nil {
				return fmt.Errorf("invalid --interval %q", args[i])
			}
			interval = d
		case strings.HasPrefix(arg, "--interval="):
			d, err := time.ParseDuration(strings.TrimPrefix(arg, "--interval="))
			if err != nil {
				return fmt.Errorf("invalid --interval %q", strings.TrimPrefix(arg, "--interval="))
			}
			interval = d
		default:
			plain = append(plain, arg)
		}
	}
	srcDir, outDir, docsDir, themeName, rest, err := parseGenerateArgs(plain)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("unknown flag %q", rest[0])
	}
	theme, err := site.ParseTheme(themeName)
	if err != nil {
		return err
	}
	rebuild := func() {
		pkg, docs, err := loadForGenerate(srcDir, docsDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "watch: rebuild failed: %v\n", err)
			return
		}
		if err := site.Generate(pkg, docs, theme, outDir); err != nil {
			fmt.Fprintf(os.Stderr, "watch: rebuild failed: %v\n", err)
			return
		}
		fmt.Printf("rebuilt %s at %s\n", outDir, time.Now().Format(time.TimeOnly))
	}
	rebuild()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if port != 0 {
		srv := serve.NewServer(fmt.Sprintf("127.0.0.1:%d", port), outDir, os.Stdout)
		fmt.Printf("serving %s at http://127.0.0.1:%d\n", outDir, port)
		go func() {
			_ = srv.ListenAndServe()
		}()
		defer func() {
			shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(shutdown)
		}()
	}
	dirs := []string{srcDir}
	if docsDir != "" {
		dirs = append(dirs, docsDir)
	}
	fmt.Printf("watching %s for changes (interval %s)\n", strings.Join(dirs, ", "), interval)
	if err := watch.Run(ctx, dirs, interval, rebuild); err != nil && err != context.Canceled {
		return err
	}
	return nil
}
