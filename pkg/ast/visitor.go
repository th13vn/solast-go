package ast

// Visitor is the interface for visiting AST nodes
type Visitor interface {
	// Return true to continue visiting children, false to stop
	VisitSourceUnit(node *SourceUnit) bool
	VisitPragmaDirective(node *PragmaDirective) bool
	VisitImportDirective(node *ImportDirective) bool
	VisitContractDefinition(node *ContractDefinition) bool
	VisitInheritanceSpecifier(node *InheritanceSpecifier) bool
	VisitFunctionDefinition(node *FunctionDefinition) bool
	VisitModifierDefinition(node *ModifierDefinition) bool
	VisitModifierInvocation(node *ModifierInvocation) bool
	VisitStateVariableDeclaration(node *StateVariableDeclaration) bool
	VisitVariableDeclaration(node *VariableDeclaration) bool
	VisitVariableDeclarationStatement(node *VariableDeclarationStatement) bool
	VisitStructDefinition(node *StructDefinition) bool
	VisitEnumDefinition(node *EnumDefinition) bool
	VisitEnumValue(node *EnumValue) bool
	VisitEventDefinition(node *EventDefinition) bool
	VisitErrorDefinition(node *ErrorDefinition) bool
	VisitUserDefinedValueTypeDefinition(node *UserDefinedValueTypeDefinition) bool
	VisitUsingForDeclaration(node *UsingForDeclaration) bool
	VisitElementaryTypeName(node *ElementaryTypeName) bool
	VisitUserDefinedTypeName(node *UserDefinedTypeName) bool
	VisitMapping(node *Mapping) bool
	VisitArrayTypeName(node *ArrayTypeName) bool
	VisitFunctionTypeName(node *FunctionTypeName) bool
	VisitBlock(node *Block) bool
	VisitUncheckedBlock(node *UncheckedBlock) bool
	VisitExpressionStatement(node *ExpressionStatement) bool
	VisitIfStatement(node *IfStatement) bool
	VisitWhileStatement(node *WhileStatement) bool
	VisitDoWhileStatement(node *DoWhileStatement) bool
	VisitForStatement(node *ForStatement) bool
	VisitContinueStatement(node *ContinueStatement) bool
	VisitBreakStatement(node *BreakStatement) bool
	VisitReturnStatement(node *ReturnStatement) bool
	VisitEmitStatement(node *EmitStatement) bool
	VisitRevertStatement(node *RevertStatement) bool
	VisitTryStatement(node *TryStatement) bool
	VisitCatchClause(node *CatchClause) bool
	VisitBinaryOperation(node *BinaryOperation) bool
	VisitUnaryOperation(node *UnaryOperation) bool
	VisitConditional(node *Conditional) bool
	VisitFunctionCall(node *FunctionCall) bool
	VisitFunctionCallOptions(node *FunctionCallOptions) bool
	VisitMemberAccess(node *MemberAccess) bool
	VisitIndexAccess(node *IndexAccess) bool
	VisitIndexRangeAccess(node *IndexRangeAccess) bool
	VisitNewExpression(node *NewExpression) bool
	VisitTupleExpression(node *TupleExpression) bool
	VisitNameValueExpression(node *NameValueExpression) bool
	VisitIdentifier(node *Identifier) bool
	VisitNumberLiteral(node *NumberLiteral) bool
	VisitBooleanLiteral(node *BooleanLiteral) bool
	VisitStringLiteral(node *StringLiteral) bool
	VisitHexLiteral(node *HexLiteral) bool
	VisitInlineAssembly(node *InlineAssembly) bool
	VisitAssemblyBlock(node *AssemblyBlock) bool
	VisitAssemblyCall(node *AssemblyCall) bool
	VisitAssemblyLocalDefinition(node *AssemblyLocalDefinition) bool
	VisitAssemblyAssignment(node *AssemblyAssignment) bool
	VisitAssemblyIdentifier(node *AssemblyIdentifier) bool
	VisitAssemblyLiteral(node *AssemblyLiteral) bool
	VisitAssemblyIf(node *AssemblyIf) bool
	VisitAssemblySwitch(node *AssemblySwitch) bool
	VisitAssemblyCase(node *AssemblyCase) bool
	VisitAssemblyFor(node *AssemblyFor) bool
	VisitAssemblyFunctionDefinition(node *AssemblyFunctionDefinition) bool
}

// BaseVisitor provides default implementations for all visitor methods
type BaseVisitor struct{}

func (v *BaseVisitor) VisitSourceUnit(node *SourceUnit) bool                       { return true }
func (v *BaseVisitor) VisitPragmaDirective(node *PragmaDirective) bool             { return true }
func (v *BaseVisitor) VisitImportDirective(node *ImportDirective) bool             { return true }
func (v *BaseVisitor) VisitContractDefinition(node *ContractDefinition) bool       { return true }
func (v *BaseVisitor) VisitInheritanceSpecifier(node *InheritanceSpecifier) bool   { return true }
func (v *BaseVisitor) VisitFunctionDefinition(node *FunctionDefinition) bool       { return true }
func (v *BaseVisitor) VisitModifierDefinition(node *ModifierDefinition) bool       { return true }
func (v *BaseVisitor) VisitModifierInvocation(node *ModifierInvocation) bool       { return true }
func (v *BaseVisitor) VisitStateVariableDeclaration(node *StateVariableDeclaration) bool { return true }
func (v *BaseVisitor) VisitVariableDeclaration(node *VariableDeclaration) bool     { return true }
func (v *BaseVisitor) VisitVariableDeclarationStatement(node *VariableDeclarationStatement) bool { return true }
func (v *BaseVisitor) VisitStructDefinition(node *StructDefinition) bool           { return true }
func (v *BaseVisitor) VisitEnumDefinition(node *EnumDefinition) bool               { return true }
func (v *BaseVisitor) VisitEnumValue(node *EnumValue) bool                         { return true }
func (v *BaseVisitor) VisitEventDefinition(node *EventDefinition) bool             { return true }
func (v *BaseVisitor) VisitErrorDefinition(node *ErrorDefinition) bool             { return true }
func (v *BaseVisitor) VisitUserDefinedValueTypeDefinition(node *UserDefinedValueTypeDefinition) bool { return true }
func (v *BaseVisitor) VisitUsingForDeclaration(node *UsingForDeclaration) bool     { return true }
func (v *BaseVisitor) VisitElementaryTypeName(node *ElementaryTypeName) bool       { return true }
func (v *BaseVisitor) VisitUserDefinedTypeName(node *UserDefinedTypeName) bool     { return true }
func (v *BaseVisitor) VisitMapping(node *Mapping) bool                             { return true }
func (v *BaseVisitor) VisitArrayTypeName(node *ArrayTypeName) bool                 { return true }
func (v *BaseVisitor) VisitFunctionTypeName(node *FunctionTypeName) bool           { return true }
func (v *BaseVisitor) VisitBlock(node *Block) bool                                 { return true }
func (v *BaseVisitor) VisitUncheckedBlock(node *UncheckedBlock) bool               { return true }
func (v *BaseVisitor) VisitExpressionStatement(node *ExpressionStatement) bool     { return true }
func (v *BaseVisitor) VisitIfStatement(node *IfStatement) bool                     { return true }
func (v *BaseVisitor) VisitWhileStatement(node *WhileStatement) bool               { return true }
func (v *BaseVisitor) VisitDoWhileStatement(node *DoWhileStatement) bool           { return true }
func (v *BaseVisitor) VisitForStatement(node *ForStatement) bool                   { return true }
func (v *BaseVisitor) VisitContinueStatement(node *ContinueStatement) bool         { return true }
func (v *BaseVisitor) VisitBreakStatement(node *BreakStatement) bool               { return true }
func (v *BaseVisitor) VisitReturnStatement(node *ReturnStatement) bool             { return true }
func (v *BaseVisitor) VisitEmitStatement(node *EmitStatement) bool                 { return true }
func (v *BaseVisitor) VisitRevertStatement(node *RevertStatement) bool             { return true }
func (v *BaseVisitor) VisitTryStatement(node *TryStatement) bool                   { return true }
func (v *BaseVisitor) VisitCatchClause(node *CatchClause) bool                     { return true }
func (v *BaseVisitor) VisitBinaryOperation(node *BinaryOperation) bool             { return true }
func (v *BaseVisitor) VisitUnaryOperation(node *UnaryOperation) bool               { return true }
func (v *BaseVisitor) VisitConditional(node *Conditional) bool                     { return true }
func (v *BaseVisitor) VisitFunctionCall(node *FunctionCall) bool                   { return true }
func (v *BaseVisitor) VisitFunctionCallOptions(node *FunctionCallOptions) bool     { return true }
func (v *BaseVisitor) VisitMemberAccess(node *MemberAccess) bool                   { return true }
func (v *BaseVisitor) VisitIndexAccess(node *IndexAccess) bool                     { return true }
func (v *BaseVisitor) VisitIndexRangeAccess(node *IndexRangeAccess) bool           { return true }
func (v *BaseVisitor) VisitNewExpression(node *NewExpression) bool                 { return true }
func (v *BaseVisitor) VisitTupleExpression(node *TupleExpression) bool             { return true }
func (v *BaseVisitor) VisitNameValueExpression(node *NameValueExpression) bool     { return true }
func (v *BaseVisitor) VisitIdentifier(node *Identifier) bool                       { return true }
func (v *BaseVisitor) VisitNumberLiteral(node *NumberLiteral) bool                 { return true }
func (v *BaseVisitor) VisitBooleanLiteral(node *BooleanLiteral) bool               { return true }
func (v *BaseVisitor) VisitStringLiteral(node *StringLiteral) bool                 { return true }
func (v *BaseVisitor) VisitHexLiteral(node *HexLiteral) bool                       { return true }
func (v *BaseVisitor) VisitInlineAssembly(node *InlineAssembly) bool               { return true }
func (v *BaseVisitor) VisitAssemblyBlock(node *AssemblyBlock) bool                 { return true }
func (v *BaseVisitor) VisitAssemblyCall(node *AssemblyCall) bool                   { return true }
func (v *BaseVisitor) VisitAssemblyLocalDefinition(node *AssemblyLocalDefinition) bool { return true }
func (v *BaseVisitor) VisitAssemblyAssignment(node *AssemblyAssignment) bool       { return true }
func (v *BaseVisitor) VisitAssemblyIdentifier(node *AssemblyIdentifier) bool       { return true }
func (v *BaseVisitor) VisitAssemblyLiteral(node *AssemblyLiteral) bool             { return true }
func (v *BaseVisitor) VisitAssemblyIf(node *AssemblyIf) bool                       { return true }
func (v *BaseVisitor) VisitAssemblySwitch(node *AssemblySwitch) bool               { return true }
func (v *BaseVisitor) VisitAssemblyCase(node *AssemblyCase) bool                   { return true }
func (v *BaseVisitor) VisitAssemblyFor(node *AssemblyFor) bool                     { return true }
func (v *BaseVisitor) VisitAssemblyFunctionDefinition(node *AssemblyFunctionDefinition) bool { return true }

// SimpleVisitor allows specifying only the callbacks you care about
type SimpleVisitor struct {
	BaseVisitor
	SourceUnitFn                       func(*SourceUnit)
	PragmaDirectiveFn                  func(*PragmaDirective)
	ImportDirectiveFn                  func(*ImportDirective)
	ContractDefinitionFn               func(*ContractDefinition)
	InheritanceSpecifierFn             func(*InheritanceSpecifier)
	FunctionDefinitionFn               func(*FunctionDefinition)
	ModifierDefinitionFn               func(*ModifierDefinition)
	ModifierInvocationFn               func(*ModifierInvocation)
	StateVariableDeclarationFn         func(*StateVariableDeclaration)
	VariableDeclarationFn              func(*VariableDeclaration)
	VariableDeclarationStatementFn     func(*VariableDeclarationStatement)
	StructDefinitionFn                 func(*StructDefinition)
	EnumDefinitionFn                   func(*EnumDefinition)
	EnumValueFn                        func(*EnumValue)
	EventDefinitionFn                  func(*EventDefinition)
	ErrorDefinitionFn                  func(*ErrorDefinition)
	UserDefinedValueTypeDefinitionFn   func(*UserDefinedValueTypeDefinition)
	UsingForDeclarationFn              func(*UsingForDeclaration)
	ElementaryTypeNameFn               func(*ElementaryTypeName)
	UserDefinedTypeNameFn              func(*UserDefinedTypeName)
	MappingFn                          func(*Mapping)
	ArrayTypeNameFn                    func(*ArrayTypeName)
	FunctionTypeNameFn                 func(*FunctionTypeName)
	BlockFn                            func(*Block)
	UncheckedBlockFn                   func(*UncheckedBlock)
	ExpressionStatementFn              func(*ExpressionStatement)
	IfStatementFn                      func(*IfStatement)
	WhileStatementFn                   func(*WhileStatement)
	DoWhileStatementFn                 func(*DoWhileStatement)
	ForStatementFn                     func(*ForStatement)
	ContinueStatementFn                func(*ContinueStatement)
	BreakStatementFn                   func(*BreakStatement)
	ReturnStatementFn                  func(*ReturnStatement)
	EmitStatementFn                    func(*EmitStatement)
	RevertStatementFn                  func(*RevertStatement)
	TryStatementFn                     func(*TryStatement)
	CatchClauseFn                      func(*CatchClause)
	BinaryOperationFn                  func(*BinaryOperation)
	UnaryOperationFn                   func(*UnaryOperation)
	ConditionalFn                      func(*Conditional)
	FunctionCallFn                     func(*FunctionCall)
	FunctionCallOptionsFn              func(*FunctionCallOptions)
	MemberAccessFn                     func(*MemberAccess)
	IndexAccessFn                      func(*IndexAccess)
	IndexRangeAccessFn                 func(*IndexRangeAccess)
	NewExpressionFn                    func(*NewExpression)
	TupleExpressionFn                  func(*TupleExpression)
	NameValueExpressionFn              func(*NameValueExpression)
	IdentifierFn                       func(*Identifier)
	NumberLiteralFn                    func(*NumberLiteral)
	BooleanLiteralFn                   func(*BooleanLiteral)
	StringLiteralFn                    func(*StringLiteral)
	HexLiteralFn                       func(*HexLiteral)
	InlineAssemblyFn                   func(*InlineAssembly)
	AssemblyBlockFn                    func(*AssemblyBlock)
	AssemblyCallFn                     func(*AssemblyCall)
	AssemblyLocalDefinitionFn          func(*AssemblyLocalDefinition)
	AssemblyAssignmentFn               func(*AssemblyAssignment)
	AssemblyIdentifierFn               func(*AssemblyIdentifier)
	AssemblyLiteralFn                  func(*AssemblyLiteral)
	AssemblyIfFn                       func(*AssemblyIf)
	AssemblySwitchFn                   func(*AssemblySwitch)
	AssemblyCaseFn                     func(*AssemblyCase)
	AssemblyForFn                      func(*AssemblyFor)
	AssemblyFunctionDefinitionFn       func(*AssemblyFunctionDefinition)
}

// Walk traverses the AST and calls the appropriate visitor method for each node
func Walk(node Node, visitor Visitor) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *SourceUnit:
		if visitor.VisitSourceUnit(n) {
			for _, child := range n.Children {
				Walk(child, visitor)
			}
		}
	case *PragmaDirective:
		visitor.VisitPragmaDirective(n)
	case *ImportDirective:
		visitor.VisitImportDirective(n)
	case *ContractDefinition:
		if visitor.VisitContractDefinition(n) {
			for _, base := range n.BaseContracts {
				Walk(base, visitor)
			}
			for _, sub := range n.SubNodes {
				Walk(sub, visitor)
			}
		}
	case *InheritanceSpecifier:
		if visitor.VisitInheritanceSpecifier(n) {
			Walk(n.BaseName, visitor)
			for _, arg := range n.Arguments {
				Walk(arg, visitor)
			}
		}
	case *FunctionDefinition:
		if visitor.VisitFunctionDefinition(n) {
			for _, param := range n.Parameters {
				Walk(param, visitor)
			}
			for _, param := range n.ReturnParameters {
				Walk(param, visitor)
			}
			for _, mod := range n.Modifiers {
				Walk(mod, visitor)
			}
			Walk(n.Body, visitor)
		}
	case *ModifierDefinition:
		if visitor.VisitModifierDefinition(n) {
			for _, param := range n.Parameters {
				Walk(param, visitor)
			}
			Walk(n.Body, visitor)
		}
	case *ModifierInvocation:
		if visitor.VisitModifierInvocation(n) {
			for _, arg := range n.Arguments {
				Walk(arg, visitor)
			}
		}
	case *StateVariableDeclaration:
		if visitor.VisitStateVariableDeclaration(n) {
			for _, v := range n.Variables {
				Walk(v, visitor)
			}
			Walk(n.InitialValue, visitor)
		}
	case *VariableDeclaration:
		if visitor.VisitVariableDeclaration(n) {
			Walk(n.TypeName, visitor)
			Walk(n.Expression, visitor)
		}
	case *VariableDeclarationStatement:
		if visitor.VisitVariableDeclarationStatement(n) {
			for _, v := range n.Variables {
				Walk(v, visitor)
			}
			Walk(n.InitialValue, visitor)
		}
	case *StructDefinition:
		if visitor.VisitStructDefinition(n) {
			for _, member := range n.Members {
				Walk(member, visitor)
			}
		}
	case *EnumDefinition:
		if visitor.VisitEnumDefinition(n) {
			for _, member := range n.Members {
				Walk(member, visitor)
			}
		}
	case *EnumValue:
		visitor.VisitEnumValue(n)
	case *EventDefinition:
		if visitor.VisitEventDefinition(n) {
			for _, param := range n.Parameters {
				Walk(param, visitor)
			}
		}
	case *ErrorDefinition:
		if visitor.VisitErrorDefinition(n) {
			for _, param := range n.Parameters {
				Walk(param, visitor)
			}
		}
	case *UserDefinedValueTypeDefinition:
		if visitor.VisitUserDefinedValueTypeDefinition(n) {
			Walk(n.UnderlyingType, visitor)
		}
	case *UsingForDeclaration:
		if visitor.VisitUsingForDeclaration(n) {
			Walk(n.TypeName, visitor)
		}
	case *ElementaryTypeName:
		visitor.VisitElementaryTypeName(n)
	case *UserDefinedTypeName:
		visitor.VisitUserDefinedTypeName(n)
	case *Mapping:
		if visitor.VisitMapping(n) {
			Walk(n.KeyType, visitor)
			Walk(n.ValueType, visitor)
		}
	case *ArrayTypeName:
		if visitor.VisitArrayTypeName(n) {
			Walk(n.BaseTypeName, visitor)
			Walk(n.Length, visitor)
		}
	case *FunctionTypeName:
		if visitor.VisitFunctionTypeName(n) {
			for _, param := range n.ParameterTypes {
				Walk(param, visitor)
			}
			for _, ret := range n.ReturnTypes {
				Walk(ret, visitor)
			}
		}
	case *Block:
		if visitor.VisitBlock(n) {
			for _, stmt := range n.Statements {
				Walk(stmt, visitor)
			}
		}
	case *UncheckedBlock:
		if visitor.VisitUncheckedBlock(n) {
			Walk(n.Body, visitor)
		}
	case *ExpressionStatement:
		if visitor.VisitExpressionStatement(n) {
			Walk(n.Expression, visitor)
		}
	case *IfStatement:
		if visitor.VisitIfStatement(n) {
			Walk(n.Condition, visitor)
			Walk(n.TrueBody, visitor)
			Walk(n.FalseBody, visitor)
		}
	case *WhileStatement:
		if visitor.VisitWhileStatement(n) {
			Walk(n.Condition, visitor)
			Walk(n.Body, visitor)
		}
	case *DoWhileStatement:
		if visitor.VisitDoWhileStatement(n) {
			Walk(n.Condition, visitor)
			Walk(n.Body, visitor)
		}
	case *ForStatement:
		if visitor.VisitForStatement(n) {
			Walk(n.InitExpression, visitor)
			Walk(n.ConditionExpression, visitor)
			Walk(n.LoopExpression, visitor)
			Walk(n.Body, visitor)
		}
	case *ContinueStatement:
		visitor.VisitContinueStatement(n)
	case *BreakStatement:
		visitor.VisitBreakStatement(n)
	case *ReturnStatement:
		if visitor.VisitReturnStatement(n) {
			Walk(n.Expression, visitor)
		}
	case *EmitStatement:
		if visitor.VisitEmitStatement(n) {
			Walk(n.EventCall, visitor)
		}
	case *RevertStatement:
		if visitor.VisitRevertStatement(n) {
			Walk(n.RevertCall, visitor)
		}
	case *TryStatement:
		if visitor.VisitTryStatement(n) {
			Walk(n.Expression, visitor)
			Walk(n.Body, visitor)
			for _, clause := range n.CatchClauses {
				Walk(clause, visitor)
			}
		}
	case *CatchClause:
		if visitor.VisitCatchClause(n) {
			for _, param := range n.Parameters {
				Walk(param, visitor)
			}
			Walk(n.Body, visitor)
		}
	case *BinaryOperation:
		if visitor.VisitBinaryOperation(n) {
			Walk(n.Left, visitor)
			Walk(n.Right, visitor)
		}
	case *UnaryOperation:
		if visitor.VisitUnaryOperation(n) {
			Walk(n.SubExpression, visitor)
		}
	case *Conditional:
		if visitor.VisitConditional(n) {
			Walk(n.Condition, visitor)
			Walk(n.TrueExpression, visitor)
			Walk(n.FalseExpression, visitor)
		}
	case *FunctionCall:
		if visitor.VisitFunctionCall(n) {
			Walk(n.Expression, visitor)
			for _, arg := range n.Arguments {
				Walk(arg, visitor)
			}
		}
	case *FunctionCallOptions:
		if visitor.VisitFunctionCallOptions(n) {
			Walk(n.Expression, visitor)
			for _, opt := range n.Options {
				Walk(opt, visitor)
			}
		}
	case *MemberAccess:
		if visitor.VisitMemberAccess(n) {
			Walk(n.Expression, visitor)
		}
	case *IndexAccess:
		if visitor.VisitIndexAccess(n) {
			Walk(n.Base, visitor)
			Walk(n.Index, visitor)
		}
	case *IndexRangeAccess:
		if visitor.VisitIndexRangeAccess(n) {
			Walk(n.Base, visitor)
			Walk(n.IndexStart, visitor)
			Walk(n.IndexEnd, visitor)
		}
	case *NewExpression:
		if visitor.VisitNewExpression(n) {
			Walk(n.TypeName, visitor)
		}
	case *TupleExpression:
		if visitor.VisitTupleExpression(n) {
			for _, comp := range n.Components {
				Walk(comp, visitor)
			}
		}
	case *NameValueExpression:
		if visitor.VisitNameValueExpression(n) {
			Walk(n.Expression, visitor)
		}
	case *Identifier:
		visitor.VisitIdentifier(n)
	case *NumberLiteral:
		visitor.VisitNumberLiteral(n)
	case *BooleanLiteral:
		visitor.VisitBooleanLiteral(n)
	case *StringLiteral:
		visitor.VisitStringLiteral(n)
	case *HexLiteral:
		visitor.VisitHexLiteral(n)
	case *InlineAssembly:
		if visitor.VisitInlineAssembly(n) {
			Walk(n.Body, visitor)
		}
	case *AssemblyBlock:
		if visitor.VisitAssemblyBlock(n) {
			for _, op := range n.Operations {
				Walk(op, visitor)
			}
		}
	case *AssemblyCall:
		if visitor.VisitAssemblyCall(n) {
			for _, arg := range n.Arguments {
				Walk(arg, visitor)
			}
		}
	case *AssemblyLocalDefinition:
		if visitor.VisitAssemblyLocalDefinition(n) {
			Walk(n.Expression, visitor)
		}
	case *AssemblyAssignment:
		if visitor.VisitAssemblyAssignment(n) {
			Walk(n.Expression, visitor)
		}
	case *AssemblyIdentifier:
		visitor.VisitAssemblyIdentifier(n)
	case *AssemblyLiteral:
		visitor.VisitAssemblyLiteral(n)
	case *AssemblyIf:
		if visitor.VisitAssemblyIf(n) {
			Walk(n.Condition, visitor)
			Walk(n.Body, visitor)
		}
	case *AssemblySwitch:
		if visitor.VisitAssemblySwitch(n) {
			Walk(n.Expression, visitor)
			for _, c := range n.Cases {
				Walk(c, visitor)
			}
		}
	case *AssemblyCase:
		if visitor.VisitAssemblyCase(n) {
			Walk(n.Value, visitor)
			Walk(n.Body, visitor)
		}
	case *AssemblyFor:
		if visitor.VisitAssemblyFor(n) {
			Walk(n.Pre, visitor)
			Walk(n.Condition, visitor)
			Walk(n.Post, visitor)
			Walk(n.Body, visitor)
		}
	case *AssemblyFunctionDefinition:
		if visitor.VisitAssemblyFunctionDefinition(n) {
			Walk(n.Body, visitor)
		}
	}
}

// WalkSimple traverses the AST using a SimpleVisitor
func WalkSimple(node Node, visitor *SimpleVisitor) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *SourceUnit:
		if visitor.SourceUnitFn != nil {
			visitor.SourceUnitFn(n)
		}
		for _, child := range n.Children {
			WalkSimple(child, visitor)
		}
	case *PragmaDirective:
		if visitor.PragmaDirectiveFn != nil {
			visitor.PragmaDirectiveFn(n)
		}
	case *ImportDirective:
		if visitor.ImportDirectiveFn != nil {
			visitor.ImportDirectiveFn(n)
		}
	case *ContractDefinition:
		if visitor.ContractDefinitionFn != nil {
			visitor.ContractDefinitionFn(n)
		}
		for _, base := range n.BaseContracts {
			WalkSimple(base, visitor)
		}
		for _, sub := range n.SubNodes {
			WalkSimple(sub, visitor)
		}
	case *InheritanceSpecifier:
		if visitor.InheritanceSpecifierFn != nil {
			visitor.InheritanceSpecifierFn(n)
		}
		WalkSimple(n.BaseName, visitor)
		for _, arg := range n.Arguments {
			WalkSimple(arg, visitor)
		}
	case *FunctionDefinition:
		if visitor.FunctionDefinitionFn != nil {
			visitor.FunctionDefinitionFn(n)
		}
		for _, param := range n.Parameters {
			WalkSimple(param, visitor)
		}
		for _, param := range n.ReturnParameters {
			WalkSimple(param, visitor)
		}
		for _, mod := range n.Modifiers {
			WalkSimple(mod, visitor)
		}
		WalkSimple(n.Body, visitor)
	case *ModifierDefinition:
		if visitor.ModifierDefinitionFn != nil {
			visitor.ModifierDefinitionFn(n)
		}
		for _, param := range n.Parameters {
			WalkSimple(param, visitor)
		}
		WalkSimple(n.Body, visitor)
	case *ModifierInvocation:
		if visitor.ModifierInvocationFn != nil {
			visitor.ModifierInvocationFn(n)
		}
		for _, arg := range n.Arguments {
			WalkSimple(arg, visitor)
		}
	case *StateVariableDeclaration:
		if visitor.StateVariableDeclarationFn != nil {
			visitor.StateVariableDeclarationFn(n)
		}
		for _, v := range n.Variables {
			WalkSimple(v, visitor)
		}
		WalkSimple(n.InitialValue, visitor)
	case *VariableDeclaration:
		if visitor.VariableDeclarationFn != nil {
			visitor.VariableDeclarationFn(n)
		}
		WalkSimple(n.TypeName, visitor)
		WalkSimple(n.Expression, visitor)
	case *VariableDeclarationStatement:
		if visitor.VariableDeclarationStatementFn != nil {
			visitor.VariableDeclarationStatementFn(n)
		}
		for _, v := range n.Variables {
			WalkSimple(v, visitor)
		}
		WalkSimple(n.InitialValue, visitor)
	case *StructDefinition:
		if visitor.StructDefinitionFn != nil {
			visitor.StructDefinitionFn(n)
		}
		for _, member := range n.Members {
			WalkSimple(member, visitor)
		}
	case *EnumDefinition:
		if visitor.EnumDefinitionFn != nil {
			visitor.EnumDefinitionFn(n)
		}
		for _, member := range n.Members {
			WalkSimple(member, visitor)
		}
	case *EnumValue:
		if visitor.EnumValueFn != nil {
			visitor.EnumValueFn(n)
		}
	case *EventDefinition:
		if visitor.EventDefinitionFn != nil {
			visitor.EventDefinitionFn(n)
		}
		for _, param := range n.Parameters {
			WalkSimple(param, visitor)
		}
	case *ErrorDefinition:
		if visitor.ErrorDefinitionFn != nil {
			visitor.ErrorDefinitionFn(n)
		}
		for _, param := range n.Parameters {
			WalkSimple(param, visitor)
		}
	case *UserDefinedValueTypeDefinition:
		if visitor.UserDefinedValueTypeDefinitionFn != nil {
			visitor.UserDefinedValueTypeDefinitionFn(n)
		}
		WalkSimple(n.UnderlyingType, visitor)
	case *UsingForDeclaration:
		if visitor.UsingForDeclarationFn != nil {
			visitor.UsingForDeclarationFn(n)
		}
		WalkSimple(n.TypeName, visitor)
	case *ElementaryTypeName:
		if visitor.ElementaryTypeNameFn != nil {
			visitor.ElementaryTypeNameFn(n)
		}
	case *UserDefinedTypeName:
		if visitor.UserDefinedTypeNameFn != nil {
			visitor.UserDefinedTypeNameFn(n)
		}
	case *Mapping:
		if visitor.MappingFn != nil {
			visitor.MappingFn(n)
		}
		WalkSimple(n.KeyType, visitor)
		WalkSimple(n.ValueType, visitor)
	case *ArrayTypeName:
		if visitor.ArrayTypeNameFn != nil {
			visitor.ArrayTypeNameFn(n)
		}
		WalkSimple(n.BaseTypeName, visitor)
		WalkSimple(n.Length, visitor)
	case *FunctionTypeName:
		if visitor.FunctionTypeNameFn != nil {
			visitor.FunctionTypeNameFn(n)
		}
		for _, param := range n.ParameterTypes {
			WalkSimple(param, visitor)
		}
		for _, ret := range n.ReturnTypes {
			WalkSimple(ret, visitor)
		}
	case *Block:
		if visitor.BlockFn != nil {
			visitor.BlockFn(n)
		}
		for _, stmt := range n.Statements {
			WalkSimple(stmt, visitor)
		}
	case *UncheckedBlock:
		if visitor.UncheckedBlockFn != nil {
			visitor.UncheckedBlockFn(n)
		}
		WalkSimple(n.Body, visitor)
	case *ExpressionStatement:
		if visitor.ExpressionStatementFn != nil {
			visitor.ExpressionStatementFn(n)
		}
		WalkSimple(n.Expression, visitor)
	case *IfStatement:
		if visitor.IfStatementFn != nil {
			visitor.IfStatementFn(n)
		}
		WalkSimple(n.Condition, visitor)
		WalkSimple(n.TrueBody, visitor)
		WalkSimple(n.FalseBody, visitor)
	case *WhileStatement:
		if visitor.WhileStatementFn != nil {
			visitor.WhileStatementFn(n)
		}
		WalkSimple(n.Condition, visitor)
		WalkSimple(n.Body, visitor)
	case *DoWhileStatement:
		if visitor.DoWhileStatementFn != nil {
			visitor.DoWhileStatementFn(n)
		}
		WalkSimple(n.Condition, visitor)
		WalkSimple(n.Body, visitor)
	case *ForStatement:
		if visitor.ForStatementFn != nil {
			visitor.ForStatementFn(n)
		}
		WalkSimple(n.InitExpression, visitor)
		WalkSimple(n.ConditionExpression, visitor)
		WalkSimple(n.LoopExpression, visitor)
		WalkSimple(n.Body, visitor)
	case *ContinueStatement:
		if visitor.ContinueStatementFn != nil {
			visitor.ContinueStatementFn(n)
		}
	case *BreakStatement:
		if visitor.BreakStatementFn != nil {
			visitor.BreakStatementFn(n)
		}
	case *ReturnStatement:
		if visitor.ReturnStatementFn != nil {
			visitor.ReturnStatementFn(n)
		}
		WalkSimple(n.Expression, visitor)
	case *EmitStatement:
		if visitor.EmitStatementFn != nil {
			visitor.EmitStatementFn(n)
		}
		WalkSimple(n.EventCall, visitor)
	case *RevertStatement:
		if visitor.RevertStatementFn != nil {
			visitor.RevertStatementFn(n)
		}
		WalkSimple(n.RevertCall, visitor)
	case *TryStatement:
		if visitor.TryStatementFn != nil {
			visitor.TryStatementFn(n)
		}
		WalkSimple(n.Expression, visitor)
		WalkSimple(n.Body, visitor)
		for _, clause := range n.CatchClauses {
			WalkSimple(clause, visitor)
		}
	case *CatchClause:
		if visitor.CatchClauseFn != nil {
			visitor.CatchClauseFn(n)
		}
		for _, param := range n.Parameters {
			WalkSimple(param, visitor)
		}
		WalkSimple(n.Body, visitor)
	case *BinaryOperation:
		if visitor.BinaryOperationFn != nil {
			visitor.BinaryOperationFn(n)
		}
		WalkSimple(n.Left, visitor)
		WalkSimple(n.Right, visitor)
	case *UnaryOperation:
		if visitor.UnaryOperationFn != nil {
			visitor.UnaryOperationFn(n)
		}
		WalkSimple(n.SubExpression, visitor)
	case *Conditional:
		if visitor.ConditionalFn != nil {
			visitor.ConditionalFn(n)
		}
		WalkSimple(n.Condition, visitor)
		WalkSimple(n.TrueExpression, visitor)
		WalkSimple(n.FalseExpression, visitor)
	case *FunctionCall:
		if visitor.FunctionCallFn != nil {
			visitor.FunctionCallFn(n)
		}
		WalkSimple(n.Expression, visitor)
		for _, arg := range n.Arguments {
			WalkSimple(arg, visitor)
		}
	case *FunctionCallOptions:
		if visitor.FunctionCallOptionsFn != nil {
			visitor.FunctionCallOptionsFn(n)
		}
		WalkSimple(n.Expression, visitor)
		for _, opt := range n.Options {
			WalkSimple(opt, visitor)
		}
	case *MemberAccess:
		if visitor.MemberAccessFn != nil {
			visitor.MemberAccessFn(n)
		}
		WalkSimple(n.Expression, visitor)
	case *IndexAccess:
		if visitor.IndexAccessFn != nil {
			visitor.IndexAccessFn(n)
		}
		WalkSimple(n.Base, visitor)
		WalkSimple(n.Index, visitor)
	case *IndexRangeAccess:
		if visitor.IndexRangeAccessFn != nil {
			visitor.IndexRangeAccessFn(n)
		}
		WalkSimple(n.Base, visitor)
		WalkSimple(n.IndexStart, visitor)
		WalkSimple(n.IndexEnd, visitor)
	case *NewExpression:
		if visitor.NewExpressionFn != nil {
			visitor.NewExpressionFn(n)
		}
		WalkSimple(n.TypeName, visitor)
	case *TupleExpression:
		if visitor.TupleExpressionFn != nil {
			visitor.TupleExpressionFn(n)
		}
		for _, comp := range n.Components {
			WalkSimple(comp, visitor)
		}
	case *NameValueExpression:
		if visitor.NameValueExpressionFn != nil {
			visitor.NameValueExpressionFn(n)
		}
		WalkSimple(n.Expression, visitor)
	case *Identifier:
		if visitor.IdentifierFn != nil {
			visitor.IdentifierFn(n)
		}
	case *NumberLiteral:
		if visitor.NumberLiteralFn != nil {
			visitor.NumberLiteralFn(n)
		}
	case *BooleanLiteral:
		if visitor.BooleanLiteralFn != nil {
			visitor.BooleanLiteralFn(n)
		}
	case *StringLiteral:
		if visitor.StringLiteralFn != nil {
			visitor.StringLiteralFn(n)
		}
	case *HexLiteral:
		if visitor.HexLiteralFn != nil {
			visitor.HexLiteralFn(n)
		}
	case *InlineAssembly:
		if visitor.InlineAssemblyFn != nil {
			visitor.InlineAssemblyFn(n)
		}
		WalkSimple(n.Body, visitor)
	case *AssemblyBlock:
		if visitor.AssemblyBlockFn != nil {
			visitor.AssemblyBlockFn(n)
		}
		for _, op := range n.Operations {
			WalkSimple(op, visitor)
		}
	case *AssemblyCall:
		if visitor.AssemblyCallFn != nil {
			visitor.AssemblyCallFn(n)
		}
		for _, arg := range n.Arguments {
			WalkSimple(arg, visitor)
		}
	case *AssemblyLocalDefinition:
		if visitor.AssemblyLocalDefinitionFn != nil {
			visitor.AssemblyLocalDefinitionFn(n)
		}
		WalkSimple(n.Expression, visitor)
	case *AssemblyAssignment:
		if visitor.AssemblyAssignmentFn != nil {
			visitor.AssemblyAssignmentFn(n)
		}
		WalkSimple(n.Expression, visitor)
	case *AssemblyIdentifier:
		if visitor.AssemblyIdentifierFn != nil {
			visitor.AssemblyIdentifierFn(n)
		}
	case *AssemblyLiteral:
		if visitor.AssemblyLiteralFn != nil {
			visitor.AssemblyLiteralFn(n)
		}
	case *AssemblyIf:
		if visitor.AssemblyIfFn != nil {
			visitor.AssemblyIfFn(n)
		}
		WalkSimple(n.Condition, visitor)
		WalkSimple(n.Body, visitor)
	case *AssemblySwitch:
		if visitor.AssemblySwitchFn != nil {
			visitor.AssemblySwitchFn(n)
		}
		WalkSimple(n.Expression, visitor)
		for _, c := range n.Cases {
			WalkSimple(c, visitor)
		}
	case *AssemblyCase:
		if visitor.AssemblyCaseFn != nil {
			visitor.AssemblyCaseFn(n)
		}
		WalkSimple(n.Value, visitor)
		WalkSimple(n.Body, visitor)
	case *AssemblyFor:
		if visitor.AssemblyForFn != nil {
			visitor.AssemblyForFn(n)
		}
		WalkSimple(n.Pre, visitor)
		WalkSimple(n.Condition, visitor)
		WalkSimple(n.Post, visitor)
		WalkSimple(n.Body, visitor)
	case *AssemblyFunctionDefinition:
		if visitor.AssemblyFunctionDefinitionFn != nil {
			visitor.AssemblyFunctionDefinitionFn(n)
		}
		WalkSimple(n.Body, visitor)
	}
}

