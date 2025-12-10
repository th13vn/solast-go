// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title ERC20 Interface
interface IERC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address to, uint256 amount) external returns (bool);
    function allowance(
        address tokenOwner,
        address spender
    ) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(
        address sender,
        address to,
        uint256 amount
    ) external returns (bool);

    event Transfer(address indexed sender, address indexed to, uint256 value);
    event Approval(
        address indexed tokenOwner,
        address indexed spender,
        uint256 value
    );
}

/// @title ERC20 Metadata Interface
interface IERC20Metadata is IERC20 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
}
