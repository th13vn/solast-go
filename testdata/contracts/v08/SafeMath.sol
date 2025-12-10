// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title Solidity 0.8.x Contract

// Custom errors (new in 0.8.4)
error InsufficientBalance(uint256 available, uint256 required);
error Unauthorized(address caller);
error ZeroAddress();

contract Token {
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;
    
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    
    address public owner;
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    
    constructor(string memory tokenName, string memory tokenSymbol, uint256 initialSupply) {
        name = tokenName;
        symbol = tokenSymbol;
        decimals = 18;
        totalSupply = initialSupply;
        balanceOf[msg.sender] = initialSupply;
        owner = msg.sender;
    }
    
    modifier onlyOwner() {
        if (msg.sender != owner) {
            revert Unauthorized(msg.sender);
        }
        _;
    }
    
    function transfer(address to, uint256 amount) public returns (bool) {
        if (to == address(0)) {
            revert ZeroAddress();
        }
        if (balanceOf[msg.sender] < amount) {
            revert InsufficientBalance(balanceOf[msg.sender], amount);
        }
        
        // Built-in overflow checking in 0.8
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        
        emit Transfer(msg.sender, to, amount);
        return true;
    }
    
    // Unchecked blocks for intentional wrapping (new in 0.8)
    function uncheckedTransfer(address to, uint256 amount) public returns (bool) {
        if (balanceOf[msg.sender] < amount) {
            revert InsufficientBalance(balanceOf[msg.sender], amount);
        }
        
        unchecked {
            balanceOf[msg.sender] -= amount;
            balanceOf[to] += amount;
        }
        
        emit Transfer(msg.sender, to, amount);
        return true;
    }
    
    function approve(address spender, uint256 amount) public returns (bool) {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }
}

// bytes.concat (new in 0.8)
contract BytesOperations {
    function concatenate(bytes memory a, bytes memory b) public pure returns (bytes memory) {
        return bytes.concat(a, b);
    }
    
    function concatStrings(string memory a, string memory b) public pure returns (string memory) {
        return string.concat(a, b);
    }
}

// Interface inheritance
interface IERC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address to, uint256 amount) external returns (bool);
}

interface IERC20Metadata is IERC20 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
}
