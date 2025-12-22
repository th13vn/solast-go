package builder

import (
	"fmt"

	"github.com/th13vn/solast-go/internal/lexer"
	"github.com/th13vn/solast-go/pkg/ast"
)

// Token navigation helpers

func (b *Builder) peek() lexer.Token {
	if b.pos >= len(b.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return b.tokens[b.pos]
}

func (b *Builder) previous() lexer.Token {
	if b.pos == 0 {
		return lexer.Token{Type: lexer.EOF}
	}
	return b.tokens[b.pos-1]
}

func (b *Builder) advance() lexer.Token {
	if !b.isAtEnd() {
		b.pos++
	}
	return b.previous()
}

func (b *Builder) check(t lexer.TokenType) bool {
	if b.isAtEnd() {
		return false
	}
	return b.peek().Type == t
}

func (b *Builder) isAtEnd() bool {
	return b.peek().Type == lexer.EOF
}

func (b *Builder) expect(t lexer.TokenType) lexer.Token {
	if b.check(t) {
		return b.advance()
	}
	b.addError(fmt.Sprintf("expected '%s', got '%s'", t.String(), b.peek().Value))
	// Advance even on failure to prevent infinite loops in non-tolerant mode
	if !b.options.Tolerant {
		b.advance()
	}
	return b.peek()
}

func (b *Builder) expectKeyword(keyword string) lexer.Token {
	if b.peek().Value == keyword {
		return b.advance()
	}
	b.addError(fmt.Sprintf("expected '%s', got '%s'", keyword, b.peek().Value))
	// Advance even on failure to prevent infinite loops in non-tolerant mode
	if !b.options.Tolerant {
		b.advance()
	}
	return b.peek()
}

// Error handling

func (b *Builder) addError(message string) {
	tok := b.peek()
	b.errors = append(b.errors, &Error{
		Message: message,
		Line:    tok.Line,
		Column:  tok.Column,
	})
	
	if b.options.Tolerant {
		// Try to recover by skipping to next statement
		b.synchronize()
	}
}

func (b *Builder) synchronize() {
	b.advance()
	
	for !b.isAtEnd() {
		if b.previous().Type == lexer.SEMICOLON {
			return
		}
		
		switch b.peek().Type {
		case lexer.CONTRACT, lexer.INTERFACE, lexer.LIBRARY, lexer.FUNCTION,
			lexer.STRUCT, lexer.ENUM, lexer.EVENT, lexer.ERROR,
			lexer.PRAGMA, lexer.IMPORT, lexer.USING:
			return
		}
		
		b.advance()
	}
}

// Type checking helpers

func (b *Builder) isTypeName() bool {
	return b.isElementaryTypeName() || b.check(lexer.MAPPING) || 
		b.check(lexer.FUNCTION) || b.check(lexer.IDENTIFIER)
}

func (b *Builder) isElementaryTypeName() bool {
	t := b.peek().Type
	return t == lexer.ADDRESS || t == lexer.BOOL || t == lexer.STRING_TYPE ||
		t == lexer.BYTES || t == lexer.INT || t == lexer.UINT ||
		t == lexer.BYTE || t == lexer.BYTES_N || t == lexer.FIXED ||
		t == lexer.UFIXED || t == lexer.FIXED_N || t == lexer.UFIXED_N
}

// isContextualKeyword checks if the current token is a keyword that can be used as an identifier
// in certain contexts (e.g., 'from' in parameter names)
// Based on official grammar: From | Error | Revert | Global | Transient | Layout | At
func (b *Builder) isContextualKeyword() bool {
	t := b.peek().Type
	// These keywords can be used as identifiers in certain contexts
	return t == lexer.FROM || t == lexer.ERROR || t == lexer.REVERT ||
		t == lexer.GLOBAL || t == lexer.TRANSIENT || t == lexer.LAYOUT || t == lexer.AT
}

// locationSetter is an interface for nodes that can have location set
type locationSetter interface {
	setLoc(*ast.Location)
	setRange(*ast.Range)
}

// Location helpers - simplified version that uses reflection-like approach
func (b *Builder) setLocation(node ast.Node, startTok, endTok lexer.Token) {
	if !b.options.Loc && !b.options.Range {
		return
	}
	
	var loc *ast.Location
	var rng *ast.Range
	
	if b.options.Loc {
		loc = &ast.Location{
			Start: ast.Position{Line: startTok.Line, Column: startTok.Column},
			End:   ast.Position{Line: endTok.Line, Column: endTok.Column + len(endTok.Value)},
		}
	}
	
	if b.options.Range {
		rng = &ast.Range{startTok.Start, endTok.End}
	}
	
	// Set location based on node type
	switch n := node.(type) {
	case *ast.SourceUnit:
		n.Loc, n.Range = loc, rng
	case *ast.PragmaDirective:
		n.Loc, n.Range = loc, rng
	case *ast.ImportDirective:
		n.Loc, n.Range = loc, rng
	case *ast.ContractDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.InheritanceSpecifier:
		n.Loc, n.Range = loc, rng
	case *ast.FunctionDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.ModifierDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.ModifierInvocation:
		n.Loc, n.Range = loc, rng
	case *ast.StateVariableDeclaration:
		n.Loc, n.Range = loc, rng
	case *ast.VariableDeclaration:
		n.Loc, n.Range = loc, rng
	case *ast.VariableDeclarationStatement:
		n.Loc, n.Range = loc, rng
	case *ast.StructDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.EnumDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.EventDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.ErrorDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.UsingForDeclaration:
		n.Loc, n.Range = loc, rng
	case *ast.UserDefinedValueTypeDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.ElementaryTypeName:
		n.Loc, n.Range = loc, rng
	case *ast.UserDefinedTypeName:
		n.Loc, n.Range = loc, rng
	case *ast.Mapping:
		n.Loc, n.Range = loc, rng
	case *ast.ArrayTypeName:
		n.Loc, n.Range = loc, rng
	case *ast.FunctionTypeName:
		n.Loc, n.Range = loc, rng
	case *ast.Block:
		n.Loc, n.Range = loc, rng
	case *ast.UncheckedBlock:
		n.Loc, n.Range = loc, rng
	case *ast.ExpressionStatement:
		n.Loc, n.Range = loc, rng
	case *ast.IfStatement:
		n.Loc, n.Range = loc, rng
	case *ast.WhileStatement:
		n.Loc, n.Range = loc, rng
	case *ast.DoWhileStatement:
		n.Loc, n.Range = loc, rng
	case *ast.ForStatement:
		n.Loc, n.Range = loc, rng
	case *ast.ContinueStatement:
		n.Loc, n.Range = loc, rng
	case *ast.BreakStatement:
		n.Loc, n.Range = loc, rng
	case *ast.ReturnStatement:
		n.Loc, n.Range = loc, rng
	case *ast.EmitStatement:
		n.Loc, n.Range = loc, rng
	case *ast.RevertStatement:
		n.Loc, n.Range = loc, rng
	case *ast.TryStatement:
		n.Loc, n.Range = loc, rng
	case *ast.CatchClause:
		n.Loc, n.Range = loc, rng
	case *ast.InlineAssembly:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyBlock:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyCall:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyLocalDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyAssignment:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyIdentifier:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyLiteral:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyIf:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblySwitch:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyFor:
		n.Loc, n.Range = loc, rng
	case *ast.AssemblyFunctionDefinition:
		n.Loc, n.Range = loc, rng
	case *ast.TupleExpression:
		n.Loc, n.Range = loc, rng
	case *ast.NewExpression:
		n.Loc, n.Range = loc, rng
	case *ast.FunctionCall:
		n.Loc, n.Range = loc, rng
	case *ast.BinaryOperation:
		n.Loc, n.Range = loc, rng
	case *ast.UnaryOperation:
		n.Loc, n.Range = loc, rng
	case *ast.Conditional:
		n.Loc, n.Range = loc, rng
	case *ast.MemberAccess:
		n.Loc, n.Range = loc, rng
	case *ast.IndexAccess:
		n.Loc, n.Range = loc, rng
	case *ast.IndexRangeAccess:
		n.Loc, n.Range = loc, rng
	case *ast.FunctionCallOptions:
		n.Loc, n.Range = loc, rng
	case *ast.Identifier:
		n.Loc, n.Range = loc, rng
	case *ast.NumberLiteral:
		n.Loc, n.Range = loc, rng
	case *ast.BooleanLiteral:
		n.Loc, n.Range = loc, rng
	case *ast.StringLiteral:
		n.Loc, n.Range = loc, rng
	case *ast.HexLiteral:
		n.Loc, n.Range = loc, rng
	}
}
