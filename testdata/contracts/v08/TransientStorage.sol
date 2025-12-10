// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title Transient Storage (EIP-1153, new in 0.8.24)

contract ReentrancyGuard {
    // Transient storage - cleared after each transaction
    uint256 transient locked;
    
    error ReentrancyDetected();
    
    modifier nonReentrant() {
        if (locked != 0) {
            revert ReentrancyDetected();
        }
        locked = 1;
        _;
        locked = 0;
    }
    
    function protectedFunction() public nonReentrant {
        // This function is protected from reentrancy
    }
}

contract TransientCounter {
    // Transient variables are cleared after transaction
    uint256 transient temporaryCount;
    uint256 transient temporarySum;
    
    // Regular storage for comparison
    uint256 public persistentCount;
    
    function processItems(uint256[] calldata items) public returns (uint256) {
        temporaryCount = 0;
        temporarySum = 0;
        
        for (uint256 i = 0; i < items.length; i++) {
            temporaryCount++;
            temporarySum += items[i];
        }
        
        persistentCount += temporaryCount;
        return temporarySum;
    }
}

