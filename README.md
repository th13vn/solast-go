# Solast-go

A Solidity Abstract Syntax Tree (AST) parser written in Go. Parse Solidity smart contracts into structured AST for analysis, transformation, and tooling.

*This project was developed with assistance from [Cursor](https://cursor.com), an AI-powered code editor.*

## Features

- **CLI Tool** - Parse Solidity files from command line with JSON output
- **Go Package** - Import and use in your Go projects
- **Multi-Version Support** - Supports Solidity 0.4.x through 0.8.x syntax
- **Location Tracking** - Optional line/column and character range information
- **Tolerant Mode** - Continue parsing despite errors
- **Version Detection** - Detect Solidity version from pragma directives

## Installation

### CLI Tool

```bash
go install github.com/th13vn/solast-go/cmd/solast@latest
```

### Go Package

```bash
go get github.com/th13vn/solast-go
```

## CLI Usage

```bash
# Parse a Solidity file and output JSON AST
solast parse contract.sol

# Parse with location information
solast parse contract.sol --loc --range

# Parse from stdin
cat contract.sol | solast parse -

# Validate syntax only (no AST output)
solast validate contract.sol

# Detect Solidity version
solast version-detect contract.sol

# Output to file
solast parse contract.sol -o output.json

# Tolerant mode (continue parsing despite errors)
solast parse contract.sol --tolerant
```

## Package Usage

```go
package main

import (
    "fmt"
    "github.com/th13vn/solast-go/pkg/parser"
    "github.com/th13vn/solast-go/pkg/ast"
)

func main() {
    input := `
        pragma solidity ^0.8.0;
        
        contract MyContract {
            uint256 public value;
            
            function setValue(uint256 _value) public {
                value = _value;
            }
        }
    `

    // Parse with options
    result, err := parser.Parse(input, &parser.Options{
        Tolerant: false,
        Loc:      true,
        Range:    true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Parsed %d children\n", len(result.Children))

    // Visit AST nodes
    visitor := &ast.SimpleVisitor{
        ContractDefinitionFn: func(node *ast.ContractDefinition) {
            fmt.Printf("Found contract: %s\n", node.Name)
        },
        FunctionDefinitionFn: func(node *ast.FunctionDefinition) {
            fmt.Printf("Found function: %s\n", node.Name)
        },
    }
    parser.VisitSimple(result, visitor)
}
```

## API Overview

### Parser Package

```go
// Parse Solidity source code
func Parse(input string, opts *Options) (*ast.SourceUnit, error)

// Parse and return JSON
func ParseToJSON(input string, opts *Options) ([]byte, error)

// Options for parsing
type Options struct {
    Tolerant bool  // Continue parsing despite errors
    Loc      bool  // Include line/column location
    Range    bool  // Include character range
}
```

### Version Package

```go
// Parse a version string
func Parse(s string) (Version, error)

// Create a version
func New(major, minor, patch int) Version

// Detect version from source code
func Detect(source string) (*DetectedVersion, error)

// Compare versions
v1.Compare(v2)        // -1, 0, or 1
v1.LessThan(v2)       // bool
v1.GreaterThan(v2)    // bool
v1.Equal(v2)          // bool
```

### AST Node Types

- **SourceUnit** - Root node
- **PragmaDirective** - `pragma solidity ^0.8.0;`
- **ImportDirective** - `import "./Other.sol";`
- **ContractDefinition** - Contract, interface, or library
- **FunctionDefinition** - Function declarations
- **VariableDeclaration** - Variable declarations
- **StructDefinition** - Struct definitions
- **EnumDefinition** - Enum definitions
- **EventDefinition** - Event definitions
- **ErrorDefinition** - Custom error definitions
- **ModifierDefinition** - Modifier definitions
- **Statement types** - Block, If, While, For, Return, etc.
- **Expression types** - BinaryOperation, FunctionCall, MemberAccess, etc.

## Supported Solidity Features

| Version | Features                                                                              |
| ------- | ------------------------------------------------------------------------------------- |
| 0.4.x   | Basic syntax, structs, enums, events, modifiers                                       |
| 0.5.x   | `constructor`, `emit`, `address payable`, `calldata`                                  |
| 0.6.x   | `abstract`, `virtual`/`override`, `try`/`catch`, `receive`/`fallback`                 |
| 0.7.x   | Free functions, file-level `using`, `gwei`, `immutable`                               |
| 0.8.x   | `unchecked`, custom errors, user-defined types, named mappings, `transient`, `layout` |

## Development

```bash
# Run tests
go test ./...

# Build CLI (with version info)
make build

# Build CLI (dev, quick)
go build -o solast ./cmd/solast

# Run CLI
./solast parse contract.sol

# Check version
./solast --version
```

### Release

```bash
# Build release binaries for all platforms
make release
# Creates dist/solast-{os}-{arch} binaries
```

See [docs/development.md](docs/development.md) for detailed development guide.

## References

- [solidity-parser/parser](https://github.com/solidity-parser/parser) - TypeScript Solidity parser (AST format reference)
- [Solidity Documentation](https://docs.soliditylang.org/) - Official Solidity documentation

## License

MIT License - see [LICENSE](LICENSE) file for details.
