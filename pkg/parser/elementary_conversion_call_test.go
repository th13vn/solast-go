package parser

import (
	"encoding/json"
	"strings"
	"testing"
)

// parseClean parses in tolerant mode – the mode tooling actually uses – and
// fails if the parser recovered from anything. A recovered error here means a
// statement was dropped, which is silent analysis loss for every consumer.
func parseClean(t *testing.T, src string) string {
	t.Helper()
	unit, errs, err := ParseWithErrors(src, &Options{Tolerant: true})
	if err != nil {
		t.Fatalf("ParseWithErrors returned a hard error: %v", err)
	}
	if len(errs) != 0 {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msgs = append(msgs, e.Error())
		}
		t.Fatalf("recovered %d parse error(s): %s", len(errs), strings.Join(msgs, "; "))
	}
	encoded, err := json.Marshal(unit)
	if err != nil {
		t.Fatalf("marshal AST: %v", err)
	}
	return string(encoded)
}

// TestElementaryTypeConversionCallStatement guards the regression where a
// statement beginning with an elementary type conversion – `address(x).call(b)`
// – was routed to the variable-declaration parser, which died on the `(` and
// dropped the entire call from the enclosing function.
//
// The statement-level lookahead treated every leading elementary type name as a
// declaration. That is wrong: no variable declaration can put `(` directly
// after the type, so `address(` is unambiguously a conversion expression.
//
// This matters far beyond the literal shape. `address(target).call`,
// `.delegatecall`, `.staticcall`, and the whole OpenZeppelin `Address` library
// idiom are written this way in production Solidity, and every dropped call is
// invisible to a consumer's call graph, taint analysis, and detectors.
func TestElementaryTypeConversionCallStatement(t *testing.T) {
	cases := []struct {
		name   string
		body   string
		member string
	}{
		{"address_call", `address(x).call(b);`, "call"},
		{"address_delegatecall", `address(x).delegatecall(b);`, "delegatecall"},
		{"address_staticcall", `address(x).staticcall(b);`, "staticcall"},
		{"address_this_call", `address(this).call("");`, "call"},
		{"address_call_with_value", `address(x).call{value: 1}("");`, "call"},
		{"oz_address_library_idiom", `address(x).functionDelegateCall(b, "err");`, "functionDelegateCall"},
		{"assigned_result", `(bool ok, ) = address(x).call(b); ok;`, "call"},
		{"uint_conversion_member_call", `uint256(n).toString();`, "toString"},
		{"bytes32_conversion_member_call", `bytes32(h).toHexString();`, "toHexString"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `pragma solidity ^0.8.20;
contract C {
    function f(address x, bytes memory b, uint256 n, bytes32 h) public {
        ` + tc.body + `
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			encoded := parseClean(t, src)
			if !strings.Contains(encoded, `"`+tc.member+`"`) {
				t.Fatalf("member %q missing from the AST; the call was dropped", tc.member)
			}
			if !strings.Contains(encoded, `"afterF"`) {
				t.Fatalf("the following function was dropped, so the parser desynced")
			}
		})
	}
}

// TestElementaryTypeDeclarationsStillParse pins the other side of the lookahead
// change: every ordinary declaration that begins with an elementary type name
// must still be parsed as a declaration, not as an expression.
func TestElementaryTypeDeclarationsStillParse(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"plain_address", `address a = x;`},
		{"address_payable", `address payable p = payable(x);`},
		{"address_array_memory", `address[] memory arr = new address[](1); arr;`},
		{"uint256", `uint256 v = n;`},
		{"bytes32", `bytes32 k = h;`},
		{"bytes_memory", `bytes memory data = b;`},
		{"string_memory", `string memory s = "a";`},
		{"bool", `bool flag = true;`},
		{"uninitialized", `uint256 z;`},
		{"tuple_destructure", `(uint256 q, bytes32 w) = (n, h); q; w;`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `pragma solidity ^0.8.20;
contract C {
    function f(address x, bytes memory b, uint256 n, bytes32 h) public {
        ` + tc.body + `
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			encoded := parseClean(t, src)
			// Parameters are VariableDeclaration nodes too, so only the
			// statement node proves the body took the declaration path.
			if !strings.Contains(encoded, `"VariableDeclarationStatement"`) {
				t.Fatalf("declaration was not parsed as a declaration: %s", encoded[:min(400, len(encoded))])
			}
			if !strings.Contains(encoded, `"afterF"`) {
				t.Fatalf("the following function was dropped, so the parser desynced")
			}
		})
	}
}
