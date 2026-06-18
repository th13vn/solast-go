# grammar — Reference ANTLR Grammar (NOT the runtime parser)

## Purpose & critical caveat

`SolidityLexer.g4` and `SolidityParser.g4` are the **official Solidity ANTLR grammar**, kept as a **specification reference only**. They are **NOT compiled into the runtime parser** — solast-go's parser is the hand-written recursive-descent code in [[builder]] + [[lexer]]. There is no `internal/gen` in the build; `make generate` (ANTLR → Go) exists for experimentation but its output is unused.

> Practical consequence: **updating the `.g4` files does NOT change parsing behavior.** They tell you WHAT the language allows; you must then implement it by hand in the lexer/builder.

## Files

- `SolidityLexer.g4` — token definitions (keywords, literals, operators) — mirror in [[lexer]].
- `SolidityParser.g4` — production rules (contract/function/statement/expression grammar) — mirror in [[builder]].

## Workflow

- `make update-grammar` — pulls the latest `.g4` from `ethereum/solidity@develop/docs/grammar/`. Run this when targeting a new Solidity release to see what syntax changed, then implement the delta by hand.
- `make generate` / `scripts/generate.sh` — runs ANTLR (needs Java) to emit Go into `internal/gen/` for reference; not part of the build.

## Using the grammar to add syntax

1. `make update-grammar`, then `git diff grammar/` to see new productions/tokens.
2. For each new token → add to [[lexer]] (`TokenType`, `keywords`, `tokenNames`).
3. For each new production → add a parse function in [[builder]] and wire it into the right dispatch.
4. Add the AST node + visitor wiring in [[ast-index]].
5. Add tests in [[parser-index]] + a fixture under `testdata/contracts/`.

See the `solast-go-update` skill for the full playbook.
