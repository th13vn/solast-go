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
	// An elementary type name at statement start is usually a declaration
	// (`address x = ...`), but it is also how an elementary type CONVERSION is
	// written: `address(x).call(...)`, `uint256(v).toString()`. Only the
	// conversion can put `(` directly after the type – no variable declaration
	// has that shape – so one token of lookahead separates them.
	//
	// Treating every leading elementary type as a declaration sent
	// `address(x).call(b);` to parseVariableDeclaration, which died on the `(`
	// and dropped the whole statement. The enclosing function then carried no
	// call at all, making it invisible to the call graph, to taint, and to
	// every consumer's detectors. `address(target).call/.delegatecall/
	// .staticcall` and the OpenZeppelin `Address` library idiom are written
	// this way throughout production Solidity.
	if b.isElementaryTypeName() {
		return !b.nextTokenIs(lexer.LPAREN)
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
	
	// Skip type path like A.B.C. The path comes BEFORE any array dimensions
	// (`Storage.Assimilator[] memory x`); skipping brackets first left the
	// `.Assimilator` unread, so the lookahead met `[` instead of a location
	// or name and sent the declaration down the expression path.
	for b.check(lexer.PERIOD) {
		b.advance() // .
		if b.check(lexer.IDENTIFIER) || b.isContextualKeyword() {
			b.advance()
		}
	}
	
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
	
	// Skip storage location (memory, storage, calldata)
	if b.check(lexer.MEMORY) || b.check(lexer.STORAGE) || b.check(lexer.CALLDATA) {
		b.advance()
	}
	
	// If followed by a declaration name, it's a variable declaration
	// e.g., "uint256 x" or "MyType myVar". The name may legally be a contextual
	// keyword — `UserInfo storage from = ...` is ordinary Solidity — and
	// parseVariableDeclaration already accepts one, so this lookahead must too.
	// Otherwise the statement falls through to the expression path, dies on the
	// storage location, and desyncs the tolerant parser for the rest of the file.
	if b.check(lexer.IDENTIFIER) || b.isContextualKeyword() {
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
	
	// Init. Use the same declaration lookahead as a statement: isTypeName()
	// is true for ANY identifier, so `for (i = 0; ...)` became a nameless
	// declaration of type `i` and `for (y >>= 1; ...)` failed on the `>>=`.
	if !b.check(lexer.SEMICOLON) {
		if b.looksLikeVariableDeclaration() {
			node.InitExpression = b.parseVariableDeclarationStatement()
		} else if b.check(lexer.LPAREN) {
			node.InitExpression = b.parseTupleVariableDeclarationOrExpression()
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
	
	// Optional assembly flags: assembly ("memory-safe") { ... }
	node.Flags = b.parseAssemblyFlags()

	node.Body = b.parseAssemblyBlock()

	b.setLocation(node, startTok, b.previous())
	return node
}

// parseAssemblyFlags consumes the optional assemblyFlags group that may follow
// the assembly keyword and its optional dialect string:
//
//	assembly ("memory-safe") { ... }
//
// Grammar (SolidityParser.g4 assemblyFlags): AssemblyBlockLParen
// AssemblyFlagString (AssemblyBlockComma AssemblyFlagString)*
// AssemblyBlockRParen. The flags carry no meaning for this parser and are
// recorded only so consumers can see them, but they MUST be consumed: leaving
// the '(' for parseAssemblyBlock desyncs tolerant parsing and silently drops
// every declaration after the block. Returns nil when no flag group is present.
func (b *Builder) parseAssemblyFlags() []string {
	if !b.check(lexer.LPAREN) {
		return nil
	}
	b.advance() // (

	var flags []string
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		if !b.check(lexer.STRING) {
			// Not a flag list after all. Report it and stop consuming so
			// synchronize() can recover instead of eating the block.
			b.expect(lexer.STRING)
			break
		}
		flags = append(flags, b.advance().Value)
		if !b.check(lexer.COMMA) {
			break
		}
		b.advance() // ,
	}

	b.expect(lexer.RPAREN)
	return flags
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
	default:
		if tok.Type == lexer.RBRACE {
			return nil
		}
		// An identifier, or a builtin spelled like a Solidity keyword
		// (`return(0, 0x20)`, `revert(p, n)`), starts a call or assignment.
		if b.isYulIdentifier() {
			return b.parseAssemblyExpressionOrAssignment()
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
		nameTok := b.expectYulIdentifier()
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
	
	nameTok := b.expectYulIdentifier()
	
	node := &ast.AssemblyFunctionDefinition{
		BaseNode: ast.BaseNode{Type: ast.NodeAssemblyFunctionDefinition},
		Name:     nameTok.Value,
	}
	
	b.expect(lexer.LPAREN)
	// Arguments
	for !b.check(lexer.RPAREN) && !b.isAtEnd() {
		argTok := b.expectYulIdentifier()
		node.Arguments = append(node.Arguments, &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     argTok.Value,
		})
		if !b.check(lexer.RPAREN) {
			b.expect(lexer.COMMA)
		}
	}
	b.expect(lexer.RPAREN)
	
	// Return values. The lexer produces RIGHT_ARROW for `->`; checking for
	// ARROW (`=>`) made every Yul function with return variables fail at the
	// arrow and swallow the rest of the file.
	if b.check(lexer.RIGHT_ARROW) {
		b.advance() // ->
		for {
			retTok := b.expectYulIdentifier()
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
	
	// Parse identifier(s). Each may be a dotted yulPath referring to a
	// declaration outside the block (`slot.offset`, `x.slot`).
	var names []*ast.Identifier
	for {
		nameTok := b.expectYulIdentifier()
		names = append(names, &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     b.parseAssemblyPathSuffix(nameTok.Value),
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

// parseAssemblyPathSuffix consumes the `.member` segments of a Yul path whose
// head identifier has already been read, and returns the full dotted name.
//
// Grammar (SolidityParser.g4): yulPath: (YulIdentifier | YulEVMBuiltin)
// (YulPeriod (YulIdentifier | YulEVMBuiltin))*. Only dot-free identifiers can be
// DECLARED inside assembly, but a path may REFER to a declaration outside the
// block: calldata slice members (`sig.offset`, `sig.length`) and storage-pointer
// members (`x.slot`, `x.offset`).
//
// The dotted text is kept in the single Name field rather than introducing a new
// node type, so the public AST shape stays backward compatible. Leaving the '.'
// unconsumed desynchronized the block and dropped the rest of the file.
func (b *Builder) parseAssemblyPathSuffix(head string) string {
	name := head
	for b.check(lexer.PERIOD) {
		b.advance() // .
		// Yul member names may be builtins (`offset`, `length`, `slot`) or
		// Solidity keywords (`g.address` on an external function pointer),
		// which the lexer classifies as keywords rather than identifiers.
		if b.isYulIdentifier() {
			name += "." + b.advance().Value
			continue
		}
		b.expect(lexer.IDENTIFIER)
		break
	}
	return name
}

func (b *Builder) parseAssemblyExpression() ast.Node {
	if b.isYulIdentifier() {
		startTok := b.advance()
		if b.check(lexer.LPAREN) {
			return b.parseAssemblyCall(startTok.Value, startTok)
		}
		node := &ast.AssemblyIdentifier{
			BaseNode: ast.BaseNode{Type: ast.NodeAssemblyIdentifier},
			Name:     b.parseAssemblyPathSuffix(startTok.Value),
		}
		b.setLocation(node, startTok, b.previous())
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
	
	// Backtrack point for the expression fallback: BEFORE the '(' so the
	// expression parser sees the whole tuple. Restoring to after it made
	// `(items[i].ok, items[i].data) = g();` parse as `items[i].ok` followed
	// by a stray ',' and drop the statement. Errors recorded while
	// speculating belong to the abandoned reading, so they are dropped too.
	savedPos := b.pos
	savedErrors := len(b.errors)
	
	b.expect(lexer.LPAREN)
	
	// Try to parse as tuple declaration first
	var variables []*ast.VariableDeclaration
	var hasTypes bool
	
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
	b.errors = b.errors[:savedErrors]
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

