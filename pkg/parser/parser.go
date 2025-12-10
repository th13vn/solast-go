// Package parser provides a Solidity parser that produces an AST compatible
// with the TypeScript solidity-parser (https://github.com/solidity-parser/parser).
package parser

import (
	"encoding/json"
	"io"

	"github.com/th13vn/solast-go/internal/builder"
	"github.com/th13vn/solast-go/pkg/ast"
)

// Options configures the parser behavior
type Options struct {
	// Tolerant mode: collect errors instead of stopping on first error
	Tolerant bool
	// Loc: add location information (line/column) to nodes
	Loc bool
	// Range: add character range information to nodes
	Range bool
}

// ParserError represents a parsing error
type ParserError struct {
	Errors []*Error
}

func (e *ParserError) Error() string {
	if len(e.Errors) == 0 {
		return "parsing error"
	}
	return e.Errors[0].Error()
}

// Error represents a single parsing error
type Error struct {
	Message string `json:"message"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}

func (e *Error) Error() string {
	return e.Message
}

// Parse parses Solidity source code and returns an AST
func Parse(input string, opts *Options) (*ast.SourceUnit, error) {
	if opts == nil {
		opts = &Options{}
	}

	b := builder.New(input, &builder.Options{
		Tolerant: opts.Tolerant,
		Loc:      opts.Loc,
		Range:    opts.Range,
	})

	result, err := b.Build()
	if err != nil {
		builderErr := err.(*builder.Error)
		return nil, &ParserError{
			Errors: []*Error{{
				Message: builderErr.Message,
				Line:    builderErr.Line,
				Column:  builderErr.Column,
			}},
		}
	}

	// Check for collected errors in tolerant mode
	if len(b.Errors()) > 0 && !opts.Tolerant {
		var errors []*Error
		for _, e := range b.Errors() {
			errors = append(errors, &Error{
				Message: e.Message,
				Line:    e.Line,
				Column:  e.Column,
			})
		}
		return nil, &ParserError{Errors: errors}
	}

	return result, nil
}

// ParseReader parses Solidity source from an io.Reader and returns an AST
func ParseReader(r io.Reader, opts *Options) (*ast.SourceUnit, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Parse(string(content), opts)
}

// ParseToJSON parses Solidity source code and returns JSON
func ParseToJSON(input string, opts *Options) ([]byte, error) {
	result, err := Parse(input, opts)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(result, "", "  ")
}

// Visit walks the AST and calls the appropriate visitor method for each node
func Visit(node ast.Node, visitor ast.Visitor) {
	ast.Walk(node, visitor)
}

// VisitSimple walks the AST using a SimpleVisitor for easier use
func VisitSimple(node ast.Node, visitor *ast.SimpleVisitor) {
	ast.WalkSimple(node, visitor)
}

// Visitor is an alias for ast.Visitor
type Visitor = ast.Visitor

// BaseVisitor is an alias for ast.BaseVisitor
type BaseVisitor = ast.BaseVisitor

// SimpleVisitor is an alias for ast.SimpleVisitor
type SimpleVisitor = ast.SimpleVisitor

