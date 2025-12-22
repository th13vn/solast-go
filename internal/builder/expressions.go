package builder

import (
	"github.com/th13vn/solast-go/internal/lexer"
	"github.com/th13vn/solast-go/pkg/ast"
)

// Expression parsing with operator precedence
// Precedence (lowest to highest):
// 1. Assignment (=, +=, etc.)
// 2. Conditional (?:)
// 3. Logical OR (||)
// 4. Logical AND (&&)
// 5. Equality (==, !=)
// 6. Relational (<, >, <=, >=)
// 7. Bitwise OR (|)
// 8. Bitwise XOR (^)
// 9. Bitwise AND (&)
// 10. Shift (<<, >>, >>>)
// 11. Addition (+, -)
// 12. Multiplication (*, /, %)
// 13. Exponentiation (**)
// 14. Unary (!, ~, -, +, ++, --, delete)
// 15. Postfix (++, --)
// 16. Call, Index, Member access

func (b *Builder) parseExpression() ast.Node {
	return b.parseAssignment()
}

func (b *Builder) parseAssignment() ast.Node {
	left := b.parseTernary()
	
	if b.isAssignmentOperator() {
		op := b.advance().Value
		right := b.parseAssignment()
		
		return &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseTernary() ast.Node {
	condition := b.parseLogicalOr()
	
	if b.check(lexer.QUESTION) {
		b.advance() // ?
		trueExpr := b.parseExpression()
		b.expect(lexer.COLON)
		falseExpr := b.parseTernary()
		
		return &ast.Conditional{
			BaseNode:        ast.BaseNode{Type: ast.NodeConditional},
			Condition:       condition,
			TrueExpression:  trueExpr,
			FalseExpression: falseExpr,
		}
	}
	
	return condition
}

func (b *Builder) parseLogicalOr() ast.Node {
	left := b.parseLogicalAnd()
	
	for b.check(lexer.OR) {
		op := b.advance().Value
		right := b.parseLogicalAnd()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseLogicalAnd() ast.Node {
	left := b.parseEquality()
	
	for b.check(lexer.AND) {
		op := b.advance().Value
		right := b.parseEquality()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseEquality() ast.Node {
	left := b.parseRelational()
	
	for b.check(lexer.EQ) || b.check(lexer.NEQ) {
		op := b.advance().Value
		right := b.parseRelational()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseRelational() ast.Node {
	left := b.parseBitwiseOr()
	
	for b.check(lexer.LT) || b.check(lexer.GT) || b.check(lexer.LTE) || b.check(lexer.GTE) {
		op := b.advance().Value
		right := b.parseBitwiseOr()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseBitwiseOr() ast.Node {
	left := b.parseBitwiseXor()
	
	for b.check(lexer.BIT_OR) {
		op := b.advance().Value
		right := b.parseBitwiseXor()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseBitwiseXor() ast.Node {
	left := b.parseBitwiseAnd()
	
	for b.check(lexer.BIT_XOR) {
		op := b.advance().Value
		right := b.parseBitwiseAnd()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseBitwiseAnd() ast.Node {
	left := b.parseShift()
	
	for b.check(lexer.BIT_AND) {
		op := b.advance().Value
		right := b.parseShift()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseShift() ast.Node {
	left := b.parseAdditive()
	
	for b.check(lexer.SHL) || b.check(lexer.SHR) || b.check(lexer.SAR) {
		op := b.advance().Value
		right := b.parseAdditive()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseAdditive() ast.Node {
	left := b.parseMultiplicative()
	
	for b.check(lexer.ADD) || b.check(lexer.SUB) {
		op := b.advance().Value
		right := b.parseMultiplicative()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseMultiplicative() ast.Node {
	left := b.parseExponentiation()
	
	for b.check(lexer.MUL) || b.check(lexer.DIV) || b.check(lexer.MOD) {
		op := b.advance().Value
		right := b.parseExponentiation()
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseExponentiation() ast.Node {
	left := b.parseUnary()
	
	if b.check(lexer.EXP) {
		op := b.advance().Value
		right := b.parseExponentiation() // Right associative
		left = &ast.BinaryOperation{
			BaseNode: ast.BaseNode{Type: ast.NodeBinaryOperation},
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	
	return left
}

func (b *Builder) parseUnary() ast.Node {
	if b.check(lexer.NOT) || b.check(lexer.BIT_NOT) || b.check(lexer.SUB) || 
	   b.check(lexer.ADD) || b.check(lexer.INC) || b.check(lexer.DEC) || b.check(lexer.DELETE) {
		op := b.advance().Value
		expr := b.parseUnary()
		
		return &ast.UnaryOperation{
			BaseNode:      ast.BaseNode{Type: ast.NodeUnaryOperation},
			Operator:      op,
			SubExpression: expr,
			IsPrefix:      true,
		}
	}
	
	return b.parsePostfix()
}

func (b *Builder) parsePostfix() ast.Node {
	expr := b.parseCallMemberIndex()
	
	for b.check(lexer.INC) || b.check(lexer.DEC) {
		op := b.advance().Value
		expr = &ast.UnaryOperation{
			BaseNode:      ast.BaseNode{Type: ast.NodeUnaryOperation},
			Operator:      op,
			SubExpression: expr,
			IsPrefix:      false,
		}
	}
	
	return expr
}

func (b *Builder) parseCallMemberIndex() ast.Node {
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
		} else if b.check(lexer.LBRACE) {
			// Named arguments for function call options
			expr = b.parseFunctionCallOptions(expr)
		} else {
			break
		}
	}
	
	return expr
}

func (b *Builder) parseFunctionCall(callee ast.Node) *ast.FunctionCall {
	b.expect(lexer.LPAREN)
	
	node := &ast.FunctionCall{
		BaseNode:   ast.BaseNode{Type: ast.NodeFunctionCall},
		Expression: callee,
		Arguments:  make([]ast.Node, 0),
	}
	
	// Check for named arguments
	if b.check(lexer.LBRACE) {
		b.advance() // {
		for !b.check(lexer.RBRACE) && !b.isAtEnd() {
			nameTok := b.expect(lexer.IDENTIFIER)
			b.expect(lexer.COLON)
			value := b.parseExpression()
			
			node.Names = append(node.Names, nameTok.Value)
			node.Arguments = append(node.Arguments, value)
			
			if !b.check(lexer.RBRACE) {
				b.expect(lexer.COMMA)
			}
		}
		b.expect(lexer.RBRACE)
	} else {
		node.Arguments = b.parseExpressionList()
	}
	
	b.expect(lexer.RPAREN)
	return node
}

func (b *Builder) parseFunctionCallOptions(expr ast.Node) *ast.FunctionCallOptions {
	b.expect(lexer.LBRACE)
	
	node := &ast.FunctionCallOptions{
		BaseNode:   ast.BaseNode{Type: ast.NodeFunctionCallOptions},
		Expression: expr,
		Names:      make([]string, 0),
		Options:    make([]ast.Node, 0),
	}
	
	for !b.check(lexer.RBRACE) && !b.isAtEnd() {
		nameTok := b.expect(lexer.IDENTIFIER)
		b.expect(lexer.COLON)
		value := b.parseExpression()
		
		node.Names = append(node.Names, nameTok.Value)
		node.Options = append(node.Options, value)
		
		if !b.check(lexer.RBRACE) {
			b.expect(lexer.COMMA)
		}
	}
	
	b.expect(lexer.RBRACE)
	return node
}

func (b *Builder) parsePrimary() ast.Node {
	tok := b.peek()
	
	switch tok.Type {
	case lexer.IDENTIFIER:
		b.advance()
		return &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     tok.Value,
		}
	
	// Contextual keywords can also be used as identifiers in expressions
	case lexer.FROM, lexer.ERROR, lexer.REVERT, lexer.GLOBAL, lexer.TRANSIENT, lexer.LAYOUT, lexer.AT:
		b.advance()
		return &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     tok.Value,
		}
	
	case lexer.NUMBER:
		b.advance()
		return b.parseNumberLiteral(tok)
	
	case lexer.HEX_NUMBER:
		b.advance()
		return &ast.NumberLiteral{
			BaseNode: ast.BaseNode{Type: ast.NodeNumberLiteral},
			Number:   tok.Value,
		}
	
	case lexer.STRING, lexer.HEX_STRING, lexer.UNICODE_STRING:
		return b.parseStringLiteral()
	
	case lexer.TRUE:
		b.advance()
		return &ast.BooleanLiteral{
			BaseNode: ast.BaseNode{Type: ast.NodeBooleanLiteral},
			Value:    true,
		}
	
	case lexer.FALSE:
		b.advance()
		return &ast.BooleanLiteral{
			BaseNode: ast.BaseNode{Type: ast.NodeBooleanLiteral},
			Value:    false,
		}
	
	case lexer.LPAREN:
		return b.parseTupleOrParenthesized()
	
	case lexer.LBRACK:
		return b.parseArrayLiteral()
	
	case lexer.NEW:
		return b.parseNewExpression()
	
	case lexer.TYPE:
		return b.parseTypeExpression()
	
	case lexer.PAYABLE:
		return b.parsePayableConversion()
	
	// Elementary type names as expressions
	case lexer.ADDRESS, lexer.BOOL, lexer.STRING_TYPE, lexer.BYTES,
		lexer.INT, lexer.UINT, lexer.BYTE, lexer.BYTES_N,
		lexer.FIXED, lexer.UFIXED, lexer.FIXED_N, lexer.UFIXED_N:
		return b.parseElementaryTypeName()
	
	default:
		b.addError("expected expression")
		b.advance()
		return &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     "",
		}
	}
}

func (b *Builder) parseNumberLiteral(tok lexer.Token) *ast.NumberLiteral {
	node := &ast.NumberLiteral{
		BaseNode: ast.BaseNode{Type: ast.NodeNumberLiteral},
		Number:   tok.Value,
	}
	
	// Check for number unit
	if b.checkNumberUnit() {
		node.SubDenomination = b.advance().Value
	}
	
	return node
}

func (b *Builder) parseStringLiteral() ast.Node {
	var parts []string
	var isUnicode bool
	var isHex bool
	
	for b.check(lexer.STRING) || b.check(lexer.UNICODE_STRING) || b.check(lexer.HEX_STRING) {
		tok := b.advance()
		parts = append(parts, tok.Value)
		if tok.Type == lexer.UNICODE_STRING {
			isUnicode = true
		}
		if tok.Type == lexer.HEX_STRING {
			isHex = true
		}
	}
	
	if isHex {
		return &ast.HexLiteral{
			BaseNode: ast.BaseNode{Type: ast.NodeHexLiteral},
			Value:    parts[0],
			Parts:    parts,
		}
	}
	
	return &ast.StringLiteral{
		BaseNode:  ast.BaseNode{Type: ast.NodeStringLiteral},
		Value:     parts[0],
		Parts:     parts,
		IsUnicode: isUnicode,
	}
}

func (b *Builder) parseTupleOrParenthesized() ast.Node {
	startTok := b.advance() // (
	
	// Empty tuple
	if b.check(lexer.RPAREN) {
		b.advance()
		return &ast.TupleExpression{
			BaseNode:   ast.BaseNode{Type: ast.NodeTupleExpression},
			Components: make([]ast.Node, 0),
			IsArray:    false,
		}
	}
	
	// Parse first expression
	var components []ast.Node
	
	// Check for empty component
	if b.check(lexer.COMMA) {
		components = append(components, nil)
	} else {
		expr := b.parseExpression()
		components = append(components, expr)
	}
	
	// If followed by comma, it's a tuple
	if b.check(lexer.COMMA) {
		for b.check(lexer.COMMA) {
			b.advance() // ,
			if b.check(lexer.RPAREN) || b.check(lexer.COMMA) {
				components = append(components, nil)
			} else {
				components = append(components, b.parseExpression())
			}
		}
		b.expect(lexer.RPAREN)
		
		node := &ast.TupleExpression{
			BaseNode:   ast.BaseNode{Type: ast.NodeTupleExpression},
			Components: components,
			IsArray:    false,
		}
		b.setLocation(node, startTok, b.previous())
		return node
	}
	
	// Single expression in parentheses
	b.expect(lexer.RPAREN)
	return components[0]
}

func (b *Builder) parseArrayLiteral() *ast.TupleExpression {
	startTok := b.advance() // [
	
	node := &ast.TupleExpression{
		BaseNode:   ast.BaseNode{Type: ast.NodeTupleExpression},
		Components: make([]ast.Node, 0),
		IsArray:    true,
	}
	
	if !b.check(lexer.RBRACK) {
		node.Components = b.parseExpressionList()
	}
	
	b.expect(lexer.RBRACK)
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseNewExpression() *ast.NewExpression {
	startTok := b.advance() // new
	
	typeName := b.parseTypeName()
	
	node := &ast.NewExpression{
		BaseNode: ast.BaseNode{Type: ast.NodeNewExpression},
		TypeName: typeName,
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseTypeExpression() ast.Node {
	b.advance() // type
	b.expect(lexer.LPAREN)
	typeName := b.parseTypeName()
	b.expect(lexer.RPAREN)
	
	// Return the type name as a member access: type(T)
	return &ast.FunctionCall{
		BaseNode: ast.BaseNode{Type: ast.NodeFunctionCall},
		Expression: &ast.Identifier{
			BaseNode: ast.BaseNode{Type: ast.NodeIdentifier},
			Name:     "type",
		},
		Arguments: []ast.Node{typeName},
	}
}

func (b *Builder) parsePayableConversion() ast.Node {
	startTok := b.advance() // payable
	b.expect(lexer.LPAREN)
	expr := b.parseExpression()
	b.expect(lexer.RPAREN)
	
	node := &ast.FunctionCall{
		BaseNode: ast.BaseNode{Type: ast.NodeFunctionCall},
		Expression: &ast.ElementaryTypeName{
			BaseNode:        ast.BaseNode{Type: ast.NodeElementaryTypeName},
			Name:            "address",
			StateMutability: "payable",
		},
		Arguments: []ast.Node{expr},
	}
	
	b.setLocation(node, startTok, b.previous())
	return node
}

func (b *Builder) parseExpressionList() []ast.Node {
	var exprs []ast.Node
	
	if b.check(lexer.RPAREN) || b.check(lexer.RBRACK) || b.check(lexer.RBRACE) {
		return exprs
	}
	
	exprs = append(exprs, b.parseExpression())
	
	for b.check(lexer.COMMA) {
		b.advance() // ,
		if b.check(lexer.RPAREN) || b.check(lexer.RBRACK) || b.check(lexer.RBRACE) {
			break
		}
		exprs = append(exprs, b.parseExpression())
	}
	
	return exprs
}

func (b *Builder) isAssignmentOperator() bool {
	return b.check(lexer.ASSIGN) || b.check(lexer.ASSIGN_ADD) || b.check(lexer.ASSIGN_SUB) ||
		b.check(lexer.ASSIGN_MUL) || b.check(lexer.ASSIGN_DIV) || b.check(lexer.ASSIGN_MOD) ||
		b.check(lexer.ASSIGN_AND) || b.check(lexer.ASSIGN_OR) || b.check(lexer.ASSIGN_XOR) ||
		b.check(lexer.ASSIGN_SHL) || b.check(lexer.ASSIGN_SHR) || b.check(lexer.ASSIGN_SAR)
}

func (b *Builder) checkNumberUnit() bool {
	if !b.check(lexer.IDENTIFIER) {
		return false
	}
	unit := b.peek().Value
	units := []string{"wei", "gwei", "ether", "seconds", "minutes", "hours", "days", "weeks", "years"}
	for _, u := range units {
		if unit == u {
			return true
		}
	}
	return false
}

