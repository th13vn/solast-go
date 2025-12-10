# Contributing

Thank you for your interest in contributing to solast-go! This document provides guidelines for contributing to the project.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the issue, not the person
- Help others learn and grow

## How to Contribute

### Reporting Issues

1. Search existing issues to avoid duplicates
2. Use a clear, descriptive title
3. Describe the problem and expected behavior
4. Include Solidity code that reproduces the issue
5. Specify your Go version and OS

**Issue Template:**

```markdown
## Description
Brief description of the issue

## Steps to Reproduce
1. Parse the following Solidity code:
```solidity
// Your code here
```

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- Go version: 
- OS: 
- solast-go version:
```

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Commit with clear messages
7. Push to your fork
8. Open a Pull Request

**PR Checklist:**

- [ ] Tests added/updated
- [ ] Documentation updated if needed
- [ ] Code follows existing style
- [ ] All tests pass
- [ ] Commit messages are clear

### Commit Message Format

```
<type>: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `test`: Adding or updating tests
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `style`: Code style changes (formatting, etc.)
- `chore`: Maintenance tasks

**Examples:**

```
feat: add support for Solidity 0.8.25 transient storage

fix: correct parsing of nested mappings

docs: update getting started guide

test: add tests for try/catch parsing
```

## Development Setup

See [development.md](development.md) for detailed setup instructions.

Quick start:

```bash
# Clone
git clone https://github.com/th13vn/solast-go.git
cd solast-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o solast ./cmd/solast
```

## Code Style

### Go Code

- Follow standard Go conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions focused and small
- Use meaningful variable names

```go
// Good
func parseContractDefinition(kind string) *ast.ContractDefinition {
    // ...
}

// Bad
func pcd(k string) *ast.ContractDefinition {
    // ...
}
```

### Tests

- Test file naming: `*_test.go`
- Use table-driven tests where appropriate
- Test both success and error cases
- Use descriptive test names

```go
func TestParseFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:  "public function",
            input: `function foo() public {}`,
        },
        {
            name:    "invalid syntax",
            input:   `function {}`,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Adding New Solidity Features

When Solidity adds new syntax:

1. **Update Lexer** (`internal/lexer/lexer.go`)
   - Add new token types if needed
   - Update keyword map

2. **Update Parser** (`internal/builder/`)
   - Add parsing logic for new syntax
   - Handle in appropriate parsing function

3. **Update AST Nodes** (`pkg/ast/nodes.go`)
   - Add new node types if needed
   - Add fields to existing nodes

4. **Add Tests** (`pkg/parser/parser_test.go`)
   - Add test cases for new syntax
   - Test edge cases

5. **Update Documentation**
   - Update README.md
   - Update docs/overview.md

### Example: Adding New Keyword

```go
// 1. Add token type in lexer.go
const (
    // ...
    NEWKEYWORD TokenType = iota
    // ...
)

// 2. Add to keywords map
var keywords = map[string]TokenType{
    // ...
    "newkeyword": NEWKEYWORD,
}

// 3. Handle in parser (builder.go or appropriate file)
case lexer.NEWKEYWORD:
    return b.parseNewKeywordFeature()

// 4. Add AST node if needed (nodes.go)
type NewFeatureNode struct {
    BaseNode
    // fields
}

// 5. Add tests
func TestNewKeyword(t *testing.T) {
    // ...
}
```

## Questions?

- Open an issue for questions
- Check existing issues and docs first
- Be specific about what you're trying to do

Thank you for contributing!

