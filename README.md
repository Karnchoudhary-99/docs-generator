# Docs Generator CLI

Generate beautiful docs from code comments and Markdown — in Go.

## Roadmap

- [x] Project scaffold (README roadmap, license, gitignore)
- [x] Go doc-comment parser (packages, functions, types)
- [x] Markdown site generator (index, search, themes)
- [x] CLI (generate, serve, watch) + tests
- [x] Self-review checklist, then public

## Go doc-comment parser

`internal/parser` parses a directory of Go files with the standard
library's `go/parser` and `go/doc` packages, so it follows exactly the
same rules as `go doc`: only exported symbols are reported, and
unexported declarations (and unexported struct fields) are dropped.
Test files are ignored, and a directory must contain exactly one package.

```sh
# Human-readable summary
go run ./cmd/docsgen parse ./internal/parser/testdata/samplepkg

# Machine-readable output for the site generator
go run ./cmd/docsgen parse ./internal/parser/testdata/samplepkg --json
```

## Markdown site generator

`internal/site` turns a parsed package into a static Markdown site —
no dependencies, no build step beyond `docsgen` itself:

```sh
# index.md + <package>.md + search.json into ./site
go run ./cmd/docsgen site ./internal/parser/testdata/samplepkg ./site

# Include hand-written Markdown pages, and use the terse theme
go run ./cmd/docsgen site ./internal/parser/testdata/samplepkg ./site \
  --docs ./examples/docs --theme compact
```

Output per site:

- `index.md` — package page links, guide links, search pointer.
- `<package>.md` — index of symbols, then one section per symbol
  (classic theme) or one-line entries (compact theme), with
  signatures in Go code blocks, struct field tables, and methods.
- `search.json` — one entry per package, symbol, method, and guide
  page (`title`, `url`, `text` with Markdown stripped), ready for
  offline client-side search.
- Guide pages from `--docs` are copied verbatim, linked from the
  index, and added to the search index.

A generated example lives in [`examples/sample-site`](examples/sample-site/).

## CLI: generate, serve, watch

`generate` is the one-command build: it parses, generates, and prints
a build summary (symbol counts and the files written). `--dry-run`
reports without writing anything.

```sh
# Full build with a summary
go run ./cmd/docsgen generate ./internal/parser/testdata/samplepkg ./site \
  --docs ./examples/docs

# Preview what a build would produce, without writing files
go run ./cmd/docsgen generate ./internal/parser/testdata/samplepkg ./site --dry-run
```

`serve` previews an already-generated site directory over HTTP.
Directories resolve to their `index.md`; directories without one
return 404 instead of a listing. Every request is logged.

```sh
go run ./cmd/docsgen serve ./site --port 8080
# serving ./site at http://127.0.0.1:8080
```

`watch` builds once, then rebuilds whenever a `.go` or `.md` file
under the source or docs directory changes. A failed rebuild is
reported and watching continues. With `--port` it also serves the
output directory while rebuilding, so the preview stays fresh.

```sh
go run ./cmd/docsgen watch ./internal/parser/testdata/samplepkg ./site \
  --docs ./examples/docs --interval 1s --port 8080
```

## License

MIT — see [LICENSE](LICENSE).
