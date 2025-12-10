// Package lexer provides a hand-written lexer for Solidity code.
// This lexer can be used standalone or swapped with ANTLR-generated lexer.
package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	EOF TokenType = iota
	ILLEGAL
	COMMENT

	// Literals
	IDENTIFIER
	NUMBER
	HEX_NUMBER
	STRING
	HEX_STRING
	UNICODE_STRING

	// Keywords
	ABSTRACT
	ADDRESS
	ANONYMOUS
	AS
	ASSEMBLY
	BOOL
	BREAK
	BYTES
	CALLDATA
	CASE
	CATCH
	CONSTANT
	CONSTRUCTOR
	CONTINUE
	CONTRACT
	DEFAULT
	DELETE
	DO
	ELSE
	EMIT
	ENUM
	ERROR
	EVENT
	EXTERNAL
	FALLBACK
	FALSE
	FIXED
	FOR
	FROM
	FUNCTION
	GLOBAL
	HEX
	IF
	IMMUTABLE
	IMPORT
	INDEXED
	INTERFACE
	INTERNAL
	IS
	LET
	LIBRARY
	MAPPING
	MEMORY
	MODIFIER
	NEW
	OVERRIDE
	PAYABLE
	PRAGMA
	PRIVATE
	PUBLIC
	PURE
	RECEIVE
	RETURN
	RETURNS
	REVERT
	STORAGE
	STRING_TYPE
	STRUCT
	SWITCH
	TRANSIENT
	TRUE
	TRY
	TYPE
	UFIXED
	UNCHECKED
	UNICODE
	USING
	VIEW
	VIRTUAL
	WHILE
	LAYOUT
	AT

	// Types
	INT
	UINT
	BYTE
	BYTES_N
	FIXED_N
	UFIXED_N

	// Operators
	LPAREN    // (
	RPAREN    // )
	LBRACK    // [
	RBRACK    // ]
	LBRACE    // {
	RBRACE    // }
	COLON     // :
	SEMICOLON // ;
	PERIOD    // .
	COMMA     // ,
	QUESTION    // ?
	ARROW       // =>
	RIGHT_ARROW // ->

	// Assignment
	ASSIGN       // =
	ASSIGN_ADD   // +=
	ASSIGN_SUB   // -=
	ASSIGN_MUL   // *=
	ASSIGN_DIV   // /=
	ASSIGN_MOD   // %=
	ASSIGN_AND   // &=
	ASSIGN_OR    // |=
	ASSIGN_XOR   // ^=
	ASSIGN_SHL   // <<=
	ASSIGN_SHR   // >>=
	ASSIGN_SAR   // >>>=

	// Comparison
	EQ  // ==
	NEQ // !=
	LT  // <
	GT  // >
	LTE // <=
	GTE // >=

	// Logical
	AND // &&
	OR  // ||
	NOT // !

	// Bitwise
	BIT_AND // &
	BIT_OR  // |
	BIT_XOR // ^
	BIT_NOT // ~
	SHL     // <<
	SHR     // >>
	SAR     // >>>

	// Arithmetic
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %
	EXP // **
	INC // ++
	DEC // --
)

var tokenNames = map[TokenType]string{
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	COMMENT:   "COMMENT",
	IDENTIFIER: "IDENTIFIER",
	NUMBER:    "NUMBER",
	HEX_NUMBER: "HEX_NUMBER",
	STRING:    "STRING",
	HEX_STRING: "HEX_STRING",
	UNICODE_STRING: "UNICODE_STRING",
	ABSTRACT:  "abstract",
	ADDRESS:   "address",
	ANONYMOUS: "anonymous",
	AS:        "as",
	ASSEMBLY:  "assembly",
	BOOL:      "bool",
	BREAK:     "break",
	BYTES:     "bytes",
	CALLDATA:  "calldata",
	CASE:      "case",
	CATCH:     "catch",
	CONSTANT:  "constant",
	CONSTRUCTOR: "constructor",
	CONTINUE:  "continue",
	CONTRACT:  "contract",
	DEFAULT:   "default",
	DELETE:    "delete",
	DO:        "do",
	ELSE:      "else",
	EMIT:      "emit",
	ENUM:      "enum",
	ERROR:     "error",
	EVENT:     "event",
	EXTERNAL:  "external",
	FALLBACK:  "fallback",
	FALSE:     "false",
	FIXED:     "fixed",
	FOR:       "for",
	FROM:      "from",
	FUNCTION:  "function",
	GLOBAL:    "global",
	HEX:       "hex",
	IF:        "if",
	IMMUTABLE: "immutable",
	IMPORT:    "import",
	INDEXED:   "indexed",
	INTERFACE: "interface",
	INTERNAL:  "internal",
	IS:        "is",
	LET:       "let",
	LIBRARY:   "library",
	MAPPING:   "mapping",
	MEMORY:    "memory",
	MODIFIER:  "modifier",
	NEW:       "new",
	OVERRIDE:  "override",
	PAYABLE:   "payable",
	PRAGMA:    "pragma",
	PRIVATE:   "private",
	PUBLIC:    "public",
	PURE:      "pure",
	RECEIVE:   "receive",
	RETURN:    "return",
	RETURNS:   "returns",
	REVERT:    "revert",
	STORAGE:   "storage",
	STRING_TYPE: "string",
	STRUCT:    "struct",
	SWITCH:    "switch",
	TRUE:      "true",
	TRY:       "try",
	TRANSIENT: "transient",
	TYPE:      "type",
	UFIXED:    "ufixed",
	UNCHECKED: "unchecked",
	LAYOUT:    "layout",
	AT:        "at",
	UNICODE:   "unicode",
	USING:     "using",
	VIEW:      "view",
	VIRTUAL:   "virtual",
	WHILE:     "while",
	INT:       "int",
	UINT:      "uint",
	BYTE:      "byte",
	BYTES_N:   "bytes<N>",
	FIXED_N:   "fixed<M>x<N>",
	UFIXED_N:  "ufixed<M>x<N>",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACK:    "[",
	RBRACK:    "]",
	LBRACE:    "{",
	RBRACE:    "}",
	COLON:     ":",
	SEMICOLON: ";",
	PERIOD:    ".",
	COMMA:     ",",
	QUESTION:    "?",
	ARROW:       "=>",
	RIGHT_ARROW: "->",
	ASSIGN:    "=",
	ASSIGN_ADD: "+=",
	ASSIGN_SUB: "-=",
	ASSIGN_MUL: "*=",
	ASSIGN_DIV: "/=",
	ASSIGN_MOD: "%=",
	ASSIGN_AND: "&=",
	ASSIGN_OR:  "|=",
	ASSIGN_XOR: "^=",
	ASSIGN_SHL: "<<=",
	ASSIGN_SHR: ">>=",
	ASSIGN_SAR: ">>>=",
	EQ:        "==",
	NEQ:       "!=",
	LT:        "<",
	GT:        ">",
	LTE:       "<=",
	GTE:       ">=",
	AND:       "&&",
	OR:        "||",
	NOT:       "!",
	BIT_AND:   "&",
	BIT_OR:    "|",
	BIT_XOR:   "^",
	BIT_NOT:   "~",
	SHL:       "<<",
	SHR:       ">>",
	SAR:       ">>>",
	ADD:       "+",
	SUB:       "-",
	MUL:       "*",
	DIV:       "/",
	MOD:       "%",
	EXP:       "**",
	INC:       "++",
	DEC:       "--",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

var keywords = map[string]TokenType{
	"abstract":    ABSTRACT,
	"address":     ADDRESS,
	"anonymous":   ANONYMOUS,
	"as":          AS,
	"assembly":    ASSEMBLY,
	"bool":        BOOL,
	"break":       BREAK,
	"bytes":       BYTES,
	"calldata":    CALLDATA,
	"case":        CASE,
	"catch":       CATCH,
	"constant":    CONSTANT,
	"constructor": CONSTRUCTOR,
	"continue":    CONTINUE,
	"contract":    CONTRACT,
	"default":     DEFAULT,
	"delete":      DELETE,
	"do":          DO,
	"else":        ELSE,
	"emit":        EMIT,
	"enum":        ENUM,
	"error":       ERROR,
	"event":       EVENT,
	"external":    EXTERNAL,
	"fallback":    FALLBACK,
	"false":       FALSE,
	"for":         FOR,
	"from":        FROM,
	"function":    FUNCTION,
	"global":      GLOBAL,
	"hex":         HEX,
	"if":          IF,
	"immutable":   IMMUTABLE,
	"import":      IMPORT,
	"indexed":     INDEXED,
	"interface":   INTERFACE,
	"internal":    INTERNAL,
	"is":          IS,
	"let":         LET,
	"library":     LIBRARY,
	"mapping":     MAPPING,
	"memory":      MEMORY,
	"modifier":    MODIFIER,
	"new":         NEW,
	"override":    OVERRIDE,
	"payable":     PAYABLE,
	"pragma":      PRAGMA,
	"private":     PRIVATE,
	"public":      PUBLIC,
	"pure":        PURE,
	"receive":     RECEIVE,
	"return":      RETURN,
	"returns":     RETURNS,
	"revert":      REVERT,
	"storage":     STORAGE,
	"string":      STRING_TYPE,
	"struct":      STRUCT,
	"switch":      SWITCH,
	"transient":   TRANSIENT,
	"true":        TRUE,
	"try":         TRY,
	"type":        TYPE,
	"unchecked":   UNCHECKED,
	"layout":      LAYOUT,
	"at":          AT,
	"unicode":     UNICODE,
	"using":       USING,
	"view":        VIEW,
	"virtual":     VIRTUAL,
	"while":       WHILE,
}

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
	Start   int
	End     int
}

// Position represents a position in the source
type Position struct {
	Line   int
	Column int
	Offset int
}

// Lexer tokenizes Solidity source code
type Lexer struct {
	input   string
	pos     int
	line    int
	column  int
	start   int
}

// New creates a new Lexer
func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 0,
	}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()

	if l.pos >= len(l.input) {
		return Token{Type: EOF, Line: l.line, Column: l.column, Start: l.pos, End: l.pos}
	}

	l.start = l.pos
	startLine := l.line
	startColumn := l.column

	ch := l.peek()

	// Identifiers and keywords
	if isIdentifierStart(ch) {
		return l.readIdentifier(startLine, startColumn)
	}

	// Hex numbers (must check before regular numbers since both start with digit)
	if ch == '0' && l.pos+1 < len(l.input) && (l.peekAt(1) == 'x' || l.peekAt(1) == 'X') {
		return l.readHexNumber(startLine, startColumn)
	}

	// Numbers
	if isDigit(ch) || (ch == '.' && l.pos+1 < len(l.input) && isDigit(l.peekAt(1))) {
		return l.readNumber(startLine, startColumn)
	}

	// Strings
	if ch == '"' || ch == '\'' {
		return l.readString(startLine, startColumn)
	}

	// Operators and punctuation
	return l.readOperator(startLine, startColumn)
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekAt(offset int) byte {
	pos := l.pos + offset
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

func (l *Lexer) advance() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.peek()

		// Whitespace
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.advance()
			continue
		}

		// Line comment
		if ch == '/' && l.peekAt(1) == '/' {
			l.advance() // /
			l.advance() // /
			for l.pos < len(l.input) && l.peek() != '\n' {
				l.advance()
			}
			continue
		}

		// Block comment
		if ch == '/' && l.peekAt(1) == '*' {
			l.advance() // /
			l.advance() // *
			for l.pos < len(l.input) {
				if l.peek() == '*' && l.peekAt(1) == '/' {
					l.advance() // *
					l.advance() // /
					break
				}
				l.advance()
			}
			continue
		}

		break
	}
}

func (l *Lexer) readIdentifier(line, column int) Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentifierPart(l.peek()) {
		l.advance()
	}
	value := l.input[start:l.pos]

	// Check for typed keywords (int, uint, bytes with size suffix)
	tokenType := IDENTIFIER
	if kw, ok := keywords[value]; ok {
		tokenType = kw
	} else if strings.HasPrefix(value, "int") && isIntType(value) {
		tokenType = INT
	} else if strings.HasPrefix(value, "uint") && isUintType(value) {
		tokenType = UINT
	} else if strings.HasPrefix(value, "bytes") && isBytesNType(value) {
		tokenType = BYTES_N
	} else if value == "fixed" {
		tokenType = FIXED
	} else if value == "ufixed" {
		tokenType = UFIXED
	} else if isFixedNType(value) {
		tokenType = FIXED_N
	} else if isUfixedNType(value) {
		tokenType = UFIXED_N
	}

	return Token{
		Type:   tokenType,
		Value:  value,
		Line:   line,
		Column: column,
		Start:  start,
		End:    l.pos,
	}
}

// isIntType checks if value is a valid int type (int8, int16, ..., int256)
func isIntType(value string) bool {
	if value == "int" {
		return true
	}
	if !strings.HasPrefix(value, "int") {
		return false
	}
	suffix := value[3:]
	return isValidIntSize(suffix)
}

// isUintType checks if value is a valid uint type (uint8, uint16, ..., uint256)
func isUintType(value string) bool {
	if value == "uint" {
		return true
	}
	if !strings.HasPrefix(value, "uint") {
		return false
	}
	suffix := value[4:]
	return isValidIntSize(suffix)
}

// isValidIntSize checks if the size is a valid integer size (8, 16, 24, ..., 256)
func isValidIntSize(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isBytesNType checks if value is a valid bytesN type (bytes1, bytes2, ..., bytes32)
func isBytesNType(value string) bool {
	if !strings.HasPrefix(value, "bytes") || len(value) <= 5 {
		return false
	}
	suffix := value[5:]
	for _, c := range suffix {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isFixedNType checks if value is a valid fixedMxN type (e.g., fixed128x18)
func isFixedNType(value string) bool {
	if !strings.HasPrefix(value, "fixed") {
		return false
	}
	rest := value[5:]
	return isValidFixedSuffix(rest)
}

// isUfixedNType checks if value is a valid ufixedMxN type (e.g., ufixed128x18)
func isUfixedNType(value string) bool {
	if !strings.HasPrefix(value, "ufixed") {
		return false
	}
	rest := value[6:]
	return isValidFixedSuffix(rest)
}

// isValidFixedSuffix checks if the suffix matches MxN pattern (numbers only)
func isValidFixedSuffix(s string) bool {
	xIdx := strings.Index(s, "x")
	if xIdx <= 0 || xIdx == len(s)-1 {
		return false
	}
	m := s[:xIdx]
	n := s[xIdx+1:]
	// M and N must be all digits
	for _, c := range m {
		if c < '0' || c > '9' {
			return false
		}
	}
	for _, c := range n {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (l *Lexer) readNumber(line, column int) Token {
	start := l.pos
	
	// Integer part
	for l.pos < len(l.input) && (isDigit(l.peek()) || l.peek() == '_') {
		l.advance()
	}

	// Fractional part
	if l.peek() == '.' && l.peekAt(1) != '.' {
		l.advance() // .
		for l.pos < len(l.input) && (isDigit(l.peek()) || l.peek() == '_') {
			l.advance()
		}
	}

	// Exponent
	if l.peek() == 'e' || l.peek() == 'E' {
		l.advance() // e/E
		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}
		for l.pos < len(l.input) && (isDigit(l.peek()) || l.peek() == '_') {
			l.advance()
		}
	}

	value := l.input[start:l.pos]
	// Remove underscores for the actual value
	value = strings.ReplaceAll(value, "_", "")

	return Token{
		Type:   NUMBER,
		Value:  value,
		Line:   line,
		Column: column,
		Start:  start,
		End:    l.pos,
	}
}

func (l *Lexer) readHexNumber(line, column int) Token {
	start := l.pos
	l.advance() // 0
	l.advance() // x/X

	for l.pos < len(l.input) && (isHexDigit(l.peek()) || l.peek() == '_') {
		l.advance()
	}

	value := l.input[start:l.pos]

	return Token{
		Type:   HEX_NUMBER,
		Value:  value,
		Line:   line,
		Column: column,
		Start:  start,
		End:    l.pos,
	}
}

func (l *Lexer) readString(line, column int) Token {
	start := l.pos
	quote := l.advance() // opening quote
	
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == quote {
			l.advance() // closing quote
			break
		}
		if ch == '\\' && l.pos+1 < len(l.input) {
			l.advance() // backslash
			escaped := l.advance()
			switch escaped {
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case '\\':
				sb.WriteByte('\\')
			case '\'':
				sb.WriteByte('\'')
			case '"':
				sb.WriteByte('"')
			default:
				sb.WriteByte(escaped)
			}
			continue
		}
		if ch == '\n' || ch == '\r' {
			break // Unterminated string
		}
		sb.WriteByte(l.advance())
	}

	return Token{
		Type:   STRING,
		Value:  sb.String(),
		Line:   line,
		Column: column,
		Start:  start,
		End:    l.pos,
	}
}

func (l *Lexer) readOperator(line, column int) Token {
	start := l.pos
	ch := l.advance()

	// Three-character operators
	if l.pos+1 < len(l.input) {
		three := string(ch) + string(l.peek()) + string(l.peekAt(1))
		switch three {
		case ">>>":
			l.advance()
			l.advance()
			if l.peek() == '=' {
				l.advance()
				return Token{Type: ASSIGN_SAR, Value: ">>>=", Line: line, Column: column, Start: start, End: l.pos}
			}
			return Token{Type: SAR, Value: ">>>", Line: line, Column: column, Start: start, End: l.pos}
		case ">>=":
			l.advance()
			l.advance()
			return Token{Type: ASSIGN_SHR, Value: ">>=", Line: line, Column: column, Start: start, End: l.pos}
		case "<<=":
			l.advance()
			l.advance()
			return Token{Type: ASSIGN_SHL, Value: "<<=", Line: line, Column: column, Start: start, End: l.pos}
		}
	}

	// Two-character operators
	if l.pos < len(l.input) {
		two := string(ch) + string(l.peek())
		switch two {
		case "=>":
			l.advance()
			return Token{Type: ARROW, Value: "=>", Line: line, Column: column, Start: start, End: l.pos}
		case "->":
			l.advance()
			return Token{Type: RIGHT_ARROW, Value: "->", Line: line, Column: column, Start: start, End: l.pos}
		case "==":
			l.advance()
			return Token{Type: EQ, Value: "==", Line: line, Column: column, Start: start, End: l.pos}
		case "!=":
			l.advance()
			return Token{Type: NEQ, Value: "!=", Line: line, Column: column, Start: start, End: l.pos}
		case "<=":
			l.advance()
			return Token{Type: LTE, Value: "<=", Line: line, Column: column, Start: start, End: l.pos}
		case ">=":
			l.advance()
			return Token{Type: GTE, Value: ">=", Line: line, Column: column, Start: start, End: l.pos}
		case "&&":
			l.advance()
			return Token{Type: AND, Value: "&&", Line: line, Column: column, Start: start, End: l.pos}
		case "||":
			l.advance()
			return Token{Type: OR, Value: "||", Line: line, Column: column, Start: start, End: l.pos}
		case "<<":
			l.advance()
			return Token{Type: SHL, Value: "<<", Line: line, Column: column, Start: start, End: l.pos}
		case ">>":
			l.advance()
			return Token{Type: SHR, Value: ">>", Line: line, Column: column, Start: start, End: l.pos}
		case "**":
			l.advance()
			return Token{Type: EXP, Value: "**", Line: line, Column: column, Start: start, End: l.pos}
		case "++":
			l.advance()
			return Token{Type: INC, Value: "++", Line: line, Column: column, Start: start, End: l.pos}
		case "--":
			l.advance()
			return Token{Type: DEC, Value: "--", Line: line, Column: column, Start: start, End: l.pos}
		case "+=":
			l.advance()
			return Token{Type: ASSIGN_ADD, Value: "+=", Line: line, Column: column, Start: start, End: l.pos}
		case "-=":
			l.advance()
			return Token{Type: ASSIGN_SUB, Value: "-=", Line: line, Column: column, Start: start, End: l.pos}
		case "*=":
			l.advance()
			return Token{Type: ASSIGN_MUL, Value: "*=", Line: line, Column: column, Start: start, End: l.pos}
		case "/=":
			l.advance()
			return Token{Type: ASSIGN_DIV, Value: "/=", Line: line, Column: column, Start: start, End: l.pos}
		case "%=":
			l.advance()
			return Token{Type: ASSIGN_MOD, Value: "%=", Line: line, Column: column, Start: start, End: l.pos}
		case "&=":
			l.advance()
			return Token{Type: ASSIGN_AND, Value: "&=", Line: line, Column: column, Start: start, End: l.pos}
		case "|=":
			l.advance()
			return Token{Type: ASSIGN_OR, Value: "|=", Line: line, Column: column, Start: start, End: l.pos}
		case "^=":
			l.advance()
			return Token{Type: ASSIGN_XOR, Value: "^=", Line: line, Column: column, Start: start, End: l.pos}
		}
	}

	// Single-character operators
	switch ch {
	case '(':
		return Token{Type: LPAREN, Value: "(", Line: line, Column: column, Start: start, End: l.pos}
	case ')':
		return Token{Type: RPAREN, Value: ")", Line: line, Column: column, Start: start, End: l.pos}
	case '[':
		return Token{Type: LBRACK, Value: "[", Line: line, Column: column, Start: start, End: l.pos}
	case ']':
		return Token{Type: RBRACK, Value: "]", Line: line, Column: column, Start: start, End: l.pos}
	case '{':
		return Token{Type: LBRACE, Value: "{", Line: line, Column: column, Start: start, End: l.pos}
	case '}':
		return Token{Type: RBRACE, Value: "}", Line: line, Column: column, Start: start, End: l.pos}
	case ':':
		return Token{Type: COLON, Value: ":", Line: line, Column: column, Start: start, End: l.pos}
	case ';':
		return Token{Type: SEMICOLON, Value: ";", Line: line, Column: column, Start: start, End: l.pos}
	case '.':
		return Token{Type: PERIOD, Value: ".", Line: line, Column: column, Start: start, End: l.pos}
	case ',':
		return Token{Type: COMMA, Value: ",", Line: line, Column: column, Start: start, End: l.pos}
	case '?':
		return Token{Type: QUESTION, Value: "?", Line: line, Column: column, Start: start, End: l.pos}
	case '=':
		return Token{Type: ASSIGN, Value: "=", Line: line, Column: column, Start: start, End: l.pos}
	case '<':
		return Token{Type: LT, Value: "<", Line: line, Column: column, Start: start, End: l.pos}
	case '>':
		return Token{Type: GT, Value: ">", Line: line, Column: column, Start: start, End: l.pos}
	case '!':
		return Token{Type: NOT, Value: "!", Line: line, Column: column, Start: start, End: l.pos}
	case '&':
		return Token{Type: BIT_AND, Value: "&", Line: line, Column: column, Start: start, End: l.pos}
	case '|':
		return Token{Type: BIT_OR, Value: "|", Line: line, Column: column, Start: start, End: l.pos}
	case '^':
		return Token{Type: BIT_XOR, Value: "^", Line: line, Column: column, Start: start, End: l.pos}
	case '~':
		return Token{Type: BIT_NOT, Value: "~", Line: line, Column: column, Start: start, End: l.pos}
	case '+':
		return Token{Type: ADD, Value: "+", Line: line, Column: column, Start: start, End: l.pos}
	case '-':
		return Token{Type: SUB, Value: "-", Line: line, Column: column, Start: start, End: l.pos}
	case '*':
		return Token{Type: MUL, Value: "*", Line: line, Column: column, Start: start, End: l.pos}
	case '/':
		return Token{Type: DIV, Value: "/", Line: line, Column: column, Start: start, End: l.pos}
	case '%':
		return Token{Type: MOD, Value: "%", Line: line, Column: column, Start: start, End: l.pos}
	}

	return Token{Type: ILLEGAL, Value: string(ch), Line: line, Column: column, Start: start, End: l.pos}
}

// Tokenize returns all tokens from the input
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}

func isIdentifierStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$'
}

func isIdentifierPart(ch byte) bool {
	return isIdentifierStart(ch) || isDigit(ch)
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// IsKeyword checks if a token type is a keyword
func IsKeyword(t TokenType) bool {
	return t >= ABSTRACT && t <= WHILE
}

// IsIdentifier checks if a rune is valid in an identifier
func IsIdentifier(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$'
}

