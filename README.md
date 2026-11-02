# Docs Generator CLI

Generate beautiful docs from code comments and Markdown — in Go.

## Roadmap

- [x] Project scaffold (README roadmap, license, gitignore)
- [x] Go doc-comment parser (packages, functions, types)
- [ ] Markdown site generator (index, search, themes)
- [ ] CLI (generate, serve, watch) + tests
- [ ] Self-review checklist, then public

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

## License

MIT — see [LICENSE](LICENSE).
