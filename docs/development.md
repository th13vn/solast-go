# Development Guide

This guide covers setting up a development environment and understanding the codebase structure.

## Prerequisites

- Go 1.21 or later
- Git

## Setup

```bash
# Clone the repository
git clone https://github.com/th13vn/solast-go.git
cd solast-go

# Download dependencies
go mod download

# Verify setup
go test ./...
```

## Project Structure

```
solast-go/
├── cmd/
│   └── solast/
│       └── main.go          # CLI entry point
├── internal/
│   ├── lexer/
│   │   ├── lexer.go         # Tokenizer
│   │   └── lexer_test.go
│   └── builder/
│       ├── builder.go       # Main parser
│       ├── expressions.go   # Expression parsing
│       ├── statements.go    # Statement parsing
│       ├── types.go         # Type parsing
│       └── helpers.go       # Utility functions
├── pkg/
│   ├── ast/
│   │   ├── nodes.go         # AST node definitions
│   │   └── visitor.go       # Visitor pattern
│   ├── parser/
│   │   ├── parser.go        # Public API
│   │   └── parser_test.go
│   └── version/
│       ├── version.go       # Version detection
│       └── version_test.go
├── testdata/
│   └── contracts/           # Test Solidity files
├── docs/                    # Documentation
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

## Building

```bash
# Build CLI with version info (recommended)
make build

# Build CLI (dev, quick, no version info)
go build -o solast ./cmd/solast

# Check version
./solast --version
# Output: solast version 0.1.0 (commit: abc1234, built: 2024-01-01T00:00:00Z)

# Build release binaries for all platforms
make release
# Creates: dist/solast-{linux,darwin,windows}-{amd64,arm64}
```

### Version Information

Version is read from the `VERSION` file and embedded at build time:

```bash
# Show current version
make version

# Update version (edit VERSION file)
echo "0.2.0" > VERSION
make build
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./pkg/parser/...

# Run specific test
go test -v -run TestParseSimpleContract ./pkg/parser/...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Understanding the Code

### Lexer (`internal/lexer/lexer.go`)

The lexer converts source code into tokens:

```go
// Token types
const (
    IDENTIFIER TokenType = iota
    NUMBER
    STRING
    CONTRACT
    FUNCTION
    // ...
)

// Lexer usage
lex := lexer.New(sourceCode)
tokens := lex.Tokenize()

// Each token has:
// - Type: token type (keyword, identifier, etc.)
// - Value: actual string value
// - Line, Column: position in source
// - Start, End: character offsets
```

### Parser/Builder (`internal/builder/`)

The parser uses recursive descent to build the AST:

```go
// Main entry point
func (b *Builder) Build() (*ast.SourceUnit, error) {
    sourceUnit := &ast.SourceUnit{...}
    
    for !b.isAtEnd() {
        node := b.parseSourceUnitElement()
        sourceUnit.Children = append(sourceUnit.Children, node)
    }
    
    return sourceUnit, nil
}

// Parse different constructs
func (b *Builder) parseContractDefinition(kind string) *ast.ContractDefinition
func (b *Builder) parseFunctionDefinition() *ast.FunctionDefinition
func (b *Builder) parseExpression() ast.Node
func (b *Builder) parseStatement() ast.Node
```

**Expression Parsing with Precedence:**

```go
// Precedence levels (lowest to highest)
// 1. Assignment (=, +=, etc.)
// 2. Conditional (?:)
// 3. Logical OR (||)
// 4. Logical AND (&&)
// 5. Equality (==, !=)
// ... and so on

func (b *Builder) parseExpression() ast.Node {
    return b.parseAssignment()
}

func (b *Builder) parseAssignment() ast.Node {
    left := b.parseTernary()
    if b.isAssignmentOperator() {
        op := b.advance().Value
        right := b.parseAssignment()
        return &ast.BinaryOperation{...}
    }
    return left
}
```

### AST Nodes (`pkg/ast/nodes.go`)

Each node type is a Go struct:

```go
type ContractDefinition struct {
    BaseNode
    Name          string
    Kind          string  // "contract", "interface", "library"
    BaseContracts []*InheritanceSpecifier
    SubNodes      []Node
}

type FunctionDefinition struct {
    BaseNode
    Name             string
    Parameters       []*VariableDeclaration
    ReturnParameters []*VariableDeclaration
    Visibility       string
    StateMutability  string
    Body             *Block
    // ...
}
```

### Visitor Pattern (`pkg/ast/visitor.go`)

For traversing the AST:

```go
type SimpleVisitor struct {
    ContractDefinitionFn func(*ContractDefinition)
    FunctionDefinitionFn func(*FunctionDefinition)
    // ... other node types
}

// Usage
visitor := &ast.SimpleVisitor{
    FunctionDefinitionFn: func(n *ast.FunctionDefinition) {
        fmt.Println(n.Name)
    },
}
parser.VisitSimple(ast, visitor)
```

## Adding Support for New Solidity Features

### Step 1: Update Lexer

If the feature introduces new keywords:

```go
// internal/lexer/lexer.go

// Add token type
const (
    // ...
    NEWKEYWORD TokenType = iota
)

// Add to keywords map
var keywords = map[string]TokenType{
    // ...
    "newkeyword": NEWKEYWORD,
}

// Add string representation
var tokenStrings = map[TokenType]string{
    // ...
    NEWKEYWORD: "newkeyword",
}
```

### Step 2: Update Parser

Add parsing logic:

```go
// internal/builder/builder.go or appropriate file

func (b *Builder) parseNewFeature() *ast.NewFeatureNode {
    startTok := b.advance() // consume keyword
    
    // Parse the feature
    // ...
    
    node := &ast.NewFeatureNode{
        BaseNode: ast.BaseNode{Type: ast.NodeNewFeature},
        // ... set fields
    }
    
    b.setLocation(node, startTok, b.previous())
    return node
}
```

### Step 3: Add AST Node

```go
// pkg/ast/nodes.go

const (
    // ...
    NodeNewFeature NodeType = "NewFeature"
)

type NewFeatureNode struct {
    BaseNode
    // Add fields as needed
}

func (n *NewFeatureNode) GetType() NodeType { return NodeNewFeature }
func (n *NewFeatureNode) GetLocation() *Location { return n.Loc }
func (n *NewFeatureNode) GetRange() *Range { return n.Range }
```

### Step 4: Update Visitor

```go
// pkg/ast/visitor.go

type SimpleVisitor struct {
    // ...
    NewFeatureNodeFn func(*NewFeatureNode)
}
```

### Step 5: Add Tests

```go
// pkg/parser/parser_test.go

func TestNewFeature(t *testing.T) {
    input := `pragma solidity ^0.9.0; contract Test { newkeyword ... }`
    
    result, err := Parse(input, nil)
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    // Assert expected structure
}
```

### Step 6: Add Test Contract

Create a test file in `testdata/contracts/`:

```solidity
// testdata/contracts/v09/NewFeature.sol
pragma solidity ^0.9.0;

contract NewFeatureExample {
    // Use the new feature
}
```

## Debugging

### Print Tokens

```go
lex := lexer.New(source)
tokens := lex.Tokenize()
for i, tok := range tokens {
    fmt.Printf("%d: %s = %q\n", i, tok.Type, tok.Value)
}
```

### Verbose Test Output

```bash
go test -v -run TestSpecificTest ./pkg/parser/...
```

### Debug with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug tests
dlv test ./pkg/parser/... -- -test.run TestParseSimpleContract

# Debug CLI
dlv debug ./cmd/solast -- parse contract.sol
```

## Performance

The parser is designed for correctness over speed, but performance considerations:

- Avoid allocations in hot paths
- Use string builders for concatenation
- Preallocate slices when size is known

```bash
# Run benchmarks
go test -bench=. ./pkg/parser/...

# Profile
go test -cpuprofile cpu.prof -memprofile mem.prof -bench=. ./pkg/parser/...
go tool pprof cpu.prof
```

## Release Process

1. Update version in `VERSION` file
2. Update CHANGELOG.md (if exists)
3. Commit changes: `git commit -am "Release v0.2.0"`
4. Create git tag: `git tag v0.2.0`
5. Push changes and tag:
   ```bash
   git push origin main
   git push origin v0.2.0
   ```
6. Build release binaries: `make release`
7. Go module proxy will automatically index the new version

## Questions?

- Open an issue
- See [contributing.md](contributing.md)

