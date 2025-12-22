package builder

import (
	"strings"

	"github.com/th13vn/solast-go/internal/lexer"
	"github.com/th13vn/solast-go/pkg/ast"
)

func (b *Builder) parseTypeName() ast.Node {
	startTok := b.peek()
	
	var typeName ast.Node
	
	// Function type
	if b.check(lexer.FUNCTION) {
		typeName = b.parseFunctionTypeName()
	} else if b.check(lexer.MAPPING) {
		typeName = b.parseMappingType()
	} else if b.isElementaryTypeName() {
		typeName = b.parseElementaryTypeName()
	} else {
		typeName = b.parseUserDefinedTypeName()
	}
	
	// Array dimensions
	for b.check(lexer.LBRACK) {
		b.advance() // [
		var length ast.Node
		if !b.check(lexer.RBRACK) {
			length = b.parseExpression()
		}
		b.expect(lexer.RBRACK)
		
		typeName = &ast.ArrayTypeName{
			BaseNode:     ast.BaseNode{Type: ast.NodeArrayTypeName},
			BaseTypeName: typeName,
			Length:       length,
		}
	}
	
	b.setLocation(typeName, startTok, b.previous())
	return typeName
}

func (b *Builder) parseElementaryTypeName() *ast.ElementaryTypeName {
	startTok := b.advance()
	
	node := &ast.ElementaryTypeName{
		BaseNode: ast.BaseNode{Type: ast.NodeElementaryTypeName},
		Name:     startTok.Value,
	}
	
	// Handle address payable
	if startTok.Type == lexer.ADDRESS && b.check(lexer.PAYABLE) {
		b.advance() // payable
		node.StateMutability = "payable"
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseUserDefinedTypeName() *ast.UserDefinedTypeName {
	startTok := b.peek()
	
	var parts []string
	nameTok := b.expect(lexer.IDENTIFIER)
	parts = append(parts, nameTok.Value)
	
	for b.check(lexer.PERIOD) {
		b.advance() // .
		partTok := b.expect(lexer.IDENTIFIER)
		parts = append(parts, partTok.Value)
	}
	
	node := &ast.UserDefinedTypeName{
		BaseNode: ast.BaseNode{Type: ast.NodeUserDefinedTypeName},
		NamePath: strings.Join(parts, "."),
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseMappingType() *ast.Mapping {
	startTok := b.advance() // mapping
	
	b.expect(lexer.LPAREN)
	
	// Key type
	var keyType ast.Node
	if b.isElementaryTypeName() {
		keyType = b.parseElementaryTypeName()
	} else {
		keyType = b.parseUserDefinedTypeName()
	}
	
	// Optional key name
	var keyName *ast.Identifier
	if b.check(lexer.IDENTIFIER) {
		keyNameTok := b.advance()
		keyName = &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     keyNameTok.Value,
		}
	}
	
	b.expect(lexer.ARROW)
	
	// Value type
	valueType := b.parseTypeName()
	
	// Optional value name
	var valueName *ast.Identifier
	if b.check(lexer.IDENTIFIER) {
		valueNameTok := b.advance()
		valueName = &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     valueNameTok.Value,
		}
	}
	
	b.expect(lexer.RPAREN)
	
	node := &ast.Mapping{
		BaseNode:  ast.BaseNode{Type: ast.NodeMapping},
		KeyType:   keyType,
		KeyName:   keyName,
		ValueType: valueType,
		ValueName: valueName,
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseFunctionTypeName() *ast.FunctionTypeName {
	startTok := b.advance() // function
	
	node := &ast.FunctionTypeName{
		BaseNode:       ast.BaseNode{Type: ast.NodeFunctionTypeName},
		ParameterTypes: make([]*ast.VariableDeclaration, 0),
	}
	
	// Parameters
	node.ParameterTypes = b.parseParameterList()
	
	// Visibility and state mutability
	for {
		if b.check(lexer.INTERNAL) {
			b.advance()
			node.Visibility = "internal"
		} else if b.check(lexer.EXTERNAL) {
			b.advance()
			node.Visibility = "external"
		} else if b.check(lexer.PURE) {
			b.advance()
			node.StateMutability = "pure"
		} else if b.check(lexer.VIEW) {
			b.advance()
			node.StateMutability = "view"
		} else if b.check(lexer.PAYABLE) {
			b.advance()
			node.StateMutability = "payable"
		} else {
			break
		}
	}
	
	// Return types
	if b.check(lexer.RETURNS) {
		b.advance() // returns
		node.ReturnTypes = b.parseParameterList()
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseParameterList() []*ast.VariableDeclaration {
	params := make([]*ast.VariableDeclaration, 0)
	
	b.expect(lexer.LPAREN)
	
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		param := b.parseParameter()
		params = append(params, param)
		
		if !b.check(lexer.RPAREN) {
			b.expect(lexer.COMMA)
		}
	}
	
	b.expect(lexer.RPAREN)
	return params
}

func (b *Builder) parseParameter() *ast.VariableDeclaration {
	startTok := b.peek()
	
	typeName := b.parseTypeName()
	
	node := &ast.VariableDeclaration{
		BaseNode: ast.BaseNode{Type: ast.NodeVariableDeclaration},
		TypeName: typeName,
	}
	
	// Storage location
	if b.check(lexer.MEMORY) || b.check(lexer.STORAGE) || b.check(lexer.CALLDATA) {
		node.StorageLocation = b.advance().Value
	}
	
	// Name (optional) - can be identifier or contextual keyword like 'from'
	if b.check(lexer.IDENTIFIER) || b.isContextualKeyword() {
		nameTok := b.advance()
		node.Name = nameTok.Value
		node.Identifier = &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     nameTok.Value,
		}
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseVariableDeclaration() *ast.VariableDeclaration {
	startTok := b.peek()
	
	typeName := b.parseTypeName()
	
	node := &ast.VariableDeclaration{
		BaseNode: ast.BaseNode{Type: ast.NodeVariableDeclaration},
		TypeName: typeName,
	}
	
	// Storage location
	if b.check(lexer.MEMORY) || b.check(lexer.STORAGE) || b.check(lexer.CALLDATA) {
		node.StorageLocation = b.advance().Value
	}
	
	// Name - can be identifier or contextual keyword like 'from'
	if b.check(lexer.IDENTIFIER) || b.isContextualKeyword() {
		nameTok := b.advance()
		node.Name = nameTok.Value
		node.Identifier = &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     nameTok.Value,
		}
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseStateVariableDeclaration() *ast.StateVariableDeclaration {
	startTok := b.peek()
	
	typeName := b.parseTypeName()
	
	varDecl := &ast.VariableDeclaration{
		BaseNode:   ast.BaseNode{Type: ast.NodeVariableDeclaration},
		TypeName:   typeName,
		IsStateVar: true,
	}
	
	// Modifiers
	for {
		if b.check(lexer.PUBLIC) {
			b.advance()
			varDecl.Visibility = "public"
		} else if b.check(lexer.PRIVATE) {
			b.advance()
			varDecl.Visibility = "private"
		} else if b.check(lexer.INTERNAL) {
			b.advance()
			varDecl.Visibility = "internal"
		} else if b.check(lexer.CONSTANT) {
			b.advance()
			varDecl.IsDeclaredConst = true
		} else if b.check(lexer.IMMUTABLE) {
			b.advance()
			varDecl.IsImmutable = true
		} else if b.check(lexer.TRANSIENT) {
			// Transient storage location (Solidity 0.8.24+)
			b.advance()
			varDecl.StorageLocation = "transient"
		} else if b.check(lexer.OVERRIDE) {
			b.advance()
			if b.check(lexer.LPAREN) {
				b.advance()
				for !b.check(lexer.RPAREN) && !b.isAtEnd() {
					override := b.parseUserDefinedTypeName()
					varDecl.Override = append(varDecl.Override, override)
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
	
	// Name
	nameTok := b.expect(lexer.IDENTIFIER)
	varDecl.Name = nameTok.Value
	varDecl.Identifier = &ast.Identifier{
		BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
		Name:     nameTok.Value,
	}
	
	node := &ast.StateVariableDeclaration{
		BaseNode:  ast.BaseNode{Type: ast.NodeStateVariableDeclaration},
		Variables: []*ast.VariableDeclaration{varDecl},
	}
	
	// Initial value
	if b.check(lexer.ASSIGN) {
		b.advance() // =
		node.InitialValue = b.parseExpression()
	}
	
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	b.setLocation(varDecl, startTok, endTok)
	return node
}

func (b *Builder) parseConstantVariableDeclaration() *ast.StateVariableDeclaration {
	return b.parseStateVariableDeclaration()
}

func (b *Builder) parseStructDefinition() *ast.StructDefinition {
	startTok := b.advance() // struct
	nameTok := b.expect(lexer.IDENTIFIER)
	
	node := &ast.StructDefinition{
		BaseNode: ast.BaseNode{Type: ast.NodeStructDefinition},
		Name:     nameTok.Value,
		Members:  make([]*ast.VariableDeclaration, 0),
	}
	
	b.expect(lexer.LBRACE)
	
	for !b.check(lexer.RBRACE) && !b.isAtEnd() {
		typeName := b.parseTypeName()
		memberNameTok := b.expect(lexer.IDENTIFIER)
		b.expect(lexer.SEMICOLON)
		
		member := &ast.VariableDeclaration{
			BaseNode: ast.BaseNode{Type: ast.NodeVariableDeclaration},
			TypeName: typeName,
			Name:     memberNameTok.Value,
			Identifier: &ast.Identifier{
				BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
				Name:     memberNameTok.Value,
			},
		}
		node.Members = append(node.Members, member)
	}
	
	endTok := b.expect(lexer.RBRACE)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseEnumDefinition() *ast.EnumDefinition {
	startTok := b.advance() // enum
	nameTok := b.expect(lexer.IDENTIFIER)
	
	node := &ast.EnumDefinition{
		BaseNode: ast.BaseNode{Type: ast.NodeEnumDefinition},
		Name:     nameTok.Value,
		Members:  make([]*ast.EnumValue, 0),
	}
	
	b.expect(lexer.LBRACE)
	
	for !b.check(lexer.RBRACE) && !b.isAtEnd() {
		valueTok := b.expect(lexer.IDENTIFIER)
		member := &ast.EnumValue{
			BaseNode: ast.BaseNode{Type: ast.NodeEnumValue},
			Name:     valueTok.Value,
		}
		node.Members = append(node.Members, member)
		
		if !b.check(lexer.RBRACE) {
			b.expect(lexer.COMMA)
		}
	}
	
	endTok := b.expect(lexer.RBRACE)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseEventDefinition() *ast.EventDefinition {
	startTok := b.advance() // event
	nameTok := b.expect(lexer.IDENTIFIER)
	
	node := &ast.EventDefinition{
		BaseNode:   ast.BaseNode{Type: ast.NodeEventDefinition},
		Name:       nameTok.Value,
		Parameters: make([]*ast.VariableDeclaration, 0),
	}
	
	b.expect(lexer.LPAREN)
	
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		typeName := b.parseTypeName()
		
		param := &ast.VariableDeclaration{
			BaseNode: ast.BaseNode{Type: ast.NodeVariableDeclaration},
			TypeName: typeName,
		}
		
		// indexed
		if b.check(lexer.INDEXED) {
			b.advance()
			param.IsIndexed = true
		}
		
		// name (optional) - can be identifier or contextual keyword like 'from'
		if b.check(lexer.IDENTIFIER) || b.isContextualKeyword() {
			paramNameTok := b.advance()
			param.Name = paramNameTok.Value
			param.Identifier = &ast.Identifier{
				BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
				Name:     paramNameTok.Value,
			}
		}
		
		node.Parameters = append(node.Parameters, param)
		
		if b.check(lexer.COMMA) {
			b.advance()
		} else if !b.check(lexer.RPAREN) {
			// Unexpected token, skip it to avoid infinite loop
			b.advance()
		}
	}
	
	b.expect(lexer.RPAREN)
	
	// anonymous
	if b.check(lexer.ANONYMOUS) {
		b.advance()
		node.IsAnonymous = true
	}
	
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseErrorDefinition() *ast.ErrorDefinition {
	startTok := b.advance() // error
	nameTok := b.expect(lexer.IDENTIFIER)
	
	node := &ast.ErrorDefinition{
		BaseNode:   ast.BaseNode{Type: ast.NodeErrorDefinition},
		Name:       nameTok.Value,
		Parameters: make([]*ast.VariableDeclaration, 0),
	}
	
	b.expect(lexer.LPAREN)
	
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		typeName := b.parseTypeName()
		
		param := &ast.VariableDeclaration{
			BaseNode: ast.BaseNode{Type: ast.NodeVariableDeclaration},
			TypeName: typeName,
		}
		
		// name (optional)
		if b.check(lexer.IDENTIFIER) {
			paramNameTok := b.advance()
			param.Name = paramNameTok.Value
			param.Identifier = &ast.Identifier{
				BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
				Name:     paramNameTok.Value,
			}
		}
		
		node.Parameters = append(node.Parameters, param)
		
		if !b.check(lexer.RPAREN) {
			b.expect(lexer.COMMA)
		}
	}
	
	b.expect(lexer.RPAREN)
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseUsingDirective() *ast.UsingForDeclaration {
	startTok := b.advance() // using
	
	node := &ast.UsingForDeclaration{
		BaseNode: ast.BaseNode{Type: ast.NodeUsingForDeclaration},
	}
	
	if b.check(lexer.LBRACE) {
		// using { func1, func2 } for Type
		b.advance() // {
		for !b.check(lexer.RBRACE) && !b.isAtEnd() {
			funcTok := b.expect(lexer.IDENTIFIER)
			node.Functions = append(node.Functions, funcTok.Value)
			
			// Check for operator
			if b.check(lexer.AS) {
				b.advance() // as
				opTok := b.advance()
				node.Operators = append(node.Operators, opTok.Value)
			}
			
			if !b.check(lexer.RBRACE) {
				b.expect(lexer.COMMA)
			}
		}
		b.expect(lexer.RBRACE)
	} else {
		// using Library for Type
		libName := b.parseUserDefinedTypeName()
		node.LibraryName = libName.NamePath
	}
	
	b.expectKeyword("for")
	
	if b.check(lexer.MUL) {
		b.advance() // *
	} else {
		node.TypeName = b.parseTypeName()
	}
	
	// global
	if b.check(lexer.GLOBAL) {
		b.advance()
		node.IsGlobal = true
	}
	
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseUserDefinedValueTypeDefinition() *ast.UserDefinedValueTypeDefinition {
	startTok := b.advance() // type
	nameTok := b.expect(lexer.IDENTIFIER)
	b.expect(lexer.IS)
	
	underlyingType := b.parseElementaryTypeName()
	
	endTok := b.expect(lexer.SEMICOLON)
	
	node := &ast.UserDefinedValueTypeDefinition{
		BaseNode:       ast.BaseNode{Type: ast.NodeUserDefinedValueTypeDefinition},
		Name:           nameTok.Value,
		UnderlyingType: underlyingType,
	}
	
	b.setLocation(node, startTok, endTok)
	return node
}

