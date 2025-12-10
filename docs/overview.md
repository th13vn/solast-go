# Overview

Solast-go is a Solidity Abstract Syntax Tree (AST) parser written in Go. It transforms Solidity source code into a structured tree representation that can be analyzed, transformed, or used for building development tools.

## What is an AST?

An Abstract Syntax Tree is a tree representation of source code structure. Each node in the tree represents a construct in the source code:

```
SourceUnit
├── PragmaDirective (pragma solidity ^0.8.0)
└── ContractDefinition (contract MyContract)
    ├── StateVariableDeclaration (uint256 value)
    └── FunctionDefinition (function setValue)
        └── Block
            └── ExpressionStatement
                └── BinaryOperation (value = _value)
```

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Source    │ ──▶ │   Lexer     │ ──▶ │   Parser    │
│   Code      │     │  (Tokens)   │     │  (Builder)  │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                                              ▼
                                        ┌─────────────┐
                                        │    AST      │
                                        │   Nodes     │
                                        └─────────────┘
```

### Components

1. **Lexer** (`internal/lexer`)
   - Tokenizes Solidity source code
   - Handles keywords, identifiers, literals, operators
   - Supports all Solidity versions

2. **Parser/Builder** (`internal/builder`)
   - Recursive descent parser
   - Builds AST from token stream
   - Handles expressions with operator precedence

3. **AST Nodes** (`pkg/ast`)
   - Go structs representing syntax tree nodes
   - Matches TypeScript solidity-parser format
   - Includes visitor pattern for traversal

4. **Parser API** (`pkg/parser`)
   - Public API for parsing
   - Options for location tracking, tolerant mode
   - JSON serialization support

5. **Version Detection** (`pkg/version`)
   - Detect Solidity version from pragma
   - Version comparison utilities

## AST Node Types

### Top-Level Nodes

| Node Type            | Description                                    | Example                   |
| -------------------- | ---------------------------------------------- | ------------------------- |
| `SourceUnit`         | Root node containing all top-level definitions | Entire file               |
| `PragmaDirective`    | Pragma statement                               | `pragma solidity ^0.8.0;` |
| `ImportDirective`    | Import statement                               | `import "./Other.sol";`   |
| `ContractDefinition` | Contract, interface, or library                | `contract MyContract {}`  |

### Contract Members

| Node Type                  | Description           | Example                    |
| -------------------------- | --------------------- | -------------------------- |
| `FunctionDefinition`       | Function declaration  | `function foo() public {}` |
| `ModifierDefinition`       | Modifier declaration  | `modifier onlyOwner() {}`  |
| `StateVariableDeclaration` | State variable        | `uint256 public value;`    |
| `StructDefinition`         | Struct declaration    | `struct Person { ... }`    |
| `EnumDefinition`           | Enum declaration      | `enum Status { ... }`      |
| `EventDefinition`          | Event declaration     | `event Transfer(...);`     |
| `ErrorDefinition`          | Custom error (0.8.4+) | `error Unauthorized();`    |

### Statements

| Node Type         | Description                   |
| ----------------- | ----------------------------- |
| `Block`           | Block of statements `{ ... }` |
| `IfStatement`     | If/else statement             |
| `ForStatement`    | For loop                      |
| `WhileStatement`  | While loop                    |
| `ReturnStatement` | Return statement              |
| `EmitStatement`   | Event emission                |
| `RevertStatement` | Revert with error             |
| `TryStatement`    | Try/catch block               |
| `UncheckedBlock`  | Unchecked arithmetic (0.8+)   |

### Expressions

| Node Type         | Description             |
| ----------------- | ----------------------- |
| `BinaryOperation` | `a + b`, `a == b`, etc. |
| `UnaryOperation`  | `!a`, `-b`, `++i`       |
| `FunctionCall`    | `foo(arg1, arg2)`       |
| `MemberAccess`    | `obj.member`            |
| `IndexAccess`     | `arr[0]`                |
| `Identifier`      | Variable name           |
| `NumberLiteral`   | `123`, `1 ether`        |
| `StringLiteral`   | `"hello"`               |
| `BooleanLiteral`  | `true`, `false`         |

## Use Cases

- **Static Analysis** - Find security vulnerabilities, code smells
- **Code Generation** - Generate documentation, interfaces
- **Linting** - Check code style and conventions
- **Refactoring** - Automated code transformations
- **IDE Support** - Syntax highlighting, go-to-definition
- **Testing** - Contract analysis and verification

## Next Steps

- [Getting Started](getting-started.md) - Quick start guide
- [Development](development.md) - Development setup
- [Contributing](contributing.md) - How to contribute

