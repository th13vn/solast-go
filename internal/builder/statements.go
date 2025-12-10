package builder

import (
	"github.com/th13vn/solast-go/internal/lexer"
	"github.com/th13vn/solast-go/pkg/ast"
)

func (b *Builder) parseBlock() *ast.Block {
	startTok := b.expect(lexer.LBRACE)
	
	node := &ast.Block{
		BaseNode:   ast.BaseNode{Type: ast.NodeBlock},
		Statements: make([]ast.Node, 0),
	}
	
	for !b.check(lexer.RBRACE) && !b.isAtEnd() {
		stmt := b.parseStatement()
		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
		}
	}
	
	endTok := b.expect(lexer.RBRACE)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseStatement() ast.Node {
	tok := b.peek()
	
	switch tok.Type {
	case lexer.LBRACE:
		return b.parseBlock()
	case lexer.IF:
		return b.parseIfStatement()
	case lexer.FOR:
		return b.parseForStatement()
	case lexer.WHILE:
		return b.parseWhileStatement()
	case lexer.DO:
		return b.parseDoWhileStatement()
	case lexer.CONTINUE:
		return b.parseContinueStatement()
	case lexer.BREAK:
		return b.parseBreakStatement()
	case lexer.RETURN:
		return b.parseReturnStatement()
	case lexer.EMIT:
		return b.parseEmitStatement()
	case lexer.REVERT:
		return b.parseRevertStatement()
	case lexer.TRY:
		return b.parseTryStatement()
	case lexer.ASSEMBLY:
		return b.parseAssemblyStatement()
	case lexer.UNCHECKED:
		return b.parseUncheckedBlock()
	default:
		// Variable declaration or expression statement
		// Need lookahead to distinguish:
		// - "Type varName;" (variable declaration)
		// - "func();" (expression statement)
		if b.looksLikeVariableDeclaration() {
			return b.parseVariableDeclarationStatement()
		}
		if b.check(lexer.LPAREN) {
			// Could be tuple declaration
			return b.parseTupleVariableDeclarationOrExpression()
		}
		return b.parseExpressionStatement()
	}
}

// looksLikeVariableDeclaration uses lookahead to determine if current position
// looks like a variable declaration (Type name) vs expression statement (func())
func (b *Builder) looksLikeVariableDeclaration() bool {
	// Elementary types are definitely type names
	if b.isElementaryTypeName() {
		return true
	}
	
	// mapping and function type keywords
	if b.check(lexer.MAPPING) || b.check(lexer.FUNCTION) {
		return true
	}
	
	// For identifiers, we need to look ahead
	if !b.check(lexer.IDENTIFIER) {
		return false
	}
	
	// Save position for backtracking
	savedPos := b.pos
	defer func() { b.pos = savedPos }()
	
	// Skip the identifier (potential type name)
	b.advance()
	
	// Skip array dimensions like [10] or []
	for b.check(lexer.LBRACK) {
		b.advance() // [
		for !b.check(lexer.RBRACK) && !b.isAtEnd() {
			b.advance()
		}
		if b.check(lexer.RBRACK) {
			b.advance() // ]
		}
	}
	
	// Skip type path like A.B.C
	for b.check(lexer.PERIOD) {
		b.advance() // .
		if b.check(lexer.IDENTIFIER) {
			b.advance()
		}
	}
	
	// Skip storage location (memory, storage, calldata)
	if b.check(lexer.MEMORY) || b.check(lexer.STORAGE) || b.check(lexer.CALLDATA) {
		b.advance()
	}
	
	// If followed by identifier, it's a variable declaration
	// e.g., "uint256 x" or "MyType myVar"
	if b.check(lexer.IDENTIFIER) {
		return true
	}
	
	// Otherwise it's likely an expression (function call, etc.)
	return false
}

func (b *Builder) parseIfStatement() *ast.IfStatement {
	startTok := b.advance() // if
	
	b.expect(lexer.LPAREN)
	condition := b.parseExpression()
	b.expect(lexer.RPAREN)
	
	trueBody := b.parseStatement()
	
	node := &ast.IfStatement{
		BaseNode:  ast.BaseNode{Type: ast.NodeIfStatement},
		Condition: condition,
		TrueBody:  trueBody,
	}
	
	if b.check(lexer.ELSE) {
		b.advance() // else
		node.FalseBody = b.parseStatement()
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseForStatement() *ast.ForStatement {
	startTok := b.advance() // for
	
	b.expect(lexer.LPAREN)
	
	node := &ast.ForStatement{
		BaseNode: ast.BaseNode{Type: ast.NodeForStatement},
	}
	
	// Init
	if !b.check(lexer.SEMICOLON) {
		if b.isTypeName() {
			node.InitExpression = b.parseVariableDeclarationStatement()
		} else {
			node.InitExpression = b.parseExpressionStatement()
		}
	} else {
		b.advance() // ;
	}
	
	// Condition
	if !b.check(lexer.SEMICOLON) {
		node.ConditionExpression = b.parseExpression()
	}
	b.expect(lexer.SEMICOLON)
	
	// Loop expression
	if !b.check(lexer.RPAREN) {
		node.LoopExpression = b.parseExpression()
	}
	b.expect(lexer.RPAREN)
	
	node.Body = b.parseStatement()
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseWhileStatement() *ast.WhileStatement {
	startTok := b.advance() // while
	
	b.expect(lexer.LPAREN)
	condition := b.parseExpression()
	b.expect(lexer.RPAREN)
	body := b.parseStatement()
	
	node := &ast.WhileStatement{
		BaseNode:  ast.BaseNode{Type: ast.NodeWhileStatement},
		Condition: condition,
		Body:      body,
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseDoWhileStatement() *ast.DoWhileStatement {
	startTok := b.advance() // do
	
	body := b.parseStatement()
	b.expect(lexer.WHILE)
	b.expect(lexer.LPAREN)
	condition := b.parseExpression()
	b.expect(lexer.RPAREN)
	b.expect(lexer.SEMICOLON)
	
	node := &ast.DoWhileStatement{
		BaseNode:  ast.BaseNode{Type: ast.NodeDoWhileStatement},
		Condition: condition,
		Body:      body,
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseContinueStatement() *ast.ContinueStatement {
	startTok := b.advance() // continue
	endTok := b.expect(lexer.SEMICOLON)
	
	node := &ast.ContinueStatement{
		BaseNode: ast.BaseNode{Type: ast.NodeContinueStatement},
	}
	
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseBreakStatement() *ast.BreakStatement {
	startTok := b.advance() // break
	endTok := b.expect(lexer.SEMICOLON)
	
	node := &ast.BreakStatement{
		BaseNode: ast.BaseNode{Type: ast.NodeBreakStatement},
	}
	
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseReturnStatement() *ast.ReturnStatement {
	startTok := b.advance() // return
	
	node := &ast.ReturnStatement{
		BaseNode: ast.BaseNode{Type: ast.NodeReturnStatement},
	}
	
	if !b.check(lexer.SEMICOLON) {
		node.Expression = b.parseExpression()
	}
	
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseEmitStatement() *ast.EmitStatement {
	startTok := b.advance() // emit
	
	eventCall := b.parseExpression()
	endTok := b.expect(lexer.SEMICOLON)
	
	node := &ast.EmitStatement{
		BaseNode:  ast.BaseNode{Type: ast.NodeEmitStatement},
		EventCall: eventCall,
	}
	
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseRevertStatement() *ast.RevertStatement {
	startTok := b.advance() // revert
	
	node := &ast.RevertStatement{
		BaseNode: ast.BaseNode{Type: ast.NodeRevertStatement},
	}
	
	if !b.check(lexer.SEMICOLON) {
		node.RevertCall = b.parseExpression()
	}
	
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseTryStatement() *ast.TryStatement {
	startTok := b.advance() // try
	
	// Parse try expression but don't consume { which is the body block
	expr := b.parseTryExpression()
	
	node := &ast.TryStatement{
		BaseNode:   ast.BaseNode{Type: ast.NodeTryStatement},
		Expression: expr,
	}
	
	if b.check(lexer.RETURNS) {
		b.advance() // returns
		node.ReturnParameters = b.parseParameterList()
	}
	
	node.Body = b.parseBlock()
	
	// Catch clauses
	for b.check(lexer.CATCH) {
		clause := b.parseCatchClause()
		node.CatchClauses = append(node.CatchClauses, clause)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

// parseTryExpression parses the expression after 'try' keyword
// It's similar to parseExpression but doesn't consume { as FunctionCallOptions
// since { is the try block body
func (b *Builder) parseTryExpression() ast.Node {
	return b.parseTryCallMemberIndex()
}

func (b *Builder) parseTryCallMemberIndex() ast.Node {
	expr := b.parsePrimary()
	
	for {
		if b.check(lexer.PERIOD) {
			b.advance() // .
			memberTok := b.advance()
			expr = &ast.MemberAccess{
				BaseNode:   ast.BaseNode{Type: ast.NodeMemberAccess},
				Expression: expr,
				MemberName: memberTok.Value,
			}
		} else if b.check(lexer.LBRACK) {
			b.advance() // [
			
			var indexStart ast.Node
			var indexEnd ast.Node
			isRange := false
			
			if !b.check(lexer.COLON) && !b.check(lexer.RBRACK) {
				indexStart = b.parseExpression()
			}
			
			if b.check(lexer.COLON) {
				isRange = true
				b.advance() // :
				if !b.check(lexer.RBRACK) {
					indexEnd = b.parseExpression()
				}
			}
			
			b.expect(lexer.RBRACK)
			
			if isRange {
				expr = &ast.IndexRangeAccess{
					BaseNode:   ast.BaseNode{Type: ast.NodeIndexRangeAccess},
					Base:       expr,
					IndexStart: indexStart,
					IndexEnd:   indexEnd,
				}
			} else {
				expr = &ast.IndexAccess{
					BaseNode: ast.BaseNode{Type: ast.NodeIndexAccess},
					Base:     expr,
					Index:    indexStart,
				}
			}
		} else if b.check(lexer.LPAREN) {
			expr = b.parseFunctionCall(expr)
		} else {
			// Don't parse { as FunctionCallOptions in try context
			break
		}
	}
	
	return expr
}

func (b *Builder) parseCatchClause() *ast.CatchClause {
	startTok := b.advance() // catch
	
	node := &ast.CatchClause{
		BaseNode: ast.BaseNode{Type: ast.NodeCatchClause},
	}
	
	// Optional catch identifier and parameters
	if b.check(lexer.IDENTIFIER) {
		kindTok := b.advance()
		node.Kind = kindTok.Value
		if kindTok.Value == "Error" {
			node.IsReasonStringType = true
		}
	}
	
	if b.check(lexer.LPAREN) {
		node.Parameters = b.parseParameterList()
	}
	
	node.Body = b.parseBlock()
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblyStatement() *ast.InlineAssembly {
	startTok := b.advance() // assembly
	
	node := &ast.InlineAssembly{
		BaseNode: ast.BaseNode{Type: ast.NodeInlineAssembly},
	}
	
	// Optional dialect string
	if b.check(lexer.STRING) {
		node.Language = b.advance().Value
	}
	
	node.Body = b.parseAssemblyBlock()
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblyBlock() *ast.AssemblyBlock {
	startTok := b.expect(lexer.LBRACE)
	
	node := &ast.AssemblyBlock{
		BaseNode:   ast.BaseNode{Type: ast.NodeAssemblyBlock},
		Operations: make([]ast.Node, 0),
	}
	
	for !b.check(lexer.RBRACE) && !b.isAtEnd() {
		op := b.parseAssemblyStatement_()
		if op != nil {
			node.Operations = append(node.Operations, op)
		}
	}
	
	endTok := b.expect(lexer.RBRACE)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseAssemblyStatement_() ast.Node {
	tok := b.peek()
	
	switch tok.Type {
	case lexer.LBRACE:
		return b.parseAssemblyBlock()
	case lexer.LET:
		return b.parseAssemblyLocalDefinition()
	case lexer.IF:
		return b.parseAssemblyIf()
	case lexer.FOR:
		return b.parseAssemblyFor()
	case lexer.SWITCH:
		return b.parseAssemblySwitch()
	case lexer.FUNCTION:
		return b.parseAssemblyFunctionDefinition()
	case lexer.IDENTIFIER:
		return b.parseAssemblyExpressionOrAssignment()
	default:
		if tok.Type == lexer.RBRACE {
			return nil
		}
		b.advance()
		return nil
	}
}

func (b *Builder) parseAssemblyLocalDefinition() *ast.AssemblyLocalDefinition {
	startTok := b.advance() // let
	
	node := &ast.AssemblyLocalDefinition{
		BaseNode: ast.BaseNode{Type: ast.NodeAssemblyLocalDefinition},
		Names:    make([]*ast.Identifier, 0),
	}
	
	// Parse identifier list
	for {
		nameTok := b.expect(lexer.IDENTIFIER)
		node.Names = append(node.Names, &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     nameTok.Value,
		})
		if !b.check(lexer.COMMA) {
			break
		}
		b.advance() // ,
	}
	
	if b.check(lexer.ASSIGN) {
		b.advance() // :=
		node.Expression = b.parseAssemblyExpression()
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblyIf() *ast.AssemblyIf {
	startTok := b.advance() // if
	
	condition := b.parseAssemblyExpression()
	body := b.parseAssemblyBlock()
	
	node := &ast.AssemblyIf{
		BaseNode:  ast.BaseNode{Type: ast.NodeAssemblyIf},
		Condition: condition,
		Body:      body,
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblyFor() *ast.AssemblyFor {
	startTok := b.advance() // for
	
	pre := b.parseAssemblyBlock()
	condition := b.parseAssemblyExpression()
	post := b.parseAssemblyBlock()
	body := b.parseAssemblyBlock()
	
	node := &ast.AssemblyFor{
		BaseNode:  ast.BaseNode{Type: ast.NodeAssemblyFor},
		Pre:       pre,
		Condition: condition,
		Post:      post,
		Body:      body,
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblySwitch() *ast.AssemblySwitch {
	startTok := b.advance() // switch
	
	expr := b.parseAssemblyExpression()
	
	node := &ast.AssemblySwitch{
		BaseNode:   ast.BaseNode{Type: ast.NodeAssemblySwitch},
		Expression: expr,
		Cases:      make([]*ast.AssemblyCase, 0),
	}
	
	for b.check(lexer.CASE) || b.check(lexer.DEFAULT) {
		isDefault := b.check(lexer.DEFAULT)
		b.advance() // case/default
		
		caseNode := &ast.AssemblyCase{
			BaseNode: ast.BaseNode{Type: ast.NodeAssemblyCase},
			Default:  isDefault,
		}
		
		if !isDefault {
			caseNode.Value = b.parseAssemblyLiteral()
		}
		
		caseNode.Body = b.parseAssemblyBlock()
		node.Cases = append(node.Cases, caseNode)
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblyFunctionDefinition() *ast.AssemblyFunctionDefinition {
	startTok := b.advance() // function
	
	nameTok := b.expect(lexer.IDENTIFIER)
	
	node := &ast.AssemblyFunctionDefinition{
		BaseNode: ast.BaseNode{Type: ast.NodeAssemblyFunctionDefinition},
		Name:     nameTok.Value,
	}
	
	b.expect(lexer.LPAREN)
	// Arguments
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		argTok := b.expect(lexer.IDENTIFIER)
		node.Arguments = append(node.Arguments, &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     argTok.Value,
		})
		if !b.check(lexer.RPAREN) {
			b.expect(lexer.COMMA)
		}
	}
	b.expect(lexer.RPAREN)
	
	// Return values
	if b.check(lexer.ARROW) {
		b.advance() // ->
		for {
			retTok := b.expect(lexer.IDENTIFIER)
			node.ReturnArguments = append(node.ReturnArguments, &ast.Identifier{
				BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
				Name:     retTok.Value,
			})
			if !b.check(lexer.COMMA) {
				break
			}
			b.advance() // ,
		}
	}
	
	node.Body = b.parseAssemblyBlock()
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblyExpressionOrAssignment() ast.Node {
	startTok := b.peek()
	
	// Parse identifier(s)
	var names []*ast.Identifier
	for {
		nameTok := b.expect(lexer.IDENTIFIER)
		names = append(names, &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     nameTok.Value,
		})
		if !b.check(lexer.COMMA) {
			break
		}
		b.advance() // ,
	}
	
	// Check for assignment
	if b.check(lexer.ASSIGN) {
		b.advance() // :=
		expr := b.parseAssemblyExpression()
		
		node := &ast.AssemblyAssignment{
			BaseNode:   ast.BaseNode{Type: ast.NodeAssemblyAssignment},
			Names:      names,
			Expression: expr,
		}
		b.setLocation(node, startTok, b.previous())
		return node
	}
	
	// Function call or just identifier
	if len(names) == 1 {
		if b.check(lexer.LPAREN) {
			return b.parseAssemblyCall(names[0].Name, startTok)
		}
		return names[0]
	}
	
	// Return first identifier if no assignment
	return names[0]
}

func (b *Builder) parseAssemblyCall(name string, startTok lexer.Token) *ast.AssemblyCall {
	b.expect(lexer.LPAREN)
	
	node := &ast.AssemblyCall{
		BaseNode:     ast.BaseNode{Type: ast.NodeAssemblyCall},
		FunctionName: name,
		Arguments:    make([]ast.Node, 0),
	}
	
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		arg := b.parseAssemblyExpression()
		node.Arguments = append(node.Arguments, arg)
		if !b.check(lexer.RPAREN) {
			b.expect(lexer.COMMA)
		}
	}
	b.expect(lexer.RPAREN)
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseAssemblyExpression() ast.Node {
	tok := b.peek()
	
	if tok.Type == lexer.IDENTIFIER {
		startTok := b.advance()
		if b.check(lexer.LPAREN) {
			return b.parseAssemblyCall(startTok.Value, startTok)
		}
		node := &ast.AssemblyIdentifier{
			BaseNode: ast.BaseNode{Type: ast.NodeAssemblyIdentifier},
			Name:     startTok.Value,
		}
		b.setLocation(node, startTok, startTok)
		return node
	}
	
	return b.parseAssemblyLiteral()
}

func (b *Builder) parseAssemblyLiteral() ast.Node {
	tok := b.advance()
	
	node := &ast.AssemblyLiteral{
		BaseNode: ast.BaseNode{Type: ast.NodeAssemblyLiteral},
		Value:    tok.Value,
	}
	
	switch tok.Type {
	case lexer.NUMBER, lexer.HEX_NUMBER:
		node.Kind = "number"
	case lexer.STRING:
		node.Kind = "string"
	case lexer.TRUE, lexer.FALSE:
		node.Kind = "boolean"
	default:
		node.Kind = "number"
	}
	
	b.setLocation(node, tok, tok)
	return node
}

func (b *Builder) parseUncheckedBlock() *ast.UncheckedBlock {
	startTok := b.advance() // unchecked
	
	body := b.parseBlock()
	
	node := &ast.UncheckedBlock{
		BaseNode: ast.BaseNode{Type: ast.NodeUncheckedBlock},
		Body:     body,
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseVariableDeclarationStatement() *ast.VariableDeclarationStatement {
	startTok := b.peek()
	
	node := &ast.VariableDeclarationStatement{
		BaseNode:  ast.BaseNode{Type: ast.NodeVariableDeclarationStatement},
		Variables: make([]*ast.VariableDeclaration, 0),
	}
	
	varDecl := b.parseVariableDeclaration()
	node.Variables = append(node.Variables, varDecl)
	
	if b.check(lexer.ASSIGN) {
		b.advance() // =
		node.InitialValue = b.parseExpression()
	}
	
	endTok := b.expect(lexer.SEMICOLON)
	b.setLocation(node, startTok, endTok)
	return node
}

func (b *Builder) parseTupleVariableDeclarationOrExpression() ast.Node {
	startTok := b.peek()
	
	// Look ahead to determine if this is a tuple declaration
	// This is a simplified version - full implementation would need more lookahead
	
	b.expect(lexer.LPAREN)
	
	// Try to parse as tuple declaration first
	var variables []*ast.VariableDeclaration
	var hasTypes bool
	
	savedPos := b.pos
	
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		if b.check(lexer.COMMA) {
			variables = append(variables, nil)
			b.advance()
			continue
		}
		
		if b.isTypeName() {
			hasTypes = true
			varDecl := b.parseVariableDeclaration()
			variables = append(variables, varDecl)
		} else {
			break
		}
		
		if !b.check(lexer.RPAREN) && !b.check(lexer.COMMA) {
			break
		}
		if b.check(lexer.COMMA) {
			b.advance()
		}
	}
	
	if hasTypes && b.check(lexer.RPAREN) {
		b.expect(lexer.RPAREN)
		
		node := &ast.VariableDeclarationStatement{
			BaseNode:  ast.BaseNode{Type: ast.NodeVariableDeclarationStatement},
			Variables: variables,
		}
		
		if b.check(lexer.ASSIGN) {
			b.advance()
			node.InitialValue = b.parseExpression()
		}
		
		endTok := b.expect(lexer.SEMICOLON)
		b.setLocation(node, startTok, endTok)
		return node
	}
	
	// Not a tuple declaration, restore and parse as expression
	b.pos = savedPos
	return b.parseExpressionStatement()
}

func (b *Builder) parseExpressionStatement() *ast.ExpressionStatement {
	startTok := b.peek()
	
	expr := b.parseExpression()
	endTok := b.expect(lexer.SEMICOLON)
	
	node := &ast.ExpressionStatement{
		BaseNode:   ast.BaseNode{Type: ast.NodeExpressionStatement},
		Expression: expr,
	}
	
	b.setLocation(node, startTok, endTok)
	return node
}

