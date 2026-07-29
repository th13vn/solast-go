package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/th13vn/solast-go/pkg/ast"
)

// =============================================================================
// Basic Parsing Tests
// =============================================================================

func TestParseSimpleContract(t *testing.T) {
	input := `
		pragma solidity ^0.8.0;
		
		contract SimpleStorage {
			uint256 public value;
			
			function setValue(uint256 _value) public {
				value = _value;
			}
			
			function getValue() public view returns (uint256) {
				return value;
			}
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.Type != ast.NodeSourceUnit {
		t.Errorf("Expected SourceUnit type, got %s", result.Type)
	}

	if len(result.Children) < 2 {
		t.Errorf("Expected at least 2 children, got %d", len(result.Children))
	}

	pragma, ok := result.Children[0].(*ast.PragmaDirective)
	if !ok {
		t.Error("First child should be PragmaDirective")
	} else if pragma.Name != "solidity" {
		t.Errorf("Expected pragma name 'solidity', got '%s'", pragma.Name)
	}

	contract, ok := result.Children[1].(*ast.ContractDefinition)
	if !ok {
		t.Error("Second child should be ContractDefinition")
	} else {
		if contract.Name != "SimpleStorage" {
			t.Errorf("Expected contract name 'SimpleStorage', got '%s'", contract.Name)
		}
		if contract.Kind != "contract" {
			t.Errorf("Expected contract kind 'contract', got '%s'", contract.Kind)
		}
	}
}

func TestParseWithLocation(t *testing.T) {
	input := `pragma solidity ^0.8.0;`

	result, err := Parse(input, &Options{Loc: true})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	pragma := result.Children[0].(*ast.PragmaDirective)
	if pragma.Loc == nil {
		t.Error("Location should be set")
	} else if pragma.Loc.Start.Line != 1 {
		t.Errorf("Expected start line 1, got %d", pragma.Loc.Start.Line)
	}
}

func TestParseWithRange(t *testing.T) {
	input := `pragma solidity ^0.8.0;`

	result, err := Parse(input, &Options{Range: true})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	pragma := result.Children[0].(*ast.PragmaDirective)
	if pragma.Range == nil {
		t.Error("Range should be set")
	}
}

func TestTolerantMode(t *testing.T) {
	input := `contract Test { invalid syntax here }`

	_, err := Parse(input, nil)
	if err == nil {
		t.Error("Expected error without tolerant mode")
	}

	_, err = Parse(input, &Options{Tolerant: true})
	if err != nil {
		t.Errorf("Tolerant mode should not return error: %v", err)
	}
}

func TestJSONOutput(t *testing.T) {
	input := `contract Test {}`

	jsonOutput, err := ParseToJSON(input, nil)
	if err != nil {
		t.Fatalf("ParseToJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonOutput, &result); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	if result["type"] != "SourceUnit" {
		t.Errorf("Expected type 'SourceUnit', got '%v'", result["type"])
	}
}

// =============================================================================
// Import Tests
// =============================================================================

func TestParseImport(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple import", `import "./Other.sol";`, "./Other.sol"},
		{"import with alias", `import "./Other.sol" as Other;`, "./Other.sol"},
		{"import all", `import * as Other from "./Other.sol";`, "./Other.sol"},
		{"import specific", `import { Symbol } from "./Other.sol";`, "./Other.sol"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			imp, ok := result.Children[0].(*ast.ImportDirective)
			if !ok {
				t.Fatal("Expected ImportDirective")
			}

			if imp.Path != tt.expected {
				t.Errorf("Expected path '%s', got '%s'", tt.expected, imp.Path)
			}
		})
	}
}

// =============================================================================
// Contract Element Tests
// =============================================================================

func TestParseFunction(t *testing.T) {
	input := `
		contract Test {
			function publicFunc() public {}
			function privateFunc() private pure returns (uint256) { return 0; }
			function externalFunc(uint256 a, string memory b) external view {}
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	if len(contract.SubNodes) != 3 {
		t.Errorf("Expected 3 functions, got %d", len(contract.SubNodes))
	}

	fn1 := contract.SubNodes[0].(*ast.FunctionDefinition)
	if fn1.Name != "publicFunc" {
		t.Errorf("Expected name 'publicFunc', got '%s'", fn1.Name)
	}
	if fn1.Visibility != "public" {
		t.Errorf("Expected visibility 'public', got '%s'", fn1.Visibility)
	}
}

func TestParseStruct(t *testing.T) {
	input := `
		contract Test {
			struct Person {
				string name;
				uint256 age;
				address wallet;
			}
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	structDef := contract.SubNodes[0].(*ast.StructDefinition)

	if structDef.Name != "Person" {
		t.Errorf("Expected name 'Person', got '%s'", structDef.Name)
	}
	if len(structDef.Members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(structDef.Members))
	}
}

func TestParseEnum(t *testing.T) {
	input := `
		contract Test {
			enum Status { Pending, Active, Completed }
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	enumDef := contract.SubNodes[0].(*ast.EnumDefinition)

	if enumDef.Name != "Status" {
		t.Errorf("Expected name 'Status', got '%s'", enumDef.Name)
	}
	if len(enumDef.Members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(enumDef.Members))
	}
}

func TestParseEvent(t *testing.T) {
	input := `
		contract Test {
			event Transfer(address indexed from, address indexed to, uint256 value);
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	eventDef := contract.SubNodes[0].(*ast.EventDefinition)

	if eventDef.Name != "Transfer" {
		t.Errorf("Expected name 'Transfer', got '%s'", eventDef.Name)
	}
	if len(eventDef.Parameters) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(eventDef.Parameters))
	}
	if !eventDef.Parameters[0].IsIndexed {
		t.Error("First parameter should be indexed")
	}
}

func TestParseMapping(t *testing.T) {
	input := `
		contract Test {
			mapping(address => uint256) public balances;
			mapping(address => mapping(address => uint256)) public allowances;
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	if len(contract.SubNodes) != 2 {
		t.Errorf("Expected 2 state variables, got %d", len(contract.SubNodes))
	}
}

func TestParseInheritance(t *testing.T) {
	input := `contract Child is Parent, Ownable {}`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	if len(contract.BaseContracts) != 2 {
		t.Errorf("Expected 2 base contracts, got %d", len(contract.BaseContracts))
	}
}

func TestParseInterface(t *testing.T) {
	input := `
		interface IERC20 {
			function totalSupply() external view returns (uint256);
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	if contract.Kind != "interface" {
		t.Errorf("Expected kind 'interface', got '%s'", contract.Kind)
	}
}

func TestParseLibrary(t *testing.T) {
	input := `
		library SafeMath {
			function add(uint256 a, uint256 b) internal pure returns (uint256) {
				return a + b;
			}
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	contract := result.Children[0].(*ast.ContractDefinition)
	if contract.Kind != "library" {
		t.Errorf("Expected kind 'library', got '%s'", contract.Kind)
	}
}

// =============================================================================
// Visitor Tests
// =============================================================================

func TestVisitor(t *testing.T) {
	input := `
		contract Test {
			function foo() public {}
			function bar() private {}
		}
	`

	result, err := Parse(input, nil)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var functionNames []string
	visitor := &ast.SimpleVisitor{
		FunctionDefinitionFn: func(node *ast.FunctionDefinition) {
			functionNames = append(functionNames, node.Name)
		},
	}

	VisitSimple(result, visitor)

	if len(functionNames) != 2 {
		t.Errorf("Expected 2 functions, found %d", len(functionNames))
	}
	if functionNames[0] != "foo" {
		t.Errorf("Expected first function 'foo', got '%s'", functionNames[0])
	}
}

// =============================================================================
// Solidity Version Tests
// =============================================================================

func TestSolidity04x(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "constructor as function name",
			input: `pragma solidity ^0.4.0; contract Test { function Test() {} }`,
		},
		{
			name:  "years unit",
			input: `pragma solidity ^0.4.0; contract Test { uint256 oneYear = 1 years; }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, &Options{Tolerant: true})
			if err != nil {
				t.Logf("Parse note: %v", err)
			}
			if result == nil || len(result.Children) == 0 {
				t.Error("No AST produced")
			}
		})
	}
}

func TestSolidity05x(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "constructor keyword",
			input: `pragma solidity ^0.5.0; contract Test { constructor() public {} }`,
		},
		{
			name:  "address payable",
			input: `pragma solidity ^0.5.0; contract Test { address payable public owner; }`,
		},
		{
			name:  "emit keyword",
			input: `pragma solidity ^0.5.0; contract Test { event E(); function f() public { emit E(); } }`,
		},
		{
			name:  "calldata location",
			input: `pragma solidity ^0.5.0; contract Test { function f(bytes calldata d) external pure {} }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if result == nil || len(result.Children) == 0 {
				t.Fatal("No AST produced")
			}
		})
	}
}

func TestSolidity06x(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "abstract contract",
			input: `pragma solidity ^0.6.0; abstract contract A { function f() public virtual; }`,
		},
		{
			name:  "virtual and override",
			input: `pragma solidity ^0.6.0; contract B { function f() public virtual {} } contract C is B { function f() public override {} }`,
		},
		{
			name:  "receive function",
			input: `pragma solidity ^0.6.0; contract Test { receive() external payable {} }`,
		},
		{
			name:  "fallback function",
			input: `pragma solidity ^0.6.0; contract Test { fallback() external payable {} }`,
		},
		{
			name:  "try catch",
			input: `pragma solidity ^0.6.0; interface I { function f() external; } contract Test { function t(I i) public { try i.f() {} catch {} } }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if result == nil || len(result.Children) == 0 {
				t.Fatal("No AST produced")
			}
		})
	}
}

func TestSolidity07x(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "gwei denomination",
			input: `pragma solidity ^0.7.0; contract Test { uint256 x = 20 gwei; }`,
		},
		{
			name:  "constructor without visibility",
			input: `pragma solidity ^0.7.0; contract Test { constructor() {} }`,
		},
		{
			name:  "free functions",
			input: `pragma solidity ^0.7.0; function helper(uint256 x) pure returns (uint256) { return x; } contract Test {}`,
		},
		{
			name:  "immutable",
			input: `pragma solidity ^0.7.0; contract Test { uint256 immutable x; constructor() { x = 1; } }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if result == nil || len(result.Children) == 0 {
				t.Fatal("No AST produced")
			}
		})
	}
}

func TestSolidity08x(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unchecked block",
			input: `pragma solidity ^0.8.0; contract Test { function f() public pure returns (uint256) { unchecked { return 1 + 1; } } }`,
		},
		{
			name:  "custom error",
			input: `pragma solidity ^0.8.4; error MyError(uint256 x); contract Test { function f() public { revert MyError(1); } }`,
		},
		{
			name:  "user defined type",
			input: `pragma solidity ^0.8.8; type Price is uint128; contract Test {}`,
		},
		{
			name:  "named mapping",
			input: `pragma solidity ^0.8.18; contract Test { mapping(address account => uint256 balance) public balances; }`,
		},
		{
			name:  "transient storage",
			input: `pragma solidity ^0.8.24; contract Test { uint256 transient x; }`,
		},
		{
			name:  "layout directive",
			input: `pragma solidity ^0.8.24; contract Test layout at 256 { uint256 x; }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input, &Options{Tolerant: true})
			if err != nil {
				t.Logf("Parse note: %v", err)
			}
			if result == nil || len(result.Children) == 0 {
				t.Fatal("No AST produced")
			}
		})
	}
}

// =============================================================================
// Integration Tests with Test Files
// =============================================================================

func TestParseContractFiles(t *testing.T) {
	testFiles := []struct {
		path string
		name string
	}{
		{"../../testdata/contracts/v04/Legacy.sol", "v04 Legacy"},
		{"../../testdata/contracts/v05/Modern.sol", "v05 Modern"},
		{"../../testdata/contracts/v06/Inheritance.sol", "v06 Inheritance"},
		{"../../testdata/contracts/v07/FreeFunctions.sol", "v07 FreeFunctions"},
		{"../../testdata/contracts/v08/SafeMath.sol", "v08 SafeMath"},
		{"../../testdata/contracts/v08/UserDefinedTypes.sol", "v08 UserDefinedTypes"},
		{"../../testdata/contracts/v08/NamedMappings.sol", "v08 NamedMappings"},
		{"../../testdata/contracts/v08/TransientStorage.sol", "v08 TransientStorage"},
		{"../../testdata/contracts/v08/StorageLayout.sol", "v08 StorageLayout"},
		{"../../testdata/contracts/complex/Assembly.sol", "complex Assembly"},
		{"../../testdata/contracts/complex/FullFeatured.sol", "complex FullFeatured"},
		{"../../testdata/contracts/complex/interfaces/IERC20.sol", "complex IERC20"},
		{"../../testdata/contracts/complex/libraries/SafeMath.sol", "complex SafeMath"},
		{"../../testdata/contracts/complex/utils/Helpers.sol", "complex Helpers"},
	}

	for _, tf := range testFiles {
		t.Run(tf.name, func(t *testing.T) {
			absPath, err := filepath.Abs(tf.path)
			if err != nil {
				t.Skipf("Cannot resolve path: %v", err)
			}

			content, err := os.ReadFile(absPath)
			if err != nil {
				t.Skipf("Cannot read file: %v", err)
			}

			result, err := Parse(string(content), nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if result == nil || len(result.Children) == 0 {
				t.Fatal("No AST produced")
			}

			t.Logf("Parsed %d top-level elements", len(result.Children))
		})
	}
}

// =============================================================================
// Complex Feature Tests
// =============================================================================

func TestComplexContract(t *testing.T) {
	input := `
		// SPDX-License-Identifier: MIT
		pragma solidity ^0.8.20;

		error Unauthorized(address caller);
		type TokenId is uint256;

		interface IERC20 {
			function transfer(address to, uint256 amount) external returns (bool);
		}

		abstract contract Ownable {
			address public owner;
			modifier onlyOwner() virtual {
				if (msg.sender != owner) revert Unauthorized(msg.sender);
				_;
			}
		}

		contract Token is Ownable, IERC20 {
			mapping(address account => uint256 balance) public balances;
			uint256 immutable totalSupply;
			uint256 transient processingLock;
			
			event Transfer(address indexed from, address indexed to, uint256 value);
			
			constructor(uint256 _supply) {
				totalSupply = _supply;
				owner = msg.sender;
			}
			
			function transfer(address to, uint256 amount) external override returns (bool) {
				unchecked {
					balances[msg.sender] -= amount;
					balances[to] += amount;
				}
				emit Transfer(msg.sender, to, amount);
				return true;
			}
			
			receive() external payable {}
			fallback() external payable {}
		}
	`

	result, err := Parse(input, &Options{Tolerant: true})
	if err != nil {
		t.Logf("Parse note: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	var counts = make(map[string]int)
	for _, child := range result.Children {
		switch n := child.(type) {
		case *ast.PragmaDirective:
			counts["pragma"]++
		case *ast.ContractDefinition:
			counts[n.Kind]++
		case *ast.ErrorDefinition:
			counts["error"]++
		case *ast.UserDefinedValueTypeDefinition:
			counts["type"]++
		}
	}

	if counts["pragma"] == 0 {
		t.Error("Expected pragma")
	}
	if counts["contract"]+counts["abstract"] == 0 {
		t.Error("Expected contracts")
	}
}

// TestParseAssemblyFlags covers `assembly ("memory-safe") { … }` — the optional
// assemblyFlags group (Solidity >= 0.8.13, pervasive in OpenZeppelin v5 and
// Solady). Before flags were parsed, the unexpected `(` desynchronized the
// tolerant parser and every declaration AFTER the assembly block was silently
// dropped, so this asserts both a clean parse and that later members survive.
func TestParseAssemblyFlags(t *testing.T) {
	tests := []struct {
		name      string
		assembly  string
		wantFlags []string
	}{
		{"no flags", `assembly { let x := 1 }`, nil},
		{"memory safe", `assembly ("memory-safe") { let x := 1 }`, []string{"memory-safe"}},
		{"dialect and flags", `assembly "evmasm" ("memory-safe") { let x := 1 }`, []string{"memory-safe"}},
		{"multiple flags", `assembly ("memory-safe", "other") { let x := 1 }`, []string{"memory-safe", "other"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
contract A {
    function first() public pure { ` + tt.assembly + ` }
    function second() public pure returns (uint256) { return 2; }
    uint256 public third;
}`
			result, errs, err := ParseWithErrors(input, &Options{Tolerant: true})
			if err != nil {
				t.Fatalf("ParseWithErrors: %v", err)
			}
			if len(errs) != 0 {
				t.Fatalf("expected a clean parse, got %d recovered error(s), first: %v", len(errs), errs[0])
			}

			contract := findContract(result, "A")
			if contract == nil {
				t.Fatal("contract A not found")
			}

			// Desync regression: everything declared after the assembly block
			// must still be present.
			var funcs, vars []string
			for _, member := range contract.SubNodes {
				switch n := member.(type) {
				case *ast.FunctionDefinition:
					funcs = append(funcs, n.Name)
				case *ast.StateVariableDeclaration:
					for _, v := range n.Variables {
						vars = append(vars, v.Name)
					}
				}
			}
			if len(funcs) != 2 || funcs[0] != "first" || funcs[1] != "second" {
				t.Errorf("functions = %v, want [first second]", funcs)
			}
			if len(vars) != 1 || vars[0] != "third" {
				t.Errorf("state variables = %v, want [third]", vars)
			}

			asm := findInlineAssembly(contract)
			if asm == nil {
				t.Fatal("inline assembly node not found")
			}
			if asm.Body == nil || len(asm.Body.Operations) == 0 {
				t.Error("assembly body was not parsed")
			}
			if len(asm.Flags) != len(tt.wantFlags) {
				t.Fatalf("Flags = %v, want %v", asm.Flags, tt.wantFlags)
			}
			for i, want := range tt.wantFlags {
				if asm.Flags[i] != want {
					t.Errorf("Flags[%d] = %q, want %q", i, asm.Flags[i], want)
				}
			}
		})
	}
}

// TestParseAssemblyFlagsDialect pins that a dialect string without flags still
// lands in Language rather than being mistaken for a flag.
func TestParseAssemblyFlagsDialect(t *testing.T) {
	input := `pragma solidity ^0.8.20;
contract A { function f() public pure { assembly "evmasm" { let x := 1 } } }`
	result, errs, err := ParseWithErrors(input, &Options{Tolerant: true})
	if err != nil {
		t.Fatalf("ParseWithErrors: %v", err)
	}
	if len(errs) != 0 {
		t.Fatalf("expected a clean parse, got %v", errs[0])
	}
	asm := findInlineAssembly(findContract(result, "A"))
	if asm == nil {
		t.Fatal("inline assembly node not found")
	}
	if asm.Language != "evmasm" {
		t.Errorf("Language = %q, want %q", asm.Language, "evmasm")
	}
	if len(asm.Flags) != 0 {
		t.Errorf("Flags = %v, want none", asm.Flags)
	}
}

func findContract(unit *ast.SourceUnit, name string) *ast.ContractDefinition {
	if unit == nil {
		return nil
	}
	for _, child := range unit.Children {
		if contract, ok := child.(*ast.ContractDefinition); ok && contract.Name == name {
			return contract
		}
	}
	return nil
}

// findInlineAssembly returns the first inline-assembly node in any function body
// of contract, without depending on the Visitor interface.
func findInlineAssembly(contract *ast.ContractDefinition) *ast.InlineAssembly {
	if contract == nil {
		return nil
	}
	for _, member := range contract.SubNodes {
		fn, ok := member.(*ast.FunctionDefinition)
		if !ok || fn.Body == nil {
			continue
		}
		for _, stmt := range fn.Body.Statements {
			if asm, ok := stmt.(*ast.InlineAssembly); ok {
				return asm
			}
		}
	}
	return nil
}

// TestParseAssemblyDottedPath covers Yul paths that contain dots. Only dot-free
// identifiers can be DECLARED inside assembly, but a path may REFER to a
// declaration outside the block — calldata slice members (`sig.offset`,
// `sig.length`) and storage-pointer members (`x.slot`, `x.offset`). Grammar:
// yulPath: (YulIdentifier|YulEVMBuiltin) (YulPeriod (YulIdentifier|YulEVMBuiltin))*.
// Unparsed, the '.' desynchronized the block and shredded the rest of the file
// (OpenZeppelin v5 ECDSA.sol/EIP712.sol).
func TestParseAssemblyDottedPath(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"calldata slice offset", `r := calldataload(sig.offset)`},
		{"nested in builtin", `s := calldataload(add(sig.offset, 0x20))`},
		{"path as assignment target", `store.slot := 1`},
		{"multi assign with path", `a, store.slot := f()`},
		{"deep path", `r := mload(a.b.c)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := `pragma solidity ^0.8.20;
contract A {
    function first(bytes calldata sig) external pure returns (bytes32 r, bytes32 s, uint256 a) {
        assembly ("memory-safe") {
            ` + tt.body + `
        }
    }
    function second() public pure returns (uint256) { return 2; }
    uint256 public third;
}`
			result, errs, err := ParseWithErrors(input, &Options{Tolerant: true})
			if err != nil {
				t.Fatalf("ParseWithErrors: %v", err)
			}
			if len(errs) != 0 {
				t.Fatalf("expected a clean parse, got %d recovered error(s), first: %v", len(errs), errs[0])
			}

			contract := findContract(result, "A")
			if contract == nil {
				t.Fatal("contract A not found")
			}
			// Desync regression: members after the assembly block must survive.
			var funcs, vars []string
			for _, member := range contract.SubNodes {
				switch n := member.(type) {
				case *ast.FunctionDefinition:
					funcs = append(funcs, n.Name)
				case *ast.StateVariableDeclaration:
					for _, v := range n.Variables {
						vars = append(vars, v.Name)
					}
				}
			}
			if len(funcs) != 2 || funcs[1] != "second" {
				t.Errorf("functions = %v, want [first second]", funcs)
			}
			if len(vars) != 1 || vars[0] != "third" {
				t.Errorf("state variables = %v, want [third]", vars)
			}
			if asm := findInlineAssembly(contract); asm == nil || asm.Body == nil || len(asm.Body.Operations) == 0 {
				t.Error("assembly body was not parsed")
			}
		})
	}
}

// TestParsePrefixedStringLiterals covers the `hex"..."` and `unicode"..."`
// literal prefixes. HEX_STRING and UNICODE_STRING token types existed and
// parseStringLiteral already branched on them, but the lexer never produced
// either: `hex` / `unicode` lexed as a bare keyword and the following quote
// became a separate STRING, so parsePrimary hit "expected expression" and
// desynchronized. `hex"19_00"` appears in OpenZeppelin v5 MessageHashUtils.
func TestParsePrefixedStringLiterals(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		wantHex   bool
		wantUni   bool
		wantValue string
		wantParts []string
	}{
		{"hex literal", `hex"1900"`, true, false, "1900", []string{"1900"}},
		{"hex with underscores", `hex"19_00"`, true, false, "19_00", []string{"19_00"}},
		{"empty hex", `hex""`, true, false, "", []string{""}},
		{"concatenated hex", `hex"00" hex"11"`, true, false, "00", []string{"00", "11"}},
		// Real non-ASCII text, which is what the unicode prefix exists for.
		// Backslash-escape decoding inside readString is a separate concern.
		{"unicode literal", `unicode"héllo ☃"`, false, true, "héllo ☃", []string{"héllo ☃"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := `pragma solidity ^0.8.20;
contract A {
    function first() public pure returns (bytes memory) { return abi.encodePacked(` + tt.expr + `); }
    function second() public pure returns (uint256) { return 2; }
    uint256 public third;
}`
			result, errs, err := ParseWithErrors(input, &Options{Tolerant: true})
			if err != nil {
				t.Fatalf("ParseWithErrors: %v", err)
			}
			if len(errs) != 0 {
				t.Fatalf("expected a clean parse, got %d recovered error(s), first: %v", len(errs), errs[0])
			}

			contract := findContract(result, "A")
			if contract == nil {
				t.Fatal("contract A not found")
			}
			// Desync regression: members after the literal must survive.
			var funcs, vars []string
			for _, member := range contract.SubNodes {
				switch n := member.(type) {
				case *ast.FunctionDefinition:
					funcs = append(funcs, n.Name)
				case *ast.StateVariableDeclaration:
					for _, v := range n.Variables {
						vars = append(vars, v.Name)
					}
				}
			}
			if len(funcs) != 2 || funcs[1] != "second" {
				t.Errorf("functions = %v, want [first second]", funcs)
			}
			if len(vars) != 1 || vars[0] != "third" {
				t.Errorf("state variables = %v, want [third]", vars)
			}

			hex, uni := findLiterals(contract)
			if tt.wantHex {
				if hex == nil {
					t.Fatal("expected a HexLiteral node")
				}
				if hex.Value != tt.wantValue {
					t.Errorf("Value = %q, want %q", hex.Value, tt.wantValue)
				}
				if len(hex.Parts) != len(tt.wantParts) {
					t.Fatalf("Parts = %v, want %v", hex.Parts, tt.wantParts)
				}
				for i, want := range tt.wantParts {
					if hex.Parts[i] != want {
						t.Errorf("Parts[%d] = %q, want %q", i, hex.Parts[i], want)
					}
				}
			}
			if tt.wantUni {
				if uni == nil {
					t.Fatal("expected a StringLiteral node")
				}
				if !uni.IsUnicode {
					t.Error("IsUnicode = false, want true")
				}
				if uni.Value != tt.wantValue {
					t.Errorf("Value = %q, want %q", uni.Value, tt.wantValue)
				}
			}
		})
	}
}

// findLiterals returns the first HexLiteral and StringLiteral found in the
// arguments of any function call in the contract's function bodies.
func findLiterals(contract *ast.ContractDefinition) (*ast.HexLiteral, *ast.StringLiteral) {
	var hexLit *ast.HexLiteral
	var strLit *ast.StringLiteral
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.HexLiteral:
			if hexLit == nil {
				hexLit = v
			}
		case *ast.StringLiteral:
			if strLit == nil {
				strLit = v
			}
		case *ast.FunctionCall:
			for _, arg := range v.Arguments {
				walk(arg)
			}
		case *ast.ReturnStatement:
			if v.Expression != nil {
				walk(v.Expression)
			}
		}
	}
	for _, member := range contract.SubNodes {
		fn, ok := member.(*ast.FunctionDefinition)
		if !ok || fn.Body == nil {
			continue
		}
		for _, stmt := range fn.Body.Statements {
			walk(stmt)
		}
	}
	return hexLit, strLit
}
