# solast-go — Repository Index (AI-readable)

A hand-written Solidity parser in Go: source text → tokens → AST. Consumed by w3goaudit. This index maps the codebase for agents; each package has its own `INDEX.md` with file:line detail.

## Pipeline

```
source .sol
  └─ internal/lexer   →  []Token              (tokenizer)
       └─ internal/builder → *ast.SourceUnit  (recursive-descent parser)
            └─ pkg/ast        (node structs + visitor)
                 └─ pkg/parser  (public API: Parse / ParseWithErrors / ParseToJSON)
                      └─ cmd/solast (CLI), external consumers (w3goaudit)
```

## Packages

| Path | Role | INDEX |
|------|------|-------|
| `internal/lexer` | Tokenizer (keywords, literals, operators) | [internal/lexer/INDEX.md](internal/lexer/INDEX.md) |
| `internal/builder` | Recursive-descent parser (authoritative) | [internal/builder/INDEX.md](internal/builder/INDEX.md) |
| `pkg/ast` | AST node types + visitor walkers | [pkg/ast/INDEX.md](pkg/ast/INDEX.md) |
| `pkg/parser` | Public API (import this) | [pkg/parser/INDEX.md](pkg/parser/INDEX.md) |
| `pkg/version` | Solidity version/pragma detection | [pkg/version/INDEX.md](pkg/version/INDEX.md) |
| `cmd/solast` | CLI (parse/validate/version-detect) | [cmd/solast/INDEX.md](cmd/solast/INDEX.md) |
| `grammar` | Reference ANTLR `.g4` (NOT runtime) | [grammar/INDEX.md](grammar/INDEX.md) |
| `scripts` | `generate.sh` (ANTLR, reference) | [scripts/INDEX.md](scripts/INDEX.md) |

## Key facts for agents

1. **The parser is hand-written.** `grammar/*.g4` + `make generate` are reference only; editing them does not change behavior. New syntax = edit lexer + builder + ast by hand.
2. **Tolerant mode is the default for tooling.** `expect()` does NOT advance on mismatch in tolerant mode; recovery is via `synchronize()`. A mis-handled token can desync and silently drop the rest of a contract — always use `expectMemberName()` for declaration names, and prefer `ParseWithErrors` to detect silent truncation.
3. **Contextual keywords** (`from`, `error`, `revert`, `global`, `transient`, `layout`, `at`) are valid identifiers in many positions — handle them in lexer + `isContextualKeyword`/`expectMemberName`.
4. **Public API is a contract.** `pkg/parser` and the JSON node shape are depended on by w3goaudit — additive changes only.
5. Module version = `VERSION` file + matching `vX.Y.Z` git tag (distinct from `pkg/version`, which is Solidity versions).

## Docs

- [docs/overview.md](docs/overview.md) — what & why
- [docs/architecture.md](docs/architecture.md) — pipeline internals (AI-oriented)
- [docs/getting-started.md](docs/getting-started.md) — install & first parse
- [docs/development.md](docs/development.md) — dev setup, build, "adding a feature" walkthrough
- [docs/contributing.md](docs/contributing.md) — contribution flow

To update grammar / add Solidity syntax, use the **`solast-go-update`** skill.
