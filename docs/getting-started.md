# Getting Started

This guide will help you get started with solast-go, whether you're using it as a CLI tool or as a Go package.

## Installation

### CLI Tool

```bash
go install github.com/th13vn/solast-go/cmd/solast@latest
```

### Go Package

```bash
go get github.com/th13vn/solast-go
```

## Quick Start with CLI

### Parse a Solidity File

```bash
# Parse and output JSON AST
solast parse contract.sol

# Pretty print with location info
solast parse contract.sol --loc --range

# Save output to file
solast parse contract.sol -o ast.json
```

### Validate Syntax

```bash
solast validate contract.sol
# Output: Syntax OK
```

### Detect Version

```bash
solast version-detect contract.sol
# Output:
# Pragma: ^0.8.20
# Version: 0.8.20
# Constraint: ^
```

### Parse from Stdin

```bash
echo 'contract Test {}' | solast parse -
```

## Quick Start with Go Package

### Basic Parsing

```go
package main

import (
    "fmt"
    "log"

    "github.com/th13vn/solast-go/pkg/parser"
)

func main() {
    source := `
        pragma solidity ^0.8.0;
        contract Counter {
            uint256 public count;
            function increment() public {
                count++;
            }
        }
    `

    ast, err := parser.Parse(source, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Parsed %d top-level elements\n", len(ast.Children))
}
```

### Parsing with Options

```go
ast, err := parser.Parse(source, &parser.Options{
    Tolerant: true,  // Continue parsing despite errors
    Loc:      true,  // Include line/column information
    Range:    true,  // Include character range
})
```

### Get JSON Output

```go
jsonBytes, err := parser.ParseToJSON(source, nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(jsonBytes))
```

### Walking the AST

```go
import "github.com/th13vn/solast-go/pkg/ast"

visitor := &ast.SimpleVisitor{
    ContractDefinitionFn: func(node *ast.ContractDefinition) {
        fmt.Printf("Contract: %s\n", node.Name)
    },
    FunctionDefinitionFn: func(node *ast.FunctionDefinition) {
        fmt.Printf("  Function: %s\n", node.Name)
    },
    StateVariableDeclarationFn: func(node *ast.StateVariableDeclaration) {
        for _, v := range node.Variables {
            fmt.Printf("  Variable: %s\n", v.Name)
        }
    },
}

parser.VisitSimple(ast, visitor)
```

### Version Detection

```go
import "github.com/th13vn/solast-go/pkg/version"

detected, err := version.Detect(source)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Pragma: %s\n", detected.Raw)
fmt.Printf("Version: %s\n", detected.Version)
fmt.Printf("Constraint: %s\n", detected.Constraint)
```

### Version Comparison

```go
v1 := version.MustParse("0.8.0")
v2 := version.MustParse("0.8.20")

if v1.LessThan(v2) {
    fmt.Println("v1 is older than v2")
}

switch v1.Compare(v2) {
case -1:
    fmt.Println("v1 < v2")
case 0:
    fmt.Println("v1 == v2")
case 1:
    fmt.Println("v1 > v2")
}
```

## Common Use Cases

### Find All Functions in a Contract

```go
var functions []string

visitor := &ast.SimpleVisitor{
    FunctionDefinitionFn: func(node *ast.FunctionDefinition) {
        functions = append(functions, node.Name)
    },
}
parser.VisitSimple(ast, visitor)

fmt.Printf("Functions: %v\n", functions)
```

### Find All Events

```go
visitor := &ast.SimpleVisitor{
    EventDefinitionFn: func(node *ast.EventDefinition) {
        fmt.Printf("Event: %s\n", node.Name)
        for _, param := range node.Parameters {
            indexed := ""
            if param.IsIndexed {
                indexed = " (indexed)"
            }
            fmt.Printf("  - %s%s\n", param.Name, indexed)
        }
    },
}
```

### Find All State Variables

```go
visitor := &ast.SimpleVisitor{
    StateVariableDeclarationFn: func(node *ast.StateVariableDeclaration) {
        for _, v := range node.Variables {
            fmt.Printf("State var: %s (visibility: %s)\n", 
                v.Name, v.Visibility)
        }
    },
}
```

### Check Contract Inheritance

```go
visitor := &ast.SimpleVisitor{
    ContractDefinitionFn: func(node *ast.ContractDefinition) {
        fmt.Printf("Contract: %s\n", node.Name)
        for _, base := range node.BaseContracts {
            fmt.Printf("  inherits: %s\n", base.BaseName.NamePath)
        }
    },
}
```

## Error Handling

### Tolerant Mode

In tolerant mode, the parser continues even when encountering syntax errors:

```go
ast, err := parser.Parse(badSource, &parser.Options{Tolerant: true})
// err will be nil, but ast may be incomplete
```

### Parser Errors

```go
ast, err := parser.Parse(source, nil)
if err != nil {
    if parserErr, ok := err.(*parser.ParserError); ok {
        for _, e := range parserErr.Errors {
            fmt.Printf("Error at line %d: %s\n", e.Line, e.Message)
        }
    }
}
```

## Next Steps

- [Overview](overview.md) - Architecture and AST node types
- [Development](development.md) - Development setup
- [Contributing](contributing.md) - How to contribute

