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
