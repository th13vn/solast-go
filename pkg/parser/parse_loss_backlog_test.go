package parser

import (
	"encoding/json"
	"testing"
)

// The tests in this file come from a real-corpus scan (507 verified contracts)
// that surfaced 166 tolerant-mode parse diagnostics across 23 files. Each test
// pins one root cause with a minimal repro reduced from the production source
// named in its comment, and each one also asserts that the declaration AFTER
// the construct survives: in tolerant mode `expect()` does not advance on a
// mismatch, so a single unhandled token can desynchronize the parser and
// silently drop the rest of the contract or file.

// decodeAST parses cleanly (see parseClean) and returns the AST as generic
// JSON so a test can assert on node structure instead of substrings.
func decodeAST(t *testing.T, src string) any {
	t.Helper()
	var root any
	if err := json.Unmarshal([]byte(parseClean(t, src)), &root); err != nil {
		t.Fatalf("decode AST: %v", err)
	}
	return root
}

// findNodes returns every JSON object in the tree whose "type" is typ, in
// document order.
func findNodes(v any, typ string) []map[string]any {
	var out []map[string]any
	var walk func(any)
	walk = func(v any) {
		switch n := v.(type) {
		case map[string]any:
			if n["type"] == typ {
				out = append(out, n)
			}
			for _, child := range n {
				walk(child)
			}
		case []any:
			for _, child := range n {
				walk(child)
			}
		}
	}
	walk(v)
	return out
}

// hasNode reports whether a node of type typ has field key equal to want.
func hasNode(root any, typ, key, want string) bool {
	for _, n := range findNodes(root, typ) {
		if s, ok := n[key].(string); ok && s == want {
			return true
		}
	}
	return false
}

func requireAfterFunction(t *testing.T, root any) {
	t.Helper()
	if !hasNode(root, "FunctionDefinition", "name", "afterF") {
		t.Fatalf("the declaration after the construct was dropped, so the parser desynced")
	}
}

// TestYulBuiltinsSpelledLikeSolidityKeywords guards the largest cascade in the
// corpus. Inside assembly, `address()`, `revert(p, n)` and `return(p, n)` are
// ordinary EVM builtins, but the lexer classifies `address`, `revert` and
// `return` as Solidity keywords, and the Yul parser only accepted IDENTIFIER
// tokens. `mstore(0x14, address())` then failed on the `(` after the literal
// `address` and desynchronized the rest of the file (solady ERC20.sol,
// euler-vault-kit Dispatch.sol / BorrowUtils.sol); `revert(fmp, 0x44)`
// desynchronized v4-core CustomRevert.sol. With literal-only arguments the
// same statements were dropped with NO recovered error at all, so every
// `return(0, 0x20)` in assembly silently vanished from the AST.
//
// Solidity-only keywords are plain Yul identifiers, so the same rule also
// covers contextual keywords referenced from assembly (`shl(96, from)`) and
// the external-function-pointer members `g.address` / `g.selector`.
func TestYulBuiltinsSpelledLikeSolidityKeywords(t *testing.T) {
	cases := []struct {
		name string
		asm  string
		typ  string
		key  string
		want string
	}{
		{"address_call_as_argument", `mstore(0x14, address())`, "AssemblyCall", "functionName", "address"},
		{"address_call_in_let", `let a := address()`, "AssemblyCall", "functionName", "address"},
		{"revert_statement_with_identifier", `let p := mload(0x40) revert(p, 0x44)`, "AssemblyCall", "functionName", "revert"},
		{"revert_statement_literal_args", `revert(0, 0)`, "AssemblyCall", "functionName", "revert"},
		{"return_statement_literal_args", `return(0, 0x20)`, "AssemblyCall", "functionName", "return"},
		{"return_in_switch_default", `switch r case 0 { revert(0, returndatasize()) } default { return(0, returndatasize()) }`, "AssemblyCall", "functionName", "return"},
		{"contextual_keyword_reference", `let x := shl(96, from)`, "AssemblyIdentifier", "name", "from"},
		{"contextual_keyword_let_name", `let error := 1 mstore(0, error)`, "AssemblyIdentifier", "name", "error"},
		{"function_pointer_address_member", `mstore(0, g.address)`, "AssemblyIdentifier", "name", "g.address"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `pragma solidity ^0.8.20;
contract C {
    function f(address from, uint256 r, function() external g) internal view {
        assembly ("memory-safe") {
            ` + tc.asm + `
        }
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			root := decodeAST(t, src)
			if !hasNode(root, tc.typ, tc.key, tc.want) {
				t.Fatalf("no %s with %s=%q in the AST; the Yul operation was dropped", tc.typ, tc.key, tc.want)
			}
			requireAfterFunction(t, root)
		})
	}
}

// TestYulFunctionReturnVariables guards `function h(a) -> b { ... }` inside
// assembly. The parser looked for `=>` (ARROW) instead of the `->`
// (RIGHT_ARROW) token the lexer produces, so every Yul function with return
// variables failed at the arrow and swallowed the rest of the file. Not seen in
// the backlog corpus; found while probing the Yul parser for the cases above.
func TestYulFunctionReturnVariables(t *testing.T) {
	src := `pragma solidity ^0.8.20;
contract C {
    function f(uint256 v) internal pure returns (uint256 r) {
        assembly {
            function inc(a) -> b { b := add(a, 1) }
            function pair(a) -> x, y { x := a y := a }
            r := inc(v)
        }
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
	root := decodeAST(t, src)
	defs := findNodes(root, "AssemblyFunctionDefinition")
	if len(defs) != 2 {
		t.Fatalf("got %d Yul function definitions, want 2", len(defs))
	}
	for i, want := range []int{1, 2} {
		rets, _ := defs[i]["returnArguments"].([]any)
		if len(rets) != want {
			t.Fatalf("Yul function %v has %d return variables, want %d", defs[i]["name"], len(rets), want)
		}
	}
	if !hasNode(root, "AssemblyCall", "functionName", "inc") {
		t.Fatalf("the statement after the Yul function definitions was dropped")
	}
	requireAfterFunction(t, root)
}

// TestUsingForQualifiedFunctionList guards `using { Lib.fn, ... } for T global`
// – the grammar's identifierPath entries – as written by every prb-math
// ValueType.sol (`using { Casting.intoSD59x18, ... } for SD1x18 global;`,
// `using { Helpers.add as +, ... } for SD59x18 global;`). The list parser read
// a single IDENTIFIER per entry, died on the `.`, and in tolerant mode lost
// every following using directive and the rest of the file.
func TestUsingForQualifiedFunctionList(t *testing.T) {
	src := `pragma solidity >=0.8.19;
import "./Casting.sol" as Casting;
import "./Helpers.sol" as Helpers;
type SD1x18 is int64;
using { Casting.intoSD59x18, Casting.unwrap } for SD1x18 global;
using { Helpers.add as +, Helpers.sub as - } for SD1x18 global;
using { plain } for SD1x18 global;
contract C {
    using { Helpers.add } for SD1x18;
    function afterF() public pure returns (uint256) { return 7; }
}`
	root := decodeAST(t, src)
	usings := findNodes(root, "UsingForDeclaration")
	if len(usings) != 4 {
		t.Fatalf("got %d using directives, want 4", len(usings))
	}
	wantFuncs := [][]string{
		{"Casting.intoSD59x18", "Casting.unwrap"},
		{"Helpers.add", "Helpers.sub"},
		{"plain"},
		{"Helpers.add"},
	}
	for i, want := range wantFuncs {
		got, _ := usings[i]["functions"].([]any)
		if len(got) != len(want) {
			t.Fatalf("using #%d functions = %v, want %v", i, got, want)
		}
		for j := range want {
			if got[j] != want[j] {
				t.Fatalf("using #%d functions = %v, want %v", i, got, want)
			}
		}
	}
	ops, _ := usings[1]["operators"].([]any)
	if len(ops) != 2 || ops[0] != "+" || ops[1] != "-" {
		t.Fatalf("operators = %v, want [+ -]", ops)
	}
	if usings[0]["isGlobal"] != true || usings[3]["isGlobal"] != false {
		t.Fatalf("global flag lost: %v / %v", usings[0]["isGlobal"], usings[3]["isGlobal"])
	}
	requireAfterFunction(t, root)
}

// TestErrorParameterContextualKeywordName guards custom-error parameters named
// with a contextual keyword, e.g. Axelar's
// `error TokenManagerDeploymentFailed(bytes error);`. Function and event
// parameters already accepted contextual keywords as names; error parameters
// only took IDENTIFIER, so the `error` name failed `expect(',')` and the
// desync dropped the rest of IInterchainTokenService (44 diagnostics).
func TestErrorParameterContextualKeywordName(t *testing.T) {
	src := `pragma solidity ^0.8.0;
interface I {
    error TokenManagerDeploymentFailed(bytes error);
    error Moved(address from, address to);
    event Deployed(bytes32 indexed id);
    function afterF() external pure returns (uint256);
}`
	root := decodeAST(t, src)
	errs := findNodes(root, "ErrorDefinition")
	if len(errs) != 2 {
		t.Fatalf("got %d error definitions, want 2", len(errs))
	}
	params, _ := errs[0]["parameters"].([]any)
	if len(params) != 1 || params[0].(map[string]any)["name"] != "error" {
		t.Fatalf("error parameter named `error` not kept: %v", params)
	}
	params, _ = errs[1]["parameters"].([]any)
	if len(params) != 2 || params[0].(map[string]any)["name"] != "from" {
		t.Fatalf("error parameter named `from` not kept: %v", params)
	}
	if !hasNode(root, "EventDefinition", "name", "Deployed") {
		t.Fatalf("the event after the error definitions was dropped")
	}
	requireAfterFunction(t, root)
}

// TestForInitExpressionStatement guards for-loop init clauses that are
// expressions starting with an identifier, as in prb-math Math.sol
// `for (yAux >>= 1; yAux > 0; yAux >>= 1)`. The init clause chose the
// declaration path whenever the first token could start a type name, so any
// identifier-led expression was parsed as a nameless declaration: a compound
// assignment such as `>>=` then failed on the missing `;`, and a plain
// `for (i = 0; ...)` was silently recorded as a declaration of type `i`. The
// init clause now uses the same lookahead as a statement.
func TestForInitExpressionStatement(t *testing.T) {
	cases := []struct {
		name     string
		init     string
		wantInit string
	}{
		{"compound_shift_assignment", `y >>= 1`, "ExpressionStatement"},
		{"plain_assignment", `i = 0`, "ExpressionStatement"},
		{"member_assignment", `s.i = 0`, "ExpressionStatement"},
		{"elementary_declaration", `uint256 j = 0`, "VariableDeclarationStatement"},
		{"user_type_declaration", `Lib.T memory t = s`, "VariableDeclarationStatement"},
		{"tuple_declaration", `(uint256 a, uint256 b) = (1, 2)`, "VariableDeclarationStatement"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `pragma solidity ^0.8.20;
contract C {
    function f(uint256 y, uint256 i, S memory s) internal pure returns (uint256 r) {
        for (` + tc.init + `; y > 0; y >>= 1) { r += 1; }
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			root := decodeAST(t, src)
			loops := findNodes(root, "ForStatement")
			if len(loops) != 1 {
				t.Fatalf("got %d for statements, want 1", len(loops))
			}
			init, _ := loops[0]["initExpression"].(map[string]any)
			if init == nil || init["type"] != tc.wantInit {
				t.Fatalf("init clause = %v, want %s", init, tc.wantInit)
			}
			requireAfterFunction(t, root)
		})
	}
}

// TestLocalDeclarationQualifiedArrayType guards local declarations whose type
// is a qualified name with array dimensions, as in Curve Swaps.sol
// `Storage.Assimilator[] memory _reserves = curve.assets;`. The statement
// lookahead skipped array brackets BEFORE the `.member` path, the reverse of
// the grammar's order, so it met `[` where it expected a location or name and
// sent the declaration down the expression path, which died on `memory`.
func TestLocalDeclarationQualifiedArrayType(t *testing.T) {
	decls := []struct {
		name string
		body string
		want string
	}{
		{"qualified_dynamic_array_memory", `Storage.Assimilator[] memory rs = s.assets; rs;`, "rs"},
		{"qualified_fixed_array_memory", `Lib.T[2] memory pair; pair;`, "pair"},
		{"qualified_nested_array_storage", `Lib.T[][] storage grid = s.grid; grid;`, "grid"},
		{"deep_qualified_array", `A.B.C[] memory deep; deep;`, "deep"},
	}
	for _, tc := range decls {
		t.Run(tc.name, func(t *testing.T) {
			src := `pragma solidity ^0.8.20;
contract C {
    function f(S storage s) internal view {
        ` + tc.body + `
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			root := decodeAST(t, src)
			stmts := findNodes(root, "VariableDeclarationStatement")
			if len(stmts) != 1 {
				t.Fatalf("got %d declaration statements, want 1", len(stmts))
			}
			vars, _ := stmts[0]["variables"].([]any)
			if len(vars) != 1 || vars[0].(map[string]any)["name"] != tc.want {
				t.Fatalf("declared variables = %v, want %s", vars, tc.want)
			}
			requireAfterFunction(t, root)
		})
	}

	// The other side of the lookahead: indexed and member targets are still
	// assignments, not declarations.
	exprs := []string{`a.b[i] = x;`, `a[i].b = x;`, `a.b[i].c = x;`}
	for _, body := range exprs {
		t.Run("expression_"+body, func(t *testing.T) {
			src := `pragma solidity ^0.8.20;
contract C {
    function f(uint256 i, uint256 x) internal {
        ` + body + `
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			root := decodeAST(t, src)
			if n := len(findNodes(root, "VariableDeclarationStatement")); n != 0 {
				t.Fatalf("assignment %q parsed as %d declaration(s)", body, n)
			}
			if len(findNodes(root, "ExpressionStatement")) == 0 {
				t.Fatalf("assignment %q produced no expression statement", body)
			}
			requireAfterFunction(t, root)
		})
	}
}

// TestTupleAssignmentIndexedMemberTargets guards tuple assignments whose
// components are member accesses on an indexed value, as in
// ethereum-vault-connector EthereumVaultConnector.sol
// `(batchItemsResult[i].success, batchItemsResult[i].result) = ...;` and
// Set.sol `(setStorage.firstElement, setStorage.elements[index2].value) = ...;`.
// The speculative tuple-declaration parse stopped at the `.` after `[i]`, then
// backtracked to the position AFTER the opening `(` rather than before it, so
// the expression parser saw `items[i].ok, ...` without its parenthesis and
// failed on the comma.
func TestTupleAssignmentIndexedMemberTargets(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"indexed_member_pair", `(items[i].ok, items[i].data) = g();`},
		{"member_then_indexed_member", `(s.first, s.elems[j].value) = (s.elems[j].value, s.first);`},
		{"call_result_member", `(items[i].ok, x) = (true, 1);`},
		{"nested_tuple", `((a, b), c) = h();`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `pragma solidity ^0.8.20;
contract C {
    function f(uint256 i, uint256 j) internal {
        ` + tc.body + `
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			root := decodeAST(t, src)
			if !hasNode(root, "BinaryOperation", "operator", "=") {
				t.Fatalf("tuple assignment %q produced no assignment", tc.body)
			}
			if len(findNodes(root, "TupleExpression")) == 0 {
				t.Fatalf("tuple assignment %q lost its tuple", tc.body)
			}
			requireAfterFunction(t, root)
		})
	}

	// A tuple that DOES declare variables must stay a declaration.
	src := `pragma solidity ^0.8.20;
contract C {
    function f() internal {
        (bool ok, bytes memory data) = g(); ok; data;
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
	root := decodeAST(t, src)
	if len(findNodes(root, "VariableDeclarationStatement")) != 1 {
		t.Fatalf("tuple declaration no longer parsed as a declaration")
	}
}
