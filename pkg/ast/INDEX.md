# pkg/ast — AST Node Types & Visitor

## Purpose

The public, serializable AST. Every node the [[builder]] produces is a struct here. Consumers (e.g. w3goaudit) import this package and traverse it with the visitor. Adding new syntax means adding a node type AND wiring it into the visitor walkers.

## nodes.go (~650 lines)

**Node interface** (nodes.go:116):
```go
type Node interface {
    GetType() NodeType
    GetLocation() *Location
    GetRange() *Range
}
```

**BaseNode** (nodes.go:123) — embedded in every node: `Type NodeType`, `Loc *Location`, `Range *Range`.

**Location/Range** (nodes.go:101-114): `Location{Start, End Position}`, `Position{Line, Column int}`, `Range [2]int` (byte offsets).

**NodeType** (nodes.go:9) — a `string` type; ~87 `Node*` string constants (nodes.go:12-99). Add one per new node.

**Node structs by category:**
- **Top-level / directives**: `SourceUnit{Children []Node}`, `PragmaDirective`, `ImportDirective` (+ `ImportSymbol`, `ImportSymbolIdentifiers`).
- **Definitions**: `ContractDefinition{Name, Kind, BaseContracts, SubNodes}`, `InheritanceSpecifier`, `FunctionDefinition{Name, Parameters, ReturnParameters, Body, Visibility, StateMutability, Modifiers, IsConstructor/IsFallback/IsReceiveEther/IsVirtual}`, `ModifierDefinition`, `StructDefinition{Members []*VariableDeclaration}`, `EnumDefinition`/`EnumValue`, `EventDefinition`, `ErrorDefinition`, `UserDefinedValueTypeDefinition`, `UsingForDeclaration`.
- **Variables**: `StateVariableDeclaration`, `VariableDeclaration{TypeName, Name, StorageLocation, IsStateVar/IsIndexed/IsImmutable/IsDeclaredConst, Visibility, Expression}`.
- **Type names**: `ElementaryTypeName`, `UserDefinedTypeName{NamePath}`, `Mapping{KeyType, ValueType, KeyName, ValueName}`, `ArrayTypeName{BaseTypeName, Length}`, `FunctionTypeName`.
- **Statements**: `Block`, `UncheckedBlock`, `ExpressionStatement`, `IfStatement`, `WhileStatement`, `DoWhileStatement`, `ForStatement`, `Continue/Break/Return/Emit/Revert Statement`, `TryStatement`, `CatchClause`.
- **Expressions**: `BinaryOperation`, `UnaryOperation`, `Conditional`, `FunctionCall{Expression, Arguments, Names, Identifiers}`, `FunctionCallOptions`, `MemberAccess`, `IndexAccess`, `IndexRangeAccess`, `NewExpression`, `TupleExpression`, `NameValueExpression`/`NameValueList`, `Identifier`, `NumberLiteral{Number, SubDenomination}`, `BooleanLiteral`, `StringLiteral{Value, Parts, IsUnicode}`, `HexLiteral`.
- **Assembly (Yul)**: `InlineAssembly`, `AssemblyBlock`, `AssemblyCall`, `AssemblyLocalDefinition`, `AssemblyAssignment`, `AssemblyIdentifier`, `AssemblyLiteral`, `AssemblyIf`, `AssemblySwitch`/`AssemblyCase`, `AssemblyFor`, `AssemblyFunctionDefinition`.
- **Misc**: `ModifierInvocation`, `ParameterList`, `Parameter`, `EventParameter`.

> JSON note: nodes serialize to JSON (CLI `parse` and w3goaudit caching rely on it). Keep field tags stable; renaming a field is a breaking change for consumers.

## visitor.go (~950 lines)

- **Visitor interface** (visitor.go:4) — one `Visit<Node>(*Node) bool` per node type (~66). Return `false` to stop descent.
- **BaseVisitor** (visitor.go:74) — no-op defaults (all return `true`); embed it to override only what you need.
- **SimpleVisitor** (visitor.go:143) — embeds `BaseVisitor`, exposes a `<Node>Fn func(*Node)` callback field per type; always descends.
- **Walk(node, Visitor)** (visitor.go:213) and **WalkSimple(node, *SimpleVisitor)** (visitor.go:543) — recursive traversal with a big per-node switch.

## Change checklist (new node type)

1. Add a `Node<Name>` constant + the struct (embed `BaseNode`, add `GetType/GetLocation/GetRange` if not provided by BaseNode pattern).
2. Add `Visit<Name>` to `Visitor`, a default to `BaseVisitor`, a `<Name>Fn` to `SimpleVisitor`.
3. Add a `case` for the node in BOTH `Walk` and `WalkSimple` (descend into its child `Node` fields).
4. Add a `setLocation` case in [[builder]] helpers.go.
