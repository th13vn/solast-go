// SPDX-License-Identifier: MIT
pragma solidity ^0.8.8;

/// @title User Defined Value Types (new in 0.8.8)
/// @notice Demonstrates type-safe wrappers around primitive types

// User-defined value types
type Price is uint128;
type Quantity is uint128;
type TokenId is uint256;
type Percentage is uint16;

// Using directive for user-defined types
using {add as +, unwrap} for Price global;
using {mul as *} for Quantity global;

function add(Price a, Price b) pure returns (Price) {
    return Price.wrap(Price.unwrap(a) + Price.unwrap(b));
}

function unwrap(Price p) pure returns (uint128) {
    return Price.unwrap(p);
}

function mul(Quantity a, Quantity b) pure returns (Quantity) {
    return Quantity.wrap(Quantity.unwrap(a) * Quantity.unwrap(b));
}

contract Marketplace {
    struct Item {
        TokenId id;
        Price price;
        Quantity available;
        address seller;
    }
    
    mapping(TokenId => Item) public items;
    TokenId private nextTokenId;
    
    error ItemNotFound(TokenId id);
    error InsufficientQuantity(Quantity available, Quantity requested);
    error InvalidPrice();
    
    event ItemListed(TokenId indexed id, Price price, Quantity quantity);
    event ItemSold(TokenId indexed id, address buyer, Quantity quantity);
    
    constructor() {
        nextTokenId = TokenId.wrap(1);
    }
    
    function listItem(Price price, Quantity quantity) public returns (TokenId) {
        if (Price.unwrap(price) == 0) {
            revert InvalidPrice();
        }
        
        TokenId id = nextTokenId;
        items[id] = Item({
            id: id,
            price: price,
            available: quantity,
            seller: msg.sender
        });
        
        nextTokenId = TokenId.wrap(TokenId.unwrap(nextTokenId) + 1);
        
        emit ItemListed(id, price, quantity);
        return id;
    }
    
    function buyItem(TokenId id, Quantity quantity) public payable {
        Item storage item = items[id];
        
        if (item.seller == address(0)) {
            revert ItemNotFound(id);
        }
        if (Quantity.unwrap(item.available) < Quantity.unwrap(quantity)) {
            revert InsufficientQuantity(item.available, quantity);
        }
        
        item.available = Quantity.wrap(
            Quantity.unwrap(item.available) - Quantity.unwrap(quantity)
        );
        
        emit ItemSold(id, msg.sender, quantity);
    }
    
    function getItemPrice(TokenId id) public view returns (Price) {
        return items[id].price;
    }
    
    // Using operator overloading
    function addPrices(Price a, Price b) public pure returns (Price) {
        return a + b;  // Uses the overloaded + operator
    }
}

// Percentage calculations with type safety
library PercentageMath {
    uint16 constant PERCENTAGE_FACTOR = 10000;  // 100.00%
    
    function percentage(uint256 value, Percentage pct) internal pure returns (uint256) {
        return (value * Percentage.unwrap(pct)) / PERCENTAGE_FACTOR;
    }
    
    function fromBasisPoints(uint16 bps) internal pure returns (Percentage) {
        return Percentage.wrap(bps);
    }
}

contract FeeCalculator {
    using PercentageMath for uint256;
    
    Percentage public feeRate;
    
    constructor(uint16 _feeRateBps) {
        feeRate = PercentageMath.fromBasisPoints(_feeRateBps);
    }
    
    function calculateFee(uint256 amount) public view returns (uint256) {
        return amount.percentage(feeRate);
    }
}

