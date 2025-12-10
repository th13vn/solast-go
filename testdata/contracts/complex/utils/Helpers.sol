// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title Helper Utilities

// File-level constants
uint256 constant MAX_UINT256 = type(uint256).max;
uint256 constant PERCENTAGE_BASE = 10000;

// File-level free functions
function toWei(uint256 amount) pure returns (uint256) {
    return amount * 1 ether;
}

function toGwei(uint256 amount) pure returns (uint256) {
    return amount * 1 gwei;
}

function isContract(address account) view returns (bool) {
    uint256 size;
    assembly {
        size := extcodesize(account)
    }
    return size > 0;
}

/// @title Address utilities
library AddressUtils {
    error CallFailed();
    error InsufficientBalance(uint256 balance, uint256 required);
    
    function isContractAddress(address account) internal view returns (bool) {
        return account.code.length > 0;
    }
    
    function sendValue(address payable recipient, uint256 amount) internal {
        if (address(this).balance < amount) {
            revert InsufficientBalance(address(this).balance, amount);
        }
        
        (bool success, ) = recipient.call{value: amount}("");
        if (!success) {
            revert CallFailed();
        }
    }
}

/// @title Math utilities
library MathUtils {
    function min(uint256 a, uint256 b) internal pure returns (uint256) {
        return a < b ? a : b;
    }
    
    function max(uint256 a, uint256 b) internal pure returns (uint256) {
        return a > b ? a : b;
    }
    
    function average(uint256 a, uint256 b) internal pure returns (uint256) {
        return (a & b) + (a ^ b) / 2;
    }
}
