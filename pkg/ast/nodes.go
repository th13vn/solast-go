// Package ast defines the AST node types for Solidity code.
// The node types are designed to be exactly compatible with the TypeScript
// solidity-parser (https://github.com/solidity-parser/parser).
package ast

import "encoding/json"

// NodeType represents the type of an AST node
type NodeType string

// Node types matching TypeScript solidity-parser exactly
const (
	// Top-level
	NodeSourceUnit      NodeType = "SourceUnit"
	NodePragmaDirective NodeType = "PragmaDirective"
	NodeImportDirective NodeType = "ImportDirective"

	// Definitions
	NodeContractDefinition          NodeType = "ContractDefinition"
	NodeInterfaceDefinition         NodeType = "InterfaceDefinition"
	NodeLibraryDefinition           NodeType = "LibraryDefinition"
	NodeFunctionDefinition          NodeType = "FunctionDefinition"
	NodeModifierDefinition          NodeType = "ModifierDefinition"
	NodeStructDefinition            NodeType = "StructDefinition"
	NodeEnumDefinition              NodeType = "EnumDefinition"
	NodeEnumValue                   NodeType = "EnumValue"
	NodeEventDefinition             NodeType = "EventDefinition"
	NodeErrorDefinition             NodeType = "ErrorDefinition"
	NodeUserDefinedValueTypeDefinition NodeType = "UserDefinedValueTypeDefinition"
	NodeUsingForDeclaration         NodeType = "UsingForDeclaration"

	// Variables
	NodeStateVariableDeclaration    NodeType = "StateVariableDeclaration"
	NodeVariableDeclaration         NodeType = "VariableDeclaration"
	NodeVariableDeclarationStatement NodeType = "VariableDeclarationStatement"

	// Types
	NodeElementaryTypeName  NodeType = "ElementaryTypeName"
	NodeUserDefinedTypeName NodeType = "UserDefinedTypeName"
	NodeMapping             NodeType = "Mapping"
	NodeArrayTypeName       NodeType = "ArrayTypeName"
	NodeFunctionTypeName    NodeType = "FunctionTypeName"

	// Statements
	NodeBlock               NodeType = "Block"
	NodeUncheckedBlock      NodeType = "UncheckedBlock"
	NodeExpressionStatement NodeType = "ExpressionStatement"
	NodeIfStatement         NodeType = "IfStatement"
	NodeWhileStatement      NodeType = "WhileStatement"
	NodeDoWhileStatement    NodeType = "DoWhileStatement"
	NodeForStatement        NodeType = "ForStatement"
	NodeContinueStatement   NodeType = "ContinueStatement"
	NodeBreakStatement      NodeType = "BreakStatement"
	NodeReturnStatement     NodeType = "ReturnStatement"
	NodeEmitStatement       NodeType = "EmitStatement"
	NodeRevertStatement     NodeType = "RevertStatement"
	NodeTryStatement        NodeType = "TryStatement"
	NodeCatchClause         NodeType = "CatchClause"

	// Expressions
	NodeBinaryOperation        NodeType = "BinaryOperation"
	NodeUnaryOperation         NodeType = "UnaryOperation"
	NodeConditional            NodeType = "Conditional"
	NodeFunctionCall           NodeType = "FunctionCall"
	NodeFunctionCallOptions    NodeType = "FunctionCallOptions"
	NodeMemberAccess           NodeType = "MemberAccess"
	NodeIndexAccess            NodeType = "IndexAccess"
	NodeIndexRangeAccess       NodeType = "IndexRangeAccess"
	NodeNewExpression          NodeType = "NewExpression"
	NodeTupleExpression        NodeType = "TupleExpression"
	NodeNameValueExpression    NodeType = "NameValueExpression"
	NodeIdentifier             NodeType = "Identifier"
	NodeNumberLiteral          NodeType = "NumberLiteral"
	NodeBooleanLiteral         NodeType = "BooleanLiteral"
	NodeStringLiteral          NodeType = "StringLiteral"
	NodeHexLiteral             NodeType = "HexLiteral"

	// Assembly
	NodeInlineAssembly        NodeType = "InlineAssembly"
	NodeAssemblyBlock         NodeType = "AssemblyBlock"
	NodeAssemblyCall          NodeType = "AssemblyCall"
	NodeAssemblyLocalDefinition NodeType = "AssemblyLocalDefinition"
	NodeAssemblyAssignment    NodeType = "AssemblyAssignment"
	NodeAssemblyIdentifier    NodeType = "AssemblyIdentifier"
	NodeAssemblyLiteral       NodeType = "AssemblyLiteral"
	NodeAssemblyIf            NodeType = "AssemblyIf"
	NodeAssemblySwitch        NodeType = "AssemblySwitch"
	NodeAssemblyCase          NodeType = "AssemblyCase"
	NodeAssemblyFor           NodeType = "AssemblyFor"
	NodeAssemblyFunctionDefinition NodeType = "AssemblyFunctionDefinition"
	NodeAssemblyFunctionReturns    NodeType = "AssemblyFunctionReturns"

	// Misc
	NodeModifierInvocation    NodeType = "ModifierInvocation"
	NodeInheritanceSpecifier  NodeType = "InheritanceSpecifier"
	NodeEventParameter        NodeType = "EventParameter"
	NodeParameterList         NodeType = "ParameterList"
	NodeParameter             NodeType = "Parameter"
)

// Location represents the source location of a node
type Location struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position represents a position in the source code
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Range represents a character range in the source code
type Range [2]int

// Node is the interface that all AST nodes implement
type Node interface {
	GetType() NodeType
	GetLocation() *Location
	GetRange() *Range
}

// BaseNode contains common fields for all AST nodes
type BaseNode struct {
	Type  NodeType  `json:"type"`
	Loc   *Location `json:"loc,omitempty"`
	Range *Range    `json:"range,omitempty"`
}

func (n *BaseNode) GetType() NodeType     { return n.Type }
func (n *BaseNode) GetLocation() *Location { return n.Loc }
func (n *BaseNode) GetRange() *Range       { return n.Range }

// SourceUnit is the root node of the AST
type SourceUnit struct {
	BaseNode
	Children []Node `json:"children"`
}

// PragmaDirective represents a pragma statement
type PragmaDirective struct {
	BaseNode
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ImportDirective represents an import statement
type ImportDirective struct {
	BaseNode
	Path         string              `json:"path"`
	PathLiteral  *StringLiteral      `json:"pathLiteral,omitempty"`
	UnitAlias    string              `json:"unitAlias,omitempty"`
	UnitAliasIdentifier *Identifier  `json:"unitAliasIdentifier,omitempty"`
	SymbolAliases []*ImportSymbol    `json:"symbolAliases,omitempty"`
	SymbolAliasesIdentifiers []*ImportSymbolIdentifiers `json:"symbolAliasesIdentifiers,omitempty"`
}

// ImportSymbol represents an imported symbol
type ImportSymbol struct {
	Symbol string `json:"symbol"`
	Alias  string `json:"alias,omitempty"`
}

// ImportSymbolIdentifiers represents imported symbol identifiers
type ImportSymbolIdentifiers struct {
	Symbol *Identifier `json:"symbol"`
	Alias  *Identifier `json:"alias,omitempty"`
}

// ContractDefinition represents a contract, interface, or library
type ContractDefinition struct {
	BaseNode
	Name          string                   `json:"name"`
	BaseContracts []*InheritanceSpecifier  `json:"baseContracts"`
	SubNodes      []Node                   `json:"subNodes"`
	Kind          string                   `json:"kind"` // "contract", "interface", "library", "abstract"
}

// InheritanceSpecifier represents a base contract
type InheritanceSpecifier struct {
	BaseNode
	BaseName  *UserDefinedTypeName `json:"baseName"`
	Arguments []Node               `json:"arguments,omitempty"`
}

// FunctionDefinition represents a function, constructor, fallback, or receive
type FunctionDefinition struct {
	BaseNode
	Name           string              `json:"name,omitempty"`
	Parameters     []*VariableDeclaration `json:"parameters"`
	ReturnParameters []*VariableDeclaration `json:"returnParameters,omitempty"`
	Body           *Block              `json:"body,omitempty"`
	Visibility     string              `json:"visibility,omitempty"`
	Modifiers      []*ModifierInvocation `json:"modifiers,omitempty"`
	Override       []Node              `json:"override,omitempty"`
	IsConstructor  bool                `json:"isConstructor"`
	IsFallback     bool                `json:"isFallback"`
	IsReceiveEther bool                `json:"isReceiveEther"`
	IsVirtual      bool                `json:"isVirtual"`
	StateMutability string             `json:"stateMutability,omitempty"`
}

// ModifierDefinition represents a modifier definition
type ModifierDefinition struct {
	BaseNode
	Name       string                 `json:"name"`
	Parameters []*VariableDeclaration `json:"parameters,omitempty"`
	Body       *Block                 `json:"body,omitempty"`
	IsVirtual  bool                   `json:"isVirtual"`
	Override   []Node                 `json:"override,omitempty"`
}

// ModifierInvocation represents a modifier invocation
type ModifierInvocation struct {
	BaseNode
	Name      string `json:"name"`
	Arguments []Node `json:"arguments,omitempty"`
}

// StateVariableDeclaration represents a state variable
type StateVariableDeclaration struct {
	BaseNode
	Variables  []*VariableDeclaration `json:"variables"`
	InitialValue Node                 `json:"initialValue,omitempty"`
}

// VariableDeclaration represents a variable declaration
type VariableDeclaration struct {
	BaseNode
	TypeName        Node   `json:"typeName"`
	Name            string `json:"name,omitempty"`
	Identifier      *Identifier `json:"identifier,omitempty"`
	StorageLocation string `json:"storageLocation,omitempty"`
	IsStateVar      bool   `json:"isStateVar"`
	IsIndexed       bool   `json:"isIndexed"`
	IsImmutable     bool   `json:"isImmutable"`
	Override        []Node `json:"override,omitempty"`
	Visibility      string `json:"visibility,omitempty"`
	IsDeclaredConst bool   `json:"isDeclaredConst,omitempty"`
	Expression      Node   `json:"expression,omitempty"`
}

// VariableDeclarationStatement represents a variable declaration statement
type VariableDeclarationStatement struct {
	BaseNode
	Variables    []*VariableDeclaration `json:"variables"`
	InitialValue Node                   `json:"initialValue,omitempty"`
}

// StructDefinition represents a struct definition
type StructDefinition struct {
	BaseNode
	Name    string                 `json:"name"`
	Members []*VariableDeclaration `json:"members"`
}

// EnumDefinition represents an enum definition
type EnumDefinition struct {
	BaseNode
	Name    string       `json:"name"`
	Members []*EnumValue `json:"members"`
}

// EnumValue represents an enum value
type EnumValue struct {
	BaseNode
	Name string `json:"name"`
}

// EventDefinition represents an event definition
type EventDefinition struct {
	BaseNode
	Name        string                 `json:"name"`
	Parameters  []*VariableDeclaration `json:"parameters"`
	IsAnonymous bool                   `json:"isAnonymous"`
}

// ErrorDefinition represents a custom error definition
type ErrorDefinition struct {
	BaseNode
	Name       string                 `json:"name"`
	Parameters []*VariableDeclaration `json:"parameters"`
}

// UserDefinedValueTypeDefinition represents a user-defined value type
type UserDefinedValueTypeDefinition struct {
	BaseNode
	Name        string `json:"name"`
	UnderlyingType Node `json:"underlyingType"`
}

// UsingForDeclaration represents a using-for directive
type UsingForDeclaration struct {
	BaseNode
	TypeName   Node     `json:"typeName,omitempty"`
	Functions  []string `json:"functions,omitempty"`
	Operators  []string `json:"operators,omitempty"`
	LibraryName string  `json:"libraryName,omitempty"`
	IsGlobal   bool     `json:"isGlobal"`
}

// ElementaryTypeName represents a built-in type
type ElementaryTypeName struct {
	BaseNode
	Name         string `json:"name"`
	StateMutability string `json:"stateMutability,omitempty"`
}

// UserDefinedTypeName represents a user-defined type
type UserDefinedTypeName struct {
	BaseNode
	NamePath string `json:"namePath"`
}

// Mapping represents a mapping type
type Mapping struct {
	BaseNode
	KeyType      Node   `json:"keyType"`
	KeyName      *Identifier `json:"keyName,omitempty"`
	ValueType    Node   `json:"valueType"`
	ValueName    *Identifier `json:"valueName,omitempty"`
}

// ArrayTypeName represents an array type
type ArrayTypeName struct {
	BaseNode
	BaseTypeName Node `json:"baseTypeName"`
	Length       Node `json:"length,omitempty"`
}

// FunctionTypeName represents a function type
type FunctionTypeName struct {
	BaseNode
	ParameterTypes      []*VariableDeclaration `json:"parameterTypes"`
	ReturnTypes         []*VariableDeclaration `json:"returnTypes,omitempty"`
	Visibility          string                 `json:"visibility,omitempty"`
	StateMutability     string                 `json:"stateMutability,omitempty"`
}

// Block represents a block of statements
type Block struct {
	BaseNode
	Statements []Node `json:"statements"`
}

// UncheckedBlock represents an unchecked block
type UncheckedBlock struct {
	BaseNode
	Body *Block `json:"body"`
}

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	BaseNode
	Expression Node `json:"expression"`
}

// IfStatement represents an if statement
type IfStatement struct {
	BaseNode
	Condition Node `json:"condition"`
	TrueBody  Node `json:"trueBody"`
	FalseBody Node `json:"falseBody,omitempty"`
}

// WhileStatement represents a while loop
type WhileStatement struct {
	BaseNode
	Condition Node `json:"condition"`
	Body      Node `json:"body"`
}

// DoWhileStatement represents a do-while loop
type DoWhileStatement struct {
	BaseNode
	Condition Node `json:"condition"`
	Body      Node `json:"body"`
}

// ForStatement represents a for loop
type ForStatement struct {
	BaseNode
	InitExpression      Node `json:"initExpression,omitempty"`
	ConditionExpression Node `json:"conditionExpression,omitempty"`
	LoopExpression      Node `json:"loopExpression,omitempty"`
	Body                Node `json:"body"`
}

// ContinueStatement represents a continue statement
type ContinueStatement struct {
	BaseNode
}

// BreakStatement represents a break statement
type BreakStatement struct {
	BaseNode
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	BaseNode
	Expression Node `json:"expression,omitempty"`
}

// EmitStatement represents an emit statement
type EmitStatement struct {
	BaseNode
	EventCall Node `json:"eventCall"`
}

// RevertStatement represents a revert statement
type RevertStatement struct {
	BaseNode
	RevertCall Node `json:"revertCall,omitempty"`
}

// TryStatement represents a try-catch statement
type TryStatement struct {
	BaseNode
	Expression       Node           `json:"expression"`
	ReturnParameters []*VariableDeclaration `json:"returnParameters,omitempty"`
	Body             *Block         `json:"body"`
	CatchClauses     []*CatchClause `json:"catchClauses"`
}

// CatchClause represents a catch clause
type CatchClause struct {
	BaseNode
	IsReasonStringType bool                   `json:"isReasonStringType"`
	Kind               string                 `json:"kind,omitempty"`
	Parameters         []*VariableDeclaration `json:"parameters,omitempty"`
	Body               *Block                 `json:"body"`
}

// BinaryOperation represents a binary operation
type BinaryOperation struct {
	BaseNode
	Operator string `json:"operator"`
	Left     Node   `json:"left"`
	Right    Node   `json:"right"`
}

// UnaryOperation represents a unary operation
type UnaryOperation struct {
	BaseNode
	Operator   string `json:"operator"`
	SubExpression Node `json:"subExpression"`
	IsPrefix   bool   `json:"isPrefix"`
}

// Conditional represents a ternary conditional expression
type Conditional struct {
	BaseNode
	Condition   Node `json:"condition"`
	TrueExpression  Node `json:"trueExpression"`
	FalseExpression Node `json:"falseExpression"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	BaseNode
	Expression Node   `json:"expression"`
	Arguments  []Node `json:"arguments"`
	Names      []string `json:"names,omitempty"`
	Identifiers []*Identifier `json:"identifiers,omitempty"`
}

// FunctionCallOptions represents function call options (e.g., {value: 1, gas: 100})
type FunctionCallOptions struct {
	BaseNode
	Expression Node     `json:"expression"`
	Names      []string `json:"names"`
	Options    []Node   `json:"options"`
}

// MemberAccess represents member access (e.g., foo.bar)
type MemberAccess struct {
	BaseNode
	Expression Node   `json:"expression"`
	MemberName string `json:"memberName"`
}

// IndexAccess represents index access (e.g., arr[i])
type IndexAccess struct {
	BaseNode
	Base  Node `json:"base"`
	Index Node `json:"index,omitempty"`
}

// IndexRangeAccess represents slice access (e.g., arr[1:3])
type IndexRangeAccess struct {
	BaseNode
	Base       Node `json:"base"`
	IndexStart Node `json:"indexStart,omitempty"`
	IndexEnd   Node `json:"indexEnd,omitempty"`
}

// NewExpression represents a new expression
type NewExpression struct {
	BaseNode
	TypeName Node `json:"typeName"`
}

// TupleExpression represents a tuple expression
type TupleExpression struct {
	BaseNode
	Components []Node `json:"components"`
	IsArray    bool   `json:"isArray"`
}

// NameValueExpression represents a named value expression
type NameValueExpression struct {
	BaseNode
	Expression Node              `json:"expression"`
	Arguments  *NameValueList    `json:"arguments"`
}

// NameValueList represents a list of named values
type NameValueList struct {
	BaseNode
	Names       []string          `json:"names"`
	Identifiers []*Identifier     `json:"identifiers"`
	Arguments   []Node            `json:"arguments"`
}

// Identifier represents an identifier
type Identifier struct {
	BaseNode
	Name string `json:"name"`
}

// NumberLiteral represents a numeric literal
type NumberLiteral struct {
	BaseNode
	Number      string `json:"number"`
	SubDenomination string `json:"subdenomination,omitempty"`
}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	BaseNode
	Value bool `json:"value"`
}

// StringLiteral represents a string literal
type StringLiteral struct {
	BaseNode
	Value    string   `json:"value"`
	Parts    []string `json:"parts,omitempty"`
	IsUnicode bool    `json:"isUnicode"`
}

// HexLiteral represents a hex literal
type HexLiteral struct {
	BaseNode
	Value string   `json:"value"`
	Parts []string `json:"parts,omitempty"`
}

// InlineAssembly represents inline assembly
type InlineAssembly struct {
	BaseNode
	Language string         `json:"language,omitempty"`
	Body     *AssemblyBlock `json:"body"`
}

// AssemblyBlock represents an assembly block
type AssemblyBlock struct {
	BaseNode
	Operations []Node `json:"operations"`
}

// AssemblyCall represents an assembly function call
type AssemblyCall struct {
	BaseNode
	FunctionName string `json:"functionName"`
	Arguments    []Node `json:"arguments"`
}

// AssemblyLocalDefinition represents a local variable definition in assembly
type AssemblyLocalDefinition struct {
	BaseNode
	Names      []*Identifier `json:"names"`
	Expression Node          `json:"expression,omitempty"`
}

// AssemblyAssignment represents an assignment in assembly
type AssemblyAssignment struct {
	BaseNode
	Names      []*Identifier `json:"names"`
	Expression Node          `json:"expression"`
}

// AssemblyIdentifier represents an identifier in assembly
type AssemblyIdentifier struct {
	BaseNode
	Name string `json:"name"`
}

// AssemblyLiteral represents a literal in assembly
type AssemblyLiteral struct {
	BaseNode
	Kind  string `json:"kind"` // "number", "string", "boolean"
	Value string `json:"value"`
}

// AssemblyIf represents an if statement in assembly
type AssemblyIf struct {
	BaseNode
	Condition Node           `json:"condition"`
	Body      *AssemblyBlock `json:"body"`
}

// AssemblySwitch represents a switch statement in assembly
type AssemblySwitch struct {
	BaseNode
	Expression Node            `json:"expression"`
	Cases      []*AssemblyCase `json:"cases"`
}

// AssemblyCase represents a case in assembly switch
type AssemblyCase struct {
	BaseNode
	Value   Node           `json:"value,omitempty"` // nil for default case
	Body    *AssemblyBlock `json:"body"`
	Default bool           `json:"default"`
}

// AssemblyFor represents a for loop in assembly
type AssemblyFor struct {
	BaseNode
	Pre       *AssemblyBlock `json:"pre"`
	Condition Node           `json:"condition"`
	Post      *AssemblyBlock `json:"post"`
	Body      *AssemblyBlock `json:"body"`
}

// AssemblyFunctionDefinition represents a function definition in assembly
type AssemblyFunctionDefinition struct {
	BaseNode
	Name       string        `json:"name"`
	Arguments  []*Identifier `json:"arguments,omitempty"`
	ReturnArguments []*Identifier `json:"returnArguments,omitempty"`
	Body       *AssemblyBlock `json:"body"`
}

// MarshalJSON implements custom JSON marshaling for SourceUnit
func (s *SourceUnit) MarshalJSON() ([]byte, error) {
	type Alias SourceUnit
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SourceUnit
func (s *SourceUnit) UnmarshalJSON(data []byte) error {
	type Alias SourceUnit
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	return json.Unmarshal(data, aux)
}

