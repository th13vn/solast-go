# Architecture

How solast-go turns Solidity source into an AST, and where each responsibility lives. For per-package detail with file:line references, read each package's `INDEX.md`.

## The pipeline

```
source string
   │  internal/lexer.New(src).Tokenize()
   ▼
[]lexer.Token                       ← flat token stream (keywords, literals, operators, positions)
   │  internal/builder.New(src, opts).Build()
   ▼
*ast.SourceUnit                     ← tree of pkg/ast nodes
   │  pkg/parser.Parse / ParseWithErrors / ParseToJSON
   ▼
consumers (cmd/solast CLI, w3goaudit, …)
```

Three stages, three packages:

1. **Lexer** (`internal/lexer`) — character scan → tokens. Knows keywords, contextual keywords, typed keywords (`uint256`/`bytes32`/`fixedMxN` classified by suffix), number/hex/string forms, comments, and longest-match operators.
2. **Builder** (`internal/builder`) — recursive-descent parser over the token stream. A function per grammar construct, an operator-precedence ladder for expressions, and tolerant error recovery. **This is the authoritative grammar implementation.**
3. **AST** (`pkg/ast`) — the output node structs (`Node` interface + `BaseNode`), plus `Walk`/`WalkSimple` visitors. `pkg/parser` is a thin public facade.

## The grammar is a reference, not the engine

`grammar/SolidityLexer.g4` / `SolidityParser.g4` are the official ethereum/solidity ANTLR grammar, kept for reference. `make generate` can emit a Go ANTLR parser into `internal/gen/`, but **that output is not compiled into solast-go** — nothing imports it. The real parser is hand-written.

So the workflow for new language features is: use the grammar to learn *what* changed (`make update-grammar` + `git diff grammar/`), then implement the delta by hand in lexer + builder + ast. The grammar keeps you honest about syntax; it does not generate code you ship.

## Tolerant parsing (and its trap)

`pkg/parser.Options.Tolerant` makes the builder collect errors and recover instead of failing — essential for tools (w3goaudit) that must analyze a project even when one file is imperfect.

Mechanics (`internal/builder/helpers.go`):
- `expect(t)` on a match advances; on a mismatch it records an error. **In tolerant mode it does NOT advance** — recovery is delegated to `synchronize()`, which skips to the next `;` or top-level keyword.
- Because a mismatch leaves the cursor in place, a token the parser fails to anticipate can cascade: the parser can walk past a closing `}` and consume the rest of an enclosing contract, dropping its functions/state **with no hard error**.

This is exactly what happened with a struct field named `from` (a contextual keyword) before v0.1.5. Two defenses exist now:
- `expectMemberName()` accepts identifiers **and** contextual keywords for declaration names (struct members, enum values), so `struct S { address from; }` parses correctly.
- `ParseWithErrors()` returns the recovered errors that tolerant `Parse` discards, so callers can warn when extraction may be incomplete (w3goaudit's builder does this).

**Rule of thumb:** any place that reads a declaration NAME must use `expectMemberName()`, not bare `expect(IDENTIFIER)`; and tooling should call `ParseWithErrors` so silent truncation is visible.

## Contextual keywords

`from`, `error`, `revert`, `global`, `transient`, `layout`, `at` are keyword tokens (the lexer emits e.g. `FROM`) that are *also* legal identifiers in Solidity. The builder permits them as names via `isContextualKeyword()` (parameters, variables, primary expressions) and `expectMemberName()` (struct/enum members). When adding a new contextual keyword, update both the lexer keyword map and these helpers.

## Source positions

`Options.Loc` / `Options.Range` attach `Location` (line/column) and `Range` (byte offsets) to nodes via `setLocation()` in the builder. `setLocation` has a per-node-type switch — a new AST node must get a case or it will lack positions (which downstream tools like w3goaudit rely on for finding locations).

## Where to make a change

| Change | Edit |
|--------|------|
| New keyword/operator/literal | `internal/lexer` (TokenType, `keywords`, `tokenNames`, scanner) |
| New statement/expression/type/definition | `internal/builder` (parse fn + dispatch) + `pkg/ast` (node + visitor) + `setLocation` |
| New contextual keyword | lexer keyword map + `isContextualKeyword` + `expectMemberName` |
| New public API | `pkg/parser` (additive only) |
| New CLI capability | `cmd/solast` |

Full step-by-step lives in `docs/development.md` ("Adding Support for New Solidity Features") and the `solast-go-update` skill.
