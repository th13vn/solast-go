// SPDX-License-Identifier: MIT
pragma solidity ^0.5.0;

/// @title Solidity 0.5.x Contract
contract Modern {
    // Explicit visibility required
    uint256 public value;
    address payable public owner;
    
    // Events
    event ValueChanged(address indexed changer, uint256 oldValue, uint256 newValue);
    event Received(address indexed sender, uint256 amount);
    
    // Constructor keyword (required in 0.5+)
    constructor() public {
        owner = msg.sender;
        value = 0;
    }
    
    // Explicit visibility required for functions
    function setValue(uint256 newValue) public {
        uint256 oldValue = value;
        value = newValue;
        emit ValueChanged(msg.sender, oldValue, newValue);
    }
    
    // view modifier (replaces constant)
    function getValue() public view returns (uint256) {
        return value;
    }
    
    // pure modifier for no state access
    function calculate(uint256 a, uint256 b) public pure returns (uint256) {
        return a + b;
    }
    
    // address payable for receiving ether
    function withdraw() public {
        require(msg.sender == owner, "Not owner");
        owner.transfer(address(this).balance);
    }
    
    // calldata for external function parameters
    function processData(bytes calldata data) external pure returns (uint256) {
        return data.length;
    }
    
    // Explicit data location required
    function stringParam(string memory str) public pure returns (uint256) {
        return bytes(str).length;
    }
    
    // abi.encode functions
    function encodeExample(uint256 a, address b) public pure returns (bytes memory) {
        return abi.encode(a, b);
    }
}

// Interface with explicit visibility
interface IModern {
    function getValue() external view returns (uint256);
    function setValue(uint256 newValue) external;
}

// Library
library MathLib {
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        return a + b;
    }
    
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a, "Underflow");
        return a - b;
    }
}
