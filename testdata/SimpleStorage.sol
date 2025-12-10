// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleStorage {
    uint256 private value;
    address public owner;
    mapping(address => uint256) public balances;
    
    event ValueChanged(uint256 indexed oldValue, uint256 indexed newValue);
    
    error NotOwner(address caller);
    
    struct Record {
        uint256 timestamp;
        uint256 value;
    }
    
    enum Status { Pending, Active, Completed }
    
    modifier onlyOwner() {
        if (msg.sender != owner) revert NotOwner(msg.sender);
        _;
    }
    
    constructor(uint256 initialValue) {
        value = initialValue;
        owner = msg.sender;
    }
    
    function setValue(uint256 newValue) public onlyOwner {
        uint256 oldValue = value;
        value = newValue;
        emit ValueChanged(oldValue, newValue);
    }
    
    function getValue() public view returns (uint256) {
        return value;
    }
}

interface ISimpleStorage {
    function setValue(uint256 newValue) external;
    function getValue() external view returns (uint256);
}

library MathLib {
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        return a + b;
    }
}
