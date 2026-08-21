# pkg/parser — Public API

## Purpose

The package external consumers import (`github.com/th13vn/solast-go/pkg/parser`). Thin facade over [[builder]] + [[ast-index]]. This is the only API surface w3goaudit depends on — keep it backward-compatible.

## parser.go

**Options** (parser.go:14):
```go
type Options struct {
    Tolerant bool // collect & recover from errors instead of stopping
    Loc      bool // attach line/column
    Range    bool // attach byte offsets
}
```

**Errors:**
- `Error{Message string; Line, Column int}` (parser.go:36) — JSON-tagged.
- `ParserError{Errors []*Error}` (parser.go:24) — implements `error` (returns first message).

**Functions:**
- `Parse(input string, opts *Options) (*ast.SourceUnit, error)` (parser.go:47) — non-tolerant: first error is fatal; tolerant: returns the AST and **discards** recovered errors.
- `ParseWithErrors(input string, opts *Options) (*ast.SourceUnit, []*Error, error)` (parser.go:92) — like `Parse` but ALSO returns recovered errors in tolerant mode (empty slice = clean). Use this when silent truncation must be detectable. *w3goaudit's builder uses this to warn on incomplete extraction.* Added in v0.1.6.
- `ParseReader(io.Reader, *Options)` (parser.go:133), `ParseToJSON(input, *Options) ([]byte, error)` (parser.go:142, 2-space indent).
- `Visit(node, Visitor)` / `VisitSimple(node, *SimpleVisitor)` (parser.go:151) — wrap `ast.Walk`/`ast.WalkSimple`.
- Type aliases (parser.go:161): `Visitor`, `BaseVisitor`, `SimpleVisitor` re-exported from `ast`.

## Compatibility rules

- Never change an existing exported signature; ADD new functions (as `ParseWithErrors` was added).
- The JSON shape of `ParseToJSON` / serialized nodes is a consumer contract — additive changes only.

## Tests

- `parser_test.go` — broad construct coverage (the main suite).
- `struct_contextual_keyword_test.go` — regression for the contextual-keyword member desync (struct field / enum value named `from`) and `ParseWithErrors` surfacing tolerant errors.
