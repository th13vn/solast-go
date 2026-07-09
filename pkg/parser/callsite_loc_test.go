package parser

import (
	"testing"

	"github.com/th13vn/solast-go/pkg/ast"
)

// TestFunctionCallLocation ensures that FunctionCall nodes built in
// parseCallMemberIndex get their Loc/Range populated when the parser is
// invoked with Loc/Range enabled. Previously these postfix expression nodes
// (FunctionCall, MemberAccess, IndexAccess, IndexRangeAccess,
// FunctionCallOptions) never had setLocation called on them.
func TestFunctionCallLocation(t *testing.T) {
	input := `pragma solidity ^0.8.0;
contract C {
    function f() public { g(1); }
    function g(uint256 x) public {}
}
`

	result, err := Parse(input, &Options{Loc: true, Range: true})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var found *ast.FunctionCall
	VisitSimple(result, &SimpleVisitor{
		FunctionCallFn: func(n *ast.FunctionCall) {
			if callee, ok := n.Expression.(*ast.Identifier); ok && callee.Name == "g" {
				found = n
			}
		},
	})

	if found == nil {
		t.Fatal("did not find FunctionCall node for g(1)")
	}

	if found.Loc == nil {
		t.Fatal("expected FunctionCall.Loc to be set, got nil")
	}
	if found.Loc.Start.Line != 3 {
		t.Errorf("expected FunctionCall.Loc.Start.Line == 3, got %d", found.Loc.Start.Line)
	}
	if found.Range == nil {
		t.Fatal("expected FunctionCall.Range to be set, got nil")
	}
}
