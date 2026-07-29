// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title Inline Assembly Examples

contract AssemblyExamples {
    // Basic assembly
    function addAssembly(
        uint256 a,
        uint256 b
    ) public pure returns (uint256 result) {
        assembly {
            result := add(a, b)
        }
    }

    // Memory operations
    function memoryOps() public pure returns (bytes32) {
        bytes32 result;
        assembly {
            let ptr := mload(0x40)
            mstore(ptr, 0x1234)
            mstore(0x40, add(ptr, 0x20))
            result := mload(ptr)
        }
        return result;
    }

    // Control flow in assembly
    function maxValue(
        uint256 a,
        uint256 b
    ) public pure returns (uint256 result) {
        assembly {
            switch gt(a, b)
            case 1 {
                result := a
            }
            default {
                result := b
            }
        }
    }

    // Loops in assembly
    function sumArray(
        uint256[] memory arr
    ) public pure returns (uint256 result) {
        assembly {
            let len := mload(arr)
            let data := add(arr, 0x20)

            for {
                let i := 0
            } lt(i, len) {
                i := add(i, 1)
            } {
                result := add(result, mload(add(data, mul(i, 0x20))))
            }
        }
    }

    // Bitwise operations
    function bitwiseOps(
        uint256 a,
        uint256 b
    )
        public
        pure
        returns (uint256 andResult, uint256 orResult, uint256 xorResult)
    {
        assembly {
            andResult := and(a, b)
            orResult := or(a, b)
            xorResult := xor(a, b)
        }
    }

    // Hash in assembly
    function hashInAssembly(
        bytes memory data
    ) public pure returns (bytes32 hashValue) {
        assembly {
            hashValue := keccak256(add(data, 0x20), mload(data))
        }
    }

    // If statement in assembly
    function ifExample(uint256 a) public pure returns (uint256 result) {
        assembly {
            if gt(a, 10) {
                result := 1
            }
        }
    }

    // Assembly flags (Solidity >= 0.8.13), as used throughout OpenZeppelin v5
    // and Solady. The flag group must be consumed, otherwise the parser desyncs
    // and silently drops every member declared below it.
    function memorySafe(uint256 a) public pure returns (uint256 result) {
        assembly ("memory-safe") {
            result := add(a, 1)
        }
    }

    // Dialect string plus a flag group.
    function dialectAndFlags(uint256 a) public pure returns (uint256 result) {
        assembly "evmasm" ("memory-safe") {
            result := add(a, 2)
        }
    }

    // Present to prove nothing after a flagged assembly block is lost.
    uint256 public afterFlaggedAssembly = 7;
}
