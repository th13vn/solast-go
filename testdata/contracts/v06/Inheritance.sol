// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

/// @title Solidity 0.6.x Contract

// Abstract contract (new keyword in 0.6)
abstract contract AbstractBase {
    uint256 internal value;
    
    // Virtual function
    function getValue() public view virtual returns (uint256) {
        return value;
    }
    
    // Pure virtual function
    function calculate(uint256 a, uint256 b) public pure virtual returns (uint256);
}

// Interface
interface ICalculator {
    function add(uint256 a, uint256 b) external pure returns (uint256);
    function sub(uint256 a, uint256 b) external pure returns (uint256);
}

// Concrete implementation
contract Calculator is AbstractBase, ICalculator {
    event Calculated(string operation, uint256 result);
    
    constructor() public {
        value = 100;
    }
    
    // Override virtual function
    function getValue() public view override returns (uint256) {
        return value * 2;
    }
    
    // Implement pure virtual
    function calculate(uint256 a, uint256 b) public pure override returns (uint256) {
        return a * b;
    }
    
    // Implement interface
    function add(uint256 a, uint256 b) external pure override returns (uint256) {
        return a + b;
    }
    
    function sub(uint256 a, uint256 b) external pure override returns (uint256) {
        return a - b;
    }
}

// Contract with receive and fallback
contract PaymentReceiver {
    event Received(address sender, uint256 amount);
    event FallbackCalled(address sender, bytes data);
    
    // receive() for plain Ether transfers (new in 0.6)
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }
    
    // New fallback syntax
    fallback() external payable {
        emit FallbackCalled(msg.sender, msg.data);
    }
    
    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}

// Try/catch example
interface IExternalContract {
    function riskyOperation() external returns (uint256);
}

contract TryCatchExample {
    event Success(uint256 result);
    event Failed(string reason);
    
    function tryExternalCall(IExternalContract ext) public {
        try ext.riskyOperation() returns (uint256 result) {
            emit Success(result);
        } catch Error(string memory reason) {
            emit Failed(reason);
        } catch {
            emit Failed("Unknown error");
        }
    }
}

// Array slicing (new in 0.6)
contract ArraySlicing {
    function slice(bytes calldata data, uint256 start, uint256 end) 
        external 
        pure 
        returns (bytes calldata) 
    {
        return data[start:end];
    }
}
