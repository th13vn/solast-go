# internal/lexer — Tokenizer

## Purpose

Converts Solidity source text into a flat `[]Token` stream consumed by `internal/builder`. Hand-written scanner (no ANTLR). This is the FIRST place to change when a new Solidity version introduces a new keyword, operator, or literal form.

## Key File

### lexer.go (~962 lines)

**Token** (lexer.go:391) — one lexed unit:
```go
type Token struct {
    Type   TokenType // classification (enum below)
    Value  string    // raw text (underscores stripped from numbers)
    Line   int       // 1-indexed
    Column int       // 0-indexed
    Start  int       // byte offset
    End    int       // byte offset + length
}
```

**TokenType** (lexer.go:14-170) — `int` iota enum, grouped:
- Special: `EOF`, `ILLEGAL`, `COMMENT` (lexer.go:15)
- Literals: `IDENTIFIER`, `NUMBER`, `HEX_NUMBER`, `STRING`, `HEX_STRING`, `UNICODE_STRING` (lexer.go:20)
- Keywords (~69): control flow, visibility, mutability, contract kinds, members, storage, type modifiers (lexer.go:28)
- **Contextual keywords**: `FROM`, `GLOBAL`, `REVERT`, `ERROR`, `TRANSIENT`, `LAYOUT`, `AT`, `UNICODE`, `HEX`, `LET` — keyword tokens that are ALSO legal identifiers in some positions (struct/enum members, params, var names). Mishandling these desyncs the parser — see [[builder]] `expectMemberName`.
- Typed keywords: `INT`, `UINT`, `BYTE`, `BYTES_N`, `FIXED_N`, `UFIXED_N` (lexer.go:102) — `uint256`/`bytes32`/`fixedMxN` are classified by suffix at scan time, not stored as one token per width.
- Operators & punctuation: assignment (13), comparison (6), logical (3), bitwise (7), arithmetic (6), unary (2), brackets/delimiters (lexer.go:110).

**keywords map** (lexer.go:318) — `map[string]TokenType` (lowercase keyword → type). Add an entry here for any new keyword.

**tokenNames map + `String()`** (lexer.go:172, 311) — `TokenType` → human text (used in parser error messages). Add a name for every new token type.

**Exported API:**
- `New(input string) *Lexer` (lexer.go:418)
- `(*Lexer) NextToken() Token` (lexer.go:428) — skips whitespace/comments, dispatches by first rune
- `(*Lexer) Tokenize() []Token` (lexer.go:925) — full stream
- `IsKeyword(TokenType) bool` (lexer.go:954) — true for the ABSTRACT..WHILE range
- `IsIdentifier(rune) bool` (lexer.go:959)
- `(TokenType) String() string` (lexer.go:311)

**Internal scanners:** `readNumber` (662, dec/frac/exp, underscores), `readHexNumber` (703, `0x…`), `readString` (724, escapes, `'`/`"`), `readIdentifier` (534, classifies typed keywords via `isIntType`/`isUintType`/`isBytesNType`/`isFixedNType`/`isUfixedNType`), `skipWhitespaceAndComments` (495, `//` and `/* */`), `readOperator` (772, longest-match: 3-char `>>>`/`>>=`/`<<=` → 2-char → 1-char).

## Change checklist (new keyword/operator/literal)

1. Add a `TokenType` constant (lexer.go:14-170).
2. If a keyword: add to `keywords` (318). If it can also be an identifier, add it to `isContextualKeyword` in [[builder]] AND `expectMemberName` coverage.
3. Add a `tokenNames` entry (172) so error messages are readable.
4. New operator → extend `readOperator` (preserve longest-match order). New literal shape → extend the relevant `read*` scanner.
5. Add a case to `lexer_test.go` (it asserts token streams).

## Tests
`lexer_test.go` — token-stream assertions per construct.
