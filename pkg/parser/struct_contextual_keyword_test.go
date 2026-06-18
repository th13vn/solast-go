package parser

import (
	"testing"

	"github.com/th13vn/solast-go/pkg/ast"
)

// countFns returns the number of FunctionDefinition subnodes in the first
// contract of a parsed source unit.
func countFns(t *testing.T, src string) int {
	t.Helper()
	res, err := Parse(src, &Options{Tolerant: true})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	var contract *ast.ContractDefinition
	for _, child := range res.Children {
		if c, ok := child.(*ast.ContractDefinition); ok {
			contract = c
			break
		}
	}
	if contract == nil {
		t.Fatal("no contract found")
	}
	n := 0
	for _, sub := range contract.SubNodes {
		if _, ok := sub.(*ast.FunctionDefinition); ok {
			n++
		}
	}
	return n
}

// TestStructFieldContextualKeyword guards the regression where a struct (or
// enum) member named with a contextual keyword such as `from` desynced the
// tolerant parser and silently consumed the rest of the enclosing contract —
// dropping every following function. A struct ending in a `from` field is a
// common shape (e.g. `struct Transfer { uint256 amount; address from; }`), so
// the body MUST survive regardless of where `from` sits.
func TestStructFieldContextualKeyword(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{"from_only_field", `pragma solidity ^0.8.20;
contract A { struct S { address from; } function f() external pure returns (uint256) { return 1; } }`},
		{"from_last_field", `pragma solidity ^0.8.20;
contract A { struct S { address to; address from; } function f() external pure returns (uint256) { return 1; } }`},
		{"from_first_field", `pragma solidity ^0.8.20;
contract A { struct S { address from; address to; } function f() external pure returns (uint256) { return 1; } }`},
		{"uint_from_field", `pragma solidity ^0.8.20;
contract A { struct S { uint256 from; } function f() external pure returns (uint256) { return 1; } }`},
		{"enum_from_value", `pragma solidity ^0.8.20;
contract A { enum E { from } function f() external pure returns (uint256) { return 1; } }`},
		{"enum_from_last_value", `pragma solidity ^0.8.20;
contract A { enum E { a, from } function f() external pure returns (uint256) { return 1; } }`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := countFns(t, tc.src); got != 1 {
				t.Fatalf("expected the contract function to survive, got %d functions (parser desynced on contextual-keyword member)", got)
			}
		})
	}
}

// TestParseWithErrorsExposesTolerantErrors verifies that ParseWithErrors surfaces
// recovered errors in tolerant mode (which plain Parse discards), and reports a
// clean parse as an empty error slice.
func TestParseWithErrorsExposesTolerantErrors(t *testing.T) {
	clean := `pragma solidity ^0.8.20;
contract A { function f() external pure returns (uint256) { return 1; } }`
	if _, errs, err := ParseWithErrors(clean, &Options{Tolerant: true}); err != nil || len(errs) != 0 {
		t.Fatalf("clean source: err=%v errs=%d, want no error and 0 recovered errors", err, len(errs))
	}

	// Malformed: missing type before the member, and a stray token — tolerant mode
	// recovers but should report the errors rather than swallow them.
	broken := `pragma solidity ^0.8.20;
contract A { function f( external { uint256 x = ; } }`
	res, errs, err := ParseWithErrors(broken, &Options{Tolerant: true})
	if err != nil {
		t.Fatalf("tolerant parse should not hard-fail: %v", err)
	}
	if res == nil {
		t.Fatal("tolerant parse should still return an AST")
	}
	if len(errs) == 0 {
		t.Fatal("ParseWithErrors must surface recovered errors in tolerant mode")
	}
}
