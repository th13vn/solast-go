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

- `peek` (12) / `previous` (19) / `advance` (26) / `check(t)` (33) / `nextTokenIs(t)` (44, one token of lookahead, consumes nothing) / `isAtEnd` (51).
- `expect(t)` (55): on match advance+return; on mismatch `addError`. **TOLERANT MODE: does NOT advance on mismatch** (lets `synchronize` recover). Non-tolerant: advances to avoid infinite loops. *This is the trap behind the historical `from`-field desync bug.*
- `synchronize` (95): skips to the next `;` or top-level keyword after an error.
- `isContextualKeyword()` (132): `FROM|ERROR|REVERT|GLOBAL|TRANSIENT|LAYOUT|AT` — keywords usable as identifiers.
- `expectMemberName()` (147): identifier **or** contextual keyword; **use this for every declaration NAME** (struct members types.go:353, enum values types.go:388) instead of bare `expect(IDENTIFIER)`, or a member named `from` desyncs the parser and silently drops the rest of the contract.
- `isYulIdentifier()` / `expectYulIdentifier()`: inside assembly, any keyword
  token that is not a Yul keyword (`let if for switch case default function
  break continue true false hex`) is an ordinary Yul identifier. Use these, not
  `check(IDENTIFIER)`, for every Yul name, operand, call head and path member.
- `looksLikeVariableDeclaration()` in `statements.go` must accept contextual
  keywords both in dotted type paths and as the local declaration name. Since
  v0.1.11 this matches `parseVariableDeclaration`, so `UserInfo storage from =
  ...` cannot fall through to expression parsing and shred later functions in
  tolerant mode. Since v0.1.12 a leading elementary type followed by `(` is a
  conversion, never a declaration (`nextTokenIs(LPAREN)`), so
  `address(x).call(b);` stays an expression statement instead of being dropped.
- `setLocation(node, start, end)` (161): fills `Loc`/`Range` when enabled; has a per-node-type switch — **add a case for every new AST node** or it won't get source positions.

## Inline assembly (statements.go)

`parseAssemblyStatement` (423) consumes, in order: the `assembly` keyword, an
optional **dialect** string (`assembly "evmasm" { … }` → `InlineAssembly.Language`),
an optional **flags** group (`assembly ("memory-safe") { … }` →
`InlineAssembly.Flags`, via `parseAssemblyFlags`), then the block.

Both optional parts are the same desync trap as `expectMemberName`: the flags
group is optional to *write*, not optional to *consume*. Leaving its `(` for
`parseAssemblyBlock` makes `expect(LBRACE)` fail, and in tolerant mode that does
not advance — the parser then shreds the remainder of the file.
`assembly ("memory-safe")` is standard from Solidity 0.8.13 and pervasive in
OpenZeppelin v5 / Solady, so this silently blanked whole library files for
downstream consumers. Any future optional pre-block syntax must be consumed here
too.

## Yul identifiers spelled like Solidity keywords (statements.go)

The lexer tokenizes assembly with the Solidity keyword table, so the EVM
builtins `address`, `return` and `revert` (and any Solidity-only keyword used
as a Yul name, such as `from` or `error`) arrive as keyword tokens. The Yul
parser used to accept only `IDENTIFIER`:

- `mstore(0x14, address())` read `address` as a literal, failed on the `(`, and
  desynchronized the rest of the file (solady `ERC20.sol`, euler `Dispatch.sol`).
- `revert(p, 0x44)` in statement position was skipped token by token until an
  identifier argument started a bogus statement (v4-core `CustomRevert.sol`).
- with literal-only arguments, `return(0, 0x20)` / `revert(0, 0)` vanished
  with NO recovered error: silent AST loss in every OpenZeppelin `Address.sol`.

`parseAssemblyStatement_`, `parseAssemblyExpression`,
`parseAssemblyExpressionOrAssignment`, `parseAssemblyLocalDefinition`,
`parseAssemblyFunctionDefinition` and `parseAssemblyPathSuffix` all go through
`isYulIdentifier` / `expectYulIdentifier`. The same rule covers the external
function pointer members `g.address` / `g.selector`.

Yul function return variables follow `->`, which the lexer emits as
`RIGHT_ARROW`; the parser used to test for `ARROW` (`=>`), so every
`function f(a) -> b { ... }` desynchronized the rest of the file.

Known gap, deliberately untouched: Yul `:=` still lexes as `COLON` + `ASSIGN`,
so `let x := e` and `x := e` come out as a bare definition/identifier followed
by a separate expression statement. w3goaudit works around this by rewriting
the colon to a space inside assembly before parsing
(`normalizeYulAssignmentsForParser`).

## Yul paths (statements.go)

Inside assembly, only dot-free identifiers can be **declared**, but a path may
**refer** to a declaration outside the block: calldata slice members
(`sig.offset`, `sig.length`) and storage-pointer members (`x.slot`, `x.offset`).
Grammar: `yulPath: (YulIdentifier|YulEVMBuiltin) (YulPeriod (YulIdentifier|YulEVMBuiltin))*`.

`parseAssemblyPathSuffix` consumes the `.member` segments after the head
identifier and returns the full dotted text in the single `Name` field (no new
node type, so the public AST shape stays compatible). It is called from **both**
`parseAssemblyExpression` (operand position) and
`parseAssemblyExpressionOrAssignment` (assignment targets). Leaving the `.`
unconsumed desynchronized the block; in assignment position it silently produced
a truncated name instead.

## Declarations vs expressions at statement start (statements.go)

- `looksLikeVariableDeclaration` skips the `.member` type path BEFORE array
  dimensions, matching the grammar (`Storage.Assimilator[] memory x`). The old
  order met `[` where it expected a location or name and sent the declaration
  down the expression path, which died on `memory`.
- The for-loop init clause uses the same dispatch as a statement
  (`looksLikeVariableDeclaration`, then tuple, then expression). It used
  `isTypeName()`, which is true for ANY identifier, so `for (y >>= 1; ...)`
  failed on `>>=` and `for (i = 0; ...)` was recorded as a nameless
  declaration of type `i` (now an assignment `ExpressionStatement`).
- `parseTupleVariableDeclarationOrExpression` backtracks to BEFORE the `(`
  when the tuple is not a declaration, and drops the errors recorded while
  speculating. Restoring to after the `(` made
  `(items[i].ok, items[i].data) = g();` fail on the first comma.
- Known gap, deliberately untouched: a tuple assignment whose components all
  parse as type names, e.g. `(a, b) = g();` or `(s.x, t) = g();`, is still
  returned as a `VariableDeclarationStatement` with nameless variables whose
  `typeName` is the assignment target. Changing that alters the AST of every
  such assignment for consumers and needs a coordinated w3goaudit change.

## Definitions (types.go)

- `using { ... } for T` list entries are identifier paths
  (`parseUsingFunctionPath`), kept as dotted text in `Functions`
  (`Casting.intoSD59x18`, prb-math `ValueType.sol`). Reading one `IDENTIFIER`
  died on the `.` and dropped every later using directive in the file.
- Custom-error parameter names accept contextual keywords like event and
  function parameters do (`error DeploymentFailed(bytes error);`, Axelar).

## Binary-operation locations (expressions.go)

Every level of the precedence ladder captures `startTok := b.peek()` before parsing
its left operand and calls `setLocation(node, startTok, b.previous())` on the
`BinaryOperation` it builds, so the span covers the whole expression.

This was missing entirely: `setLocation` already had a `*ast.BinaryOperation` case,
but no ladder level ever called it, so **every** binary and comparison expression had
a nil `Loc`/`Range`. Consumers that anchor on those nodes lost all position
information — w3goaudit reported divide-before-multiply, boolean-cst, incorrect-exp,
and unchecked-arithmetic findings against the enclosing FUNCTION's line instead of the
operation's. When adding a new operator level, capture the start token the same way.

## Expression precedence ladder (expressions.go) — lowest → highest

`parseExpression`(27) → `parseAssignment`(31) → `parseTernary`(49) → `parseLogicalOr`(69) → `parseLogicalAnd`(86) → `parseEquality`(103) → `parseRelational`(120) → `parseBitwiseOr`(137) → `parseBitwiseXor`(154) → `parseBitwiseAnd`(171) → `parseShift`(188) → `parseAdditive`(205) → `parseMultiplicative`(222) → `parseExponentiation`(239, right-assoc) → `parseUnary`(256) → `parsePostfix`(273) → `parseCallMemberIndex`(289) → `parsePrimary`(409). A new binary operator slots into the level matching its precedence; a new primary form (literal/keyword-expr) goes in `parsePrimary`.

## Change checklist (new statement / type / definition)

1. Add the AST node + `NodeType` in [[ast-index]] and a `setLocation` case (helpers.go).
2. Add the parse function in the matching file; wire it into the right dispatch (`parseSourceUnitElement` / `parseContractBodyElement` / `parseStatement` statements.go:28 / `parseTypeName` types.go:10).
3. Use `expectMemberName()` for any declaration name that could be a contextual keyword.
4. Add a `Walk`/`WalkSimple` case + `SimpleVisitor` callback in [[ast-index]].
5. Add tests in [[parser-index]].
