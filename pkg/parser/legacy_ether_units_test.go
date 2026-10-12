package parser

import "testing"

// TestLegacyEtherUnits guards the pre-0.7 ether subdenominations `szabo` and
// `finney`. Verified contracts compiled with Solidity < 0.7 still use them
// (MUSystem, pragma =0.6.4: `if (userAmt >= 100 * 1 finney)`). The unit list
// lacked both, so the identifier after the number was an unexpected token and
// the rest of the file desynchronized: one `finney` produced all 195
// parse-recovery diagnostics of that contract.
func TestLegacyEtherUnits(t *testing.T) {
	for _, unit := range []string{"szabo", "finney"} {
		t.Run(unit, func(t *testing.T) {
			src := `pragma solidity ^0.6.4;
contract C {
    function f(uint256 amount) public pure returns (bool) {
        return amount >= 100 * 1 ` + unit + `;
    }
    function afterF() public pure returns (uint256) { return 7; }
}`
			root := decodeAST(t, src)
			if !hasNode(root, "NumberLiteral", "subdenomination", unit) {
				t.Fatalf("no NumberLiteral with subdenomination %q in the AST", unit)
			}
			requireAfterFunction(t, root)
		})
	}
}
