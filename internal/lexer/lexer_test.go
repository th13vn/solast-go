package lexer

import (
	"testing"
)

func TestEventLexing(t *testing.T) {
	input := `event Transfer(address indexed from);`
	
	lex := New(input)
	var tokens []Token
	for {
		tok := lex.NextToken()
		tokens = append(tokens, tok)
		t.Logf("Token: %s Value: %q Line: %d Col: %d", tok.Type, tok.Value, tok.Line, tok.Column)
		if tok.Type == EOF {
			break
		}
	}
	
	// Verify expected tokens
	// Note: 'from' is tokenized as FROM keyword, not IDENTIFIER
	expected := []TokenType{EVENT, IDENTIFIER, LPAREN, ADDRESS, INDEXED, FROM, RPAREN, SEMICOLON, EOF}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	
	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("Token %d: expected %s, got %s (value: %q)", i, exp, tokens[i].Type, tokens[i].Value)
		}
	}
}

func TestTransientKeyword(t *testing.T) {
	input := `uint256 transient x`
	lex := New(input)
	
	var tokens []Token
	for {
		tok := lex.NextToken()
		tokens = append(tokens, tok)
		t.Logf("Token: Type=%s Value=%q", tok.Type, tok.Value)
		if tok.Type == EOF {
			break
		}
	}
	
	// Expected: UINT, TRANSIENT, IDENTIFIER, EOF
	expected := []TokenType{UINT, TRANSIENT, IDENTIFIER, EOF}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	
	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("Token %d: expected %s, got %s (value: %q)", i, exp, tokens[i].Type, tokens[i].Value)
		}
	}
}

func TestHexNumber(t *testing.T) {
	input := `0x100`
	lex := New(input)
	tok := lex.NextToken()
	t.Logf("Token: Type=%s Value=%q", tok.Type, tok.Value)
	if tok.Type != HEX_NUMBER {
		t.Errorf("Expected HEX_NUMBER, got %s", tok.Type)
	}
}

func TestBasicTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"address", ADDRESS},
		{"bool", BOOL},
		{"string", STRING_TYPE},
		{"bytes", BYTES},
		{"uint256", UINT},
		{"int256", INT},
		{"bytes32", BYTES_N},
		{"indexed", INDEXED},
		{"transient", TRANSIENT},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lex := New(tt.input)
			tok := lex.NextToken()
			if tok.Type != tt.expected {
				t.Errorf("Expected %s, got %s (value: %q)", tt.expected, tok.Type, tok.Value)
			}
		})
	}
}

