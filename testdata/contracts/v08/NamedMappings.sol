// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

/// @title Named Parameters in Mappings (new in 0.8.18)

contract NamedMappings {
    // Named parameters in mappings
    mapping(address account => uint256 balance) public balances;
    mapping(address owner => mapping(address spender => uint256 allowance)) public allowances;
    mapping(uint256 tokenId => address owner) public tokenOwners;
    
    event Deposit(address indexed account, uint256 amount);
    event Withdrawal(address indexed account, uint256 amount);
    event Approval(address indexed owner, address indexed spender, uint256 amount);
    
    function deposit() public payable {
        balances[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }
    
    function withdraw(uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        balances[msg.sender] -= amount;
        payable(msg.sender).transfer(amount);
        emit Withdrawal(msg.sender, amount);
    }
    
    function approve(address spender, uint256 amount) public {
        allowances[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
    }
    
    function setTokenOwner(uint256 tokenId, address owner) public {
        tokenOwners[tokenId] = owner;
    }
}
