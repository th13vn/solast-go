// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title Full Featured Contract

// Custom errors
error Unauthorized(address caller);
error InsufficientFunds(uint256 available, uint256 required);
error ZeroAddress();

// User-defined value types
type TokenAmount is uint256;
type Timestamp is uint48;

// Free functions
function toTokenAmount(uint256 value) pure returns (TokenAmount) {
    return TokenAmount.wrap(value);
}

// Enums
enum OrderStatus { Pending, Active, Filled, Cancelled }

// Structs
struct Order {
    uint256 id;
    address maker;
    TokenAmount amount;
    OrderStatus status;
}

// Events
event OrderCreated(uint256 indexed orderId, address indexed maker, uint256 amount);
event OrderFilled(uint256 indexed orderId, address indexed taker);

// Interface
interface IOrderBook {
    function createOrder(uint256 amount) external returns (uint256);
    function cancelOrder(uint256 orderId) external;
}

// Abstract contract
abstract contract Ownable {
    address private ownerAddress;
    
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
    
    constructor() {
        ownerAddress = msg.sender;
    }
    
    modifier onlyOwner() virtual {
        if (msg.sender != ownerAddress) {
            revert Unauthorized(msg.sender);
        }
        _;
    }
    
    function owner() public view returns (address) {
        return ownerAddress;
    }
}

// Main contract
contract Exchange is IOrderBook, Ownable {
    mapping(uint256 => Order) public orders;
    mapping(address => TokenAmount) public balances;
    
    uint256 private nextOrderId;
    uint256 public totalVolume;
    
    // Transient storage for reentrancy guard
    uint256 transient locked;
    
    modifier nonReentrant() {
        require(locked == 0, "Reentrant");
        locked = 1;
        _;
        locked = 0;
    }
    
    constructor() {
        nextOrderId = 1;
    }
    
    function deposit() external payable nonReentrant {
        balances[msg.sender] = TokenAmount.wrap(
            TokenAmount.unwrap(balances[msg.sender]) + msg.value
        );
    }
    
    function createOrder(uint256 amount) external override returns (uint256) {
        uint256 orderId = nextOrderId++;
        
        orders[orderId] = Order({
            id: orderId,
            maker: msg.sender,
            amount: toTokenAmount(amount),
            status: OrderStatus.Active
        });
        
        emit OrderCreated(orderId, msg.sender, amount);
        return orderId;
    }
    
    function cancelOrder(uint256 orderId) external override {
        Order storage order = orders[orderId];
        
        if (order.maker != msg.sender) {
            revert Unauthorized(msg.sender);
        }
        
        order.status = OrderStatus.Cancelled;
    }
    
    function fillOrder(uint256 orderId, uint256 amount) external nonReentrant {
        Order storage order = orders[orderId];
        require(order.status == OrderStatus.Active, "Invalid order");
        
        unchecked {
            order.amount = TokenAmount.wrap(TokenAmount.unwrap(order.amount) - amount);
        }
        
        if (TokenAmount.unwrap(order.amount) == 0) {
            order.status = OrderStatus.Filled;
        }
        
        totalVolume += amount;
        emit OrderFilled(orderId, msg.sender);
    }
    
    receive() external payable {}
    fallback() external payable {}
}
