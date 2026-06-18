# internal/builder — Recursive-Descent Parser

## Purpose

Turns the `[]Token` stream from [[lexer]] into an `*ast.SourceUnit`. This is the **authoritative parser** — the ANTLR `.g4` files in `grammar/` are reference only and are NOT used at runtime (see [[grammar-index]]). Adding new Solidity syntax means editing the files here.

## Builder core (builder.go)

```go
type Builder struct {
    tokens  []lexer.Token
    pos     int
    errors  []*Error
    options *Options
}
type Options struct { Tolerant, Loc, Range bool } // builder.go:32
type Error   struct { Message string; Line, Column int } // builder.go:13
```

- `New(input string, opts *Options) *Builder` (builder.go:39) — tokenizes immediately.
- `(*Builder) Build() (*ast.SourceUnit, error)` (builder.go:56) — top loop over `parseSourceUnitElement`.
- `(*Builder) Errors() []*Error` (builder.go:89) — recovered errors (surfaced to callers via [[parser-index]] `ParseWithErrors`).

**Dispatch tables (the map of "keyword → parse function"):**
- `parseSourceUnitElement` (builder.go:93) — pragma / import / contract|interface|library|abstract / struct / enum / function / event / error / using / type / file-level const.
- `parseContractBodyElement` (builder.go:309) — function / constructor / modifier / fallback / receive / struct / enum / event / error / using / type / state-variable.

## Files (by construct)

| File | Lines | Parses |
|------|------:|--------|
| builder.go | ~574 | entry, dispatch, contract/function/modifier/constructor/fallback/receive, pragma, import, inheritance |
| expressions.go | ~692 | the precedence ladder + primary expressions, calls, literals |
| statements.go | ~848 | blocks, if/for/while/do, return/emit/revert, try/catch, **assembly (Yul)**, unchecked, var-decls, tuple-decls |
| types.go | ~576 | type names, mappings, function types, arrays, struct/enum/event/error/using/UDVT definitions, params, state vars |
| helpers.go | ~296 | token navigation, error recovery, contextual-keyword handling, `setLocation` |

## Token navigation & recovery (helpers.go) — READ BEFORE EDITING

- `peek` (12) / `previous` (19) / `advance` (26) / `check(t)` (33) / `isAtEnd` (40).
- `expect(t)` (44): on match advance+return; on mismatch `addError`. **TOLERANT MODE: does NOT advance on mismatch** (lets `synchronize` recover). Non-tolerant: advances to avoid infinite loops. *This is the trap behind the historical `from`-field desync bug.*
- `synchronize` (84): skips to the next `;` or top-level keyword after an error.
- `isContextualKeyword()` (121): `FROM|ERROR|REVERT|GLOBAL|TRANSIENT|LAYOUT|AT` — keywords usable as identifiers.
- `expectMemberName()` (136): identifier **or** contextual keyword; **use this for every declaration NAME** (struct members types.go:353, enum values types.go:388) instead of bare `expect(IDENTIFIER)`, or a member named `from` desyncs the parser and silently drops the rest of the contract.
- `setLocation(node, start, end)` (150): fills `Loc`/`Range` when enabled; has a per-node-type switch — **add a case for every new AST node** or it won't get source positions.

## Expression precedence ladder (expressions.go) — lowest → highest

`parseExpression`(27) → `parseAssignment`(31) → `parseTernary`(49) → `parseLogicalOr`(69) → `parseLogicalAnd`(86) → `parseEquality`(103) → `parseRelational`(120) → `parseBitwiseOr`(137) → `parseBitwiseXor`(154) → `parseBitwiseAnd`(171) → `parseShift`(188) → `parseAdditive`(205) → `parseMultiplicative`(222) → `parseExponentiation`(239, right-assoc) → `parseUnary`(256) → `parsePostfix`(273) → `parseCallMemberIndex`(289) → `parsePrimary`(409). A new binary operator slots into the level matching its precedence; a new primary form (literal/keyword-expr) goes in `parsePrimary`.

## Change checklist (new statement / type / definition)

1. Add the AST node + `NodeType` in [[ast-index]] and a `setLocation` case (helpers.go).
2. Add the parse function in the matching file; wire it into the right dispatch (`parseSourceUnitElement` / `parseContractBodyElement` / `parseStatement` statements.go:28 / `parseTypeName` types.go:10).
3. Use `expectMemberName()` for any declaration name that could be a contextual keyword.
4. Add a `Walk`/`WalkSimple` case + `SimpleVisitor` callback in [[ast-index]].
5. Add tests in [[parser-index]].
