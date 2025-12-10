// SPDX-License-Identifier: MIT
pragma solidity ^0.4.0;

/// @title Legacy Solidity 0.4.x Contract
contract Legacy {
    // State variables
    uint256 public value;
    address public owner;

    // Time units
    uint256 public oneWeek = 1 weeks;
    uint256 public oneDay = 1 days;
    uint256 public oneHour = 1 hours;

    // Wei units
    uint256 public oneWei = 1 wei;
    uint256 public oneEther = 1 ether;

    // Basic arrays
    uint256[] public dynamicArray;
    uint256[10] public staticArray;

    // Mappings
    mapping(address => uint256) public balances;

    // Events
    event ValueUpdated(uint256 oldValue, uint256 newValue);
    event Transfer(address from, address to, uint256 amount);

    // Modifiers
    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }

    // Constructor (0.4.x style - function with contract name)
    function Legacy() public {
        owner = msg.sender;
    }

    // Structs
    struct Person {
        string name;
        uint256 age;
        address wallet;
    }

    mapping(uint256 => Person) public people;

    // Enums
    enum Status {
        Pending,
        Active,
        Closed
    }
    Status public status;

    // Functions with explicit visibility
    function getValue() public view returns (uint256) {
        return value;
    }

    function setValue(uint256 newValue) public {
        value = newValue;
    }

    function setBalance(address account, uint256 amount) public {
        balances[account] = amount;
    }

    function setStatus(Status newStatus) public {
        status = newStatus;
    }
}
