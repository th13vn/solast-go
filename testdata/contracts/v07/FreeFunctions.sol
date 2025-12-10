// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

/// @title Solidity 0.7.x Contract

// Free functions (functions at file level, new in 0.7)
function add(uint256 a, uint256 b) pure returns (uint256) {
    return a + b;
}

function sub(uint256 a, uint256 b) pure returns (uint256) {
    require(b <= a, "Underflow");
    return a - b;
}

// Library
library SafeMath {
    function safeAdd(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a, "Overflow");
        return c;
    }
}

// Using directive at file level (new in 0.7)
using SafeMath for uint256;

contract Calculator {
    uint256 public result;
    
    // Constructor without visibility (visibility removed in 0.7)
    constructor() {
        result = 0;
    }
    
    // Using free functions
    function addNumbers(uint256 a, uint256 b) public {
        result = add(a, b);
    }
    
    // gwei denomination (finney/szabo removed in 0.7)
    uint256 public gasPrice = 20 gwei;
    uint256 public minValue = 1 wei;
    uint256 public oneEther = 1 ether;
    
    // Exponentiation is right-associative now
    function power() public pure returns (uint256) {
        return 2 ** 3 ** 2;
    }
}

// Immutable variables
contract ImmutableExample {
    uint256 public immutable maxSupply;
    address public immutable owner;
    
    constructor(uint256 supply) {
        maxSupply = supply;
        owner = msg.sender;
    }
    
    function isOwner(address addr) public view returns (bool) {
        return addr == owner;
    }
}

// Inheritance
contract Base {
    uint256 internal x;
}

contract Derived is Base {
    function setX(uint256 newX) public {
        x = newX;
    }
}
