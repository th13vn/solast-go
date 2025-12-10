// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title Storage Layout Directive (new in 0.8.24)
/// @notice Demonstrates explicit storage slot positioning

// Contract with explicit storage layout
contract ExplicitStorage layout at 0 {
    // Variables start at slot 0
    uint256 public value0;      // slot 0
    uint256 public value1;      // slot 1
    address public owner;       // slot 2
    bool public flag;           // slot 2 (packed)
}

// Upgradeable-friendly storage layout
contract StorageV1 layout at 100 {
    // Start at slot 100, leaving room for proxy storage
    uint256 public version;     // slot 100
    address public admin;       // slot 101
    mapping(address => uint256) public balances;  // slot 102
}

// Another contract inheriting storage layout
contract StorageV2 layout at 100 {
    // Same layout as V1 for upgradeability
    uint256 public version;     // slot 100
    address public admin;       // slot 101
    mapping(address => uint256) public balances;  // slot 102
    
    // New variables in V2 continue from slot 103
    uint256 public newFeature;  // slot 103
}

// Diamond proxy compatible storage
contract DiamondStorage layout at 256 {
    // Leave slots 0-255 for diamond/proxy internals
    struct AppStorage {
        uint256 totalSupply;
        mapping(address => uint256) balances;
        mapping(address => mapping(address => uint256)) allowances;
    }
    
    AppStorage internal s;
}

// Predictable storage for cross-contract calls
contract PredictableStorage layout at 1000 {
    bytes32 constant ADMIN_SLOT = bytes32(uint256(1000));
    bytes32 constant IMPLEMENTATION_SLOT = bytes32(uint256(1001));
    
    uint256 public data;
    
    function getAdminSlot() public pure returns (bytes32) {
        return ADMIN_SLOT;
    }
}

