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
- `IsKeyword(TokenType) bool` (lexer.go:982) – true for the ABSTRACT..AT range (every entry of the keyword table; `TestIsKeywordCoversEveryWordToken` guards it). The old ABSTRACT..WHILE range missed `layout` and `at`, which were appended after `WHILE`. The Yul parser uses it to treat Solidity-only keywords as Yul identifiers, so keep any new keyword token inside the range
- `IsIdentifier(rune) bool` (lexer.go:959)
- `(TokenType) String() string` (lexer.go:311)

**Internal scanners:** `readNumber` (662, dec/frac/exp, underscores), `readHexNumber` (703, `0x…`), `readString` (724, escapes, `'`/`"`), `readIdentifier` (534, classifies typed keywords via `isIntType`/`isUintType`/`isBytesNType`/`isFixedNType`/`isUfixedNType`), `skipWhitespaceAndComments` (495, `//` and `/* */`), `readOperator` (772, longest-match: 3-char `>>>`/`>>=`/`<<=` → 2-char → 1-char).

**Prefixed string literals:** `readIdentifier` also emits `HEX_STRING` /
`UNICODE_STRING` when the scanned word is exactly `hex` or `unicode` **and** the
very next character is a quote (`hex"1900"`, `unicode"..."`). Solidity binds the
prefix to the immediately following quote, so `hex` / `unicode` stay ordinary
identifiers everywhere else. Both token types existed and [[builder]]
`parseStringLiteral` already branched on them, but nothing ever produced them:
the prefix lexed as a bare keyword and the quote became a separate `STRING`,
which `parsePrimary` rejected with "expected expression" and — in tolerant mode —
desynchronized the rest of the file. Underscore separators inside the body
(`hex"19_00"`, OpenZeppelin v5 `MessageHashUtils`) need no special handling since
the body is scanned as a string.

## Change checklist (new keyword/operator/literal)

1. Add a `TokenType` constant (lexer.go:14-170).
2. If a keyword: add to `keywords` (318). If it can also be an identifier, add it to `isContextualKeyword` in [[builder]] AND `expectMemberName` coverage.
3. Add a `tokenNames` entry (172) so error messages are readable.
4. New operator → extend `readOperator` (preserve longest-match order). New literal shape → extend the relevant `read*` scanner.
5. Add a case to `lexer_test.go` (it asserts token streams).

## Tests
`lexer_test.go` — token-stream assertions per construct.
