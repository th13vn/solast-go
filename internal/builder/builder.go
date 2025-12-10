// Package builder provides the AST builder from tokens.
package builder

import (
	"fmt"
	"strings"

	"github.com/th13vn/solast-go/internal/lexer"
	"github.com/th13vn/solast-go/pkg/ast"
)

// Error represents a parsing error
type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Error() string {
	return fmt.Sprintf("line %d:%d: %s", e.Line, e.Column, e.Message)
}

// Builder builds an AST from Solidity source code
type Builder struct {
	tokens   []lexer.Token
	pos      int
	errors   []*Error
	options  *Options
}

// Options configures the parser behavior
type Options struct {
	Tolerant bool // Collect errors instead of stopping
	Loc      bool // Add location information
	Range    bool // Add range information
}

// New creates a new Builder
func New(input string, opts *Options) *Builder {
	lex := lexer.New(input)
	tokens := lex.Tokenize()
	
	if opts == nil {
		opts = &Options{}
	}
	
	return &Builder{
		tokens:  tokens,
		pos:     0,
		errors:  make([]*Error, 0),
		options: opts,
	}
}

// Build parses the source and returns the AST
func (b *Builder) Build() (*ast.SourceUnit, error) {
	sourceUnit := &ast.SourceUnit{
		BaseNode: ast.BaseNode{Type: ast.NodeSourceUnit},
		Children: make([]ast.Node, 0),
	}

	for !b.isAtEnd() {
		node := b.parseSourceUnitElement()
		if node != nil {
			sourceUnit.Children = append(sourceUnit.Children, node)
		}
		if len(b.errors) > 0 && !b.options.Tolerant {
			return nil, b.errors[0]
		}
	}

	if b.options.Loc {
		if len(sourceUnit.Children) > 0 {
			first := sourceUnit.Children[0]
			last := sourceUnit.Children[len(sourceUnit.Children)-1]
			if first.GetLocation() != nil && last.GetLocation() != nil {
				sourceUnit.Loc = &ast.Location{
					Start: first.GetLocation().Start,
					End:   last.GetLocation().End,
				}
			}
		}
	}

	return sourceUnit, nil
}

// Errors returns the collected parsing errors
func (b *Builder) Errors() []*Error {
	return b.errors
}

func (b *Builder) parseSourceUnitElement() ast.Node {
	tok := b.peek()
	
	switch tok.Type {
	case lexer.PRAGMA:
		return b.parsePragmaDirective()
	case lexer.IMPORT:
		return b.parseImportDirective()
	case lexer.CONTRACT:
		return b.parseContractDefinition("contract")
	case lexer.ABSTRACT:
		b.advance() // abstract
		if b.check(lexer.CONTRACT) {
			return b.parseContractDefinition("abstract")
		}
		b.addError("expected 'contract' after 'abstract'")
		return nil
	case lexer.INTERFACE:
		return b.parseContractDefinition("interface")
	case lexer.LIBRARY:
		return b.parseContractDefinition("library")
	case lexer.STRUCT:
		return b.parseStructDefinition()
	case lexer.ENUM:
		return b.parseEnumDefinition()
	case lexer.FUNCTION:
		return b.parseFunctionDefinition()
	case lexer.EVENT:
		return b.parseEventDefinition()
	case lexer.ERROR:
		return b.parseErrorDefinition()
	case lexer.USING:
		return b.parseUsingDirective()
	case lexer.TYPE:
		return b.parseUserDefinedValueTypeDefinition()
	case lexer.EOF:
		return nil
	default:
		// Try to parse as constant variable
		if b.isTypeName() {
			return b.parseConstantVariableDeclaration()
		}
		b.addError(fmt.Sprintf("unexpected token: %s", tok.Value))
		b.advance()
		return nil
	}
}

func (b *Builder) parsePragmaDirective() *ast.PragmaDirective {
	startTok := b.advance() // pragma
	
	nameTok := b.advance()
	name := nameTok.Value
	
	// Read pragma value until semicolon
	var valueParts []string
	for !b.check(lexer.SEMICOLON) && !b.isAtEnd() {
		valueParts = append(valueParts, b.advance().Value)
	}
	value := strings.Join(valueParts, " ")
	
	endTok := b.expect(lexer.SEMICOLON)
	
	node := &ast.PragmaDirective{
		BaseNode: ast.BaseNode{Type: ast.NodePragmaDirective},
		Name:     name,
		Value:    value,
	}
	
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseImportDirective() *ast.ImportDirective {
	startTok := b.advance() // import
	
	node := &ast.ImportDirective{
		BaseNode: ast.BaseNode{Type: ast.NodeImportDirective},
	}
	
	if b.check(lexer.STRING) {
		// import "path";
		pathTok := b.advance()
		node.Path = pathTok.Value
		
		if b.check(lexer.AS) {
			b.advance() // as
			aliasTok := b.expect(lexer.IDENTIFIER)
			node.UnitAlias = aliasTok.Value
		}
	} else if b.check(lexer.MUL) {
		// import * as alias from "path";
		b.advance() // *
		b.expect(lexer.AS)
		aliasTok := b.expect(lexer.IDENTIFIER)
		node.UnitAlias = aliasTok.Value
		b.expectKeyword("from")
		pathTok := b.expect(lexer.STRING)
		node.Path = pathTok.Value
	} else if b.check(lexer.LBRACE) {
		// import { sym1, sym2 as alias } from "path";
		b.advance() // {
		node.SymbolAliases = make([]*ast.ImportSymbol, 0)
		
		for !b.check(lexer.RBRACE) && !b.isAtEnd() {
			sym := &ast.ImportSymbol{}
			symTok := b.expect(lexer.IDENTIFIER)
			sym.Symbol = symTok.Value
			
			if b.check(lexer.AS) {
				b.advance() // as
				aliasTok := b.expect(lexer.IDENTIFIER)
				sym.Alias = aliasTok.Value
			}
			
			node.SymbolAliases = append(node.SymbolAliases, sym)
			
			if !b.check(lexer.RBRACE) {
				b.expect(lexer.COMMA)
			}
		}
		b.expect(lexer.RBRACE)
		b.expectKeyword("from")
		pathTok := b.expect(lexer.STRING)
		node.Path = pathTok.Value
	} else if b.check(lexer.IDENTIFIER) {
		// import Identifier from "path";
		nameTok := b.advance()
		node.UnitAlias = nameTok.Value
		b.expectKeyword("from")
		pathTok := b.expect(lexer.STRING)
		node.Path = pathTok.Value
	}
	
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseContractDefinition(kind string) *ast.ContractDefinition {
	var startTok lexer.Token
	if kind == "abstract" {
		startTok = b.previous()
		b.advance() // skip 'contract' keyword
	} else {
		startTok = b.advance() // contract/interface/library
	}
	
	nameTok := b.expect(lexer.IDENTIFIER)
	
	node := &ast.ContractDefinition{
		BaseNode:      ast.BaseNode{Type: ast.NodeContractDefinition},
		Name:          nameTok.Value,
		Kind:          kind,
		BaseContracts: make([]*ast.InheritanceSpecifier, 0),
		SubNodes:      make([]ast.Node, 0),
	}
	
	// Layout directive (Solidity 0.8.24+): contract Foo layout at 0x100 { }
	if b.check(lexer.LAYOUT) {
		b.advance() // layout
		b.expect(lexer.AT)
		// Parse the layout expression - just a primary expression (number/identifier)
		// We don't use parseExpression() because { would be parsed as FunctionCallOptions
		_ = b.parsePrimary() // Layout expression - could be stored in AST if needed
	}
	
	// Inheritance
	if b.check(lexer.IS) {
		b.advance() // is
		for {
			base := b.parseInheritanceSpecifier()
			node.BaseContracts = append(node.BaseContracts, base)
			if !b.check(lexer.COMMA) {
				break
			}
			b.advance() // ,
		}
	}
	
	b.expect(lexer.LBRACE)
	
	// Contract body
	for !b.check(lexer.RBRACE) && !b.isAtEnd() {
		subNode := b.parseContractBodyElement()
		if subNode != nil {
			node.SubNodes = append(node.SubNodes, subNode)
		}
	}
	
	endTok := b.expect(lexer.RBRACE)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseInheritanceSpecifier() *ast.InheritanceSpecifier {
	startTok := b.peek()
	
	baseName := b.parseUserDefinedTypeName()
	
	node := &ast.InheritanceSpecifier{
		BaseNode: ast.BaseNode{Type: ast.NodeInheritanceSpecifier},
		BaseName: baseName,
	}
	
	// Constructor arguments
	if b.check(lexer.LPAREN) {
		b.advance() // (
		node.Arguments = b.parseExpressionList()
		b.expect(lexer.RPAREN)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseContractBodyElement() ast.Node {
	tok := b.peek()
	
	switch tok.Type {
	case lexer.FUNCTION:
		return b.parseFunctionDefinition()
	case lexer.CONSTRUCTOR:
		return b.parseConstructorDefinition()
	case lexer.MODIFIER:
		return b.parseModifierDefinition()
	case lexer.FALLBACK:
		return b.parseFallbackDefinition()
	case lexer.RECEIVE:
		return b.parseReceiveDefinition()
	case lexer.STRUCT:
		return b.parseStructDefinition()
	case lexer.ENUM:
		return b.parseEnumDefinition()
	case lexer.EVENT:
		return b.parseEventDefinition()
	case lexer.ERROR:
		return b.parseErrorDefinition()
	case lexer.USING:
		return b.parseUsingDirective()
	case lexer.TYPE:
		return b.parseUserDefinedValueTypeDefinition()
	default:
		// State variable
		if b.isTypeName() {
			return b.parseStateVariableDeclaration()
		}
		b.addError(fmt.Sprintf("unexpected token in contract body: %s", tok.Value))
		b.advance()
		return nil
	}
}

func (b *Builder) parseFunctionDefinition() *ast.FunctionDefinition {
	startTok := b.advance() // function
	
	node := &ast.FunctionDefinition{
		BaseNode:   ast.BaseNode{Type: ast.NodeFunctionDefinition},
		Parameters: make([]*ast.VariableDeclaration, 0),
		Modifiers:  make([]*ast.ModifierInvocation, 0),
	}
	
	// Function name (optional for fallback/receive)
	if b.check(lexer.IDENTIFIER) || b.check(lexer.FALLBACK) || b.check(lexer.RECEIVE) {
		nameTok := b.advance()
		node.Name = nameTok.Value
		if nameTok.Type == lexer.FALLBACK {
			node.IsFallback = true
		} else if nameTok.Type == lexer.RECEIVE {
			node.IsReceiveEther = true
		}
	}
	
	// Parameters
	node.Parameters = b.parseParameterList()
	
	// Modifiers, visibility, state mutability
	b.parseFunctionModifiers(node)
	
	// Return parameters
	if b.check(lexer.RETURNS) {
		b.advance() // returns
		node.ReturnParameters = b.parseParameterList()
	}
	
	// Body or semicolon
	if b.check(lexer.LBRACE) {
		node.Body = b.parseBlock()
	} else {
		b.expect(lexer.SEMICOLON)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseConstructorDefinition() *ast.FunctionDefinition {
	startTok := b.advance() // constructor
	
	node := &ast.FunctionDefinition{
		BaseNode:      ast.BaseNode{Type: ast.NodeFunctionDefinition},
		IsConstructor: true,
		Parameters:    make([]*ast.VariableDeclaration, 0),
		Modifiers:     make([]*ast.ModifierInvocation, 0),
	}
	
	node.Parameters = b.parseParameterList()
	b.parseFunctionModifiers(node)
	node.Body = b.parseBlock()
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseFallbackDefinition() *ast.FunctionDefinition {
	startTok := b.advance() // fallback
	
	node := &ast.FunctionDefinition{
		BaseNode:   ast.BaseNode{Type: ast.NodeFunctionDefinition},
		IsFallback: true,
		Parameters: make([]*ast.VariableDeclaration, 0),
		Modifiers:  make([]*ast.ModifierInvocation, 0),
	}
	
	node.Parameters = b.parseParameterList()
	b.parseFunctionModifiers(node)
	
	if b.check(lexer.RETURNS) {
		b.advance()
		node.ReturnParameters = b.parseParameterList()
	}
	
	if b.check(lexer.LBRACE) {
		node.Body = b.parseBlock()
	} else {
		b.expect(lexer.SEMICOLON)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseReceiveDefinition() *ast.FunctionDefinition {
	startTok := b.advance() // receive
	
	node := &ast.FunctionDefinition{
		BaseNode:       ast.BaseNode{Type: ast.NodeFunctionDefinition},
		IsReceiveEther: true,
		Parameters:     make([]*ast.VariableDeclaration, 0),
		Modifiers:      make([]*ast.ModifierInvocation, 0),
	}
	
	node.Parameters = b.parseParameterList()
	b.parseFunctionModifiers(node)
	
	if b.check(lexer.LBRACE) {
		node.Body = b.parseBlock()
	} else {
		b.expect(lexer.SEMICOLON)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseFunctionModifiers(node *ast.FunctionDefinition) {
	for {
		tok := b.peek()
		switch tok.Type {
		case lexer.PUBLIC:
			b.advance()
			node.Visibility = "public"
		case lexer.PRIVATE:
			b.advance()
			node.Visibility = "private"
		case lexer.INTERNAL:
			b.advance()
			node.Visibility = "internal"
		case lexer.EXTERNAL:
			b.advance()
			node.Visibility = "external"
		case lexer.PURE:
			b.advance()
			node.StateMutability = "pure"
		case lexer.VIEW:
			b.advance()
			node.StateMutability = "view"
		case lexer.PAYABLE:
			b.advance()
			node.StateMutability = "payable"
		case lexer.VIRTUAL:
			b.advance()
			node.IsVirtual = true
		case lexer.OVERRIDE:
			b.advance()
			// Parse override specifier if present
			if b.check(lexer.LPAREN) {
				b.advance()
				for !b.check(lexer.RPAREN) && !b.isAtEnd() {
					typeName := b.parseUserDefinedTypeName()
					node.Override = append(node.Override, typeName)
					if !b.check(lexer.RPAREN) {
						b.expect(lexer.COMMA)
					}
				}
				b.expect(lexer.RPAREN)
			}
		case lexer.IDENTIFIER:
			// Modifier invocation
			mod := b.parseModifierInvocation()
			node.Modifiers = append(node.Modifiers, mod)
		default:
			return
		}
	}
}

func (b *Builder) parseModifierInvocation() *ast.ModifierInvocation {
	startTok := b.advance() // identifier
	
	node := &ast.ModifierInvocation{
		BaseNode: ast.BaseNode{Type: ast.NodeModifierInvocation},
		Name:     startTok.Value,
	}
	
	if b.check(lexer.LPAREN) {
		b.advance() // (
		node.Arguments = b.parseExpressionList()
		b.expect(lexer.RPAREN)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseModifierDefinition() *ast.ModifierDefinition {
	startTok := b.advance() // modifier
	nameTok := b.expect(lexer.IDENTIFIER)
	
	node := &ast.ModifierDefinition{
		BaseNode: ast.BaseNode{Type: ast.NodeModifierDefinition},
		Name:     nameTok.Value,
	}
	
	if b.check(lexer.LPAREN) {
		node.Parameters = b.parseParameterList()
	}
	
	// virtual/override
	for {
		if b.check(lexer.VIRTUAL) {
			b.advance()
			node.IsVirtual = true
		} else if b.check(lexer.OVERRIDE) {
			b.advance()
			if b.check(lexer.LPAREN) {
				b.advance()
				for !b.check(lexer.RPAREN) && !b.isAtEnd() {
					typeName := b.parseUserDefinedTypeName()
					node.Override = append(node.Override, typeName)
					if !b.check(lexer.RPAREN) {
						b.expect(lexer.COMMA)
					}
				}
				b.expect(lexer.RPAREN)
			}
		} else {
			break
		}
	}
	
	if b.check(lexer.LBRACE) {
		node.Body = b.parseBlock()
	} else {
		b.expect(lexer.SEMICOLON)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

