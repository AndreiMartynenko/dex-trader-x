// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;


// Task 1: Write a Simple Contract
// Objective:

// Write a Solidity contract called DexTraderX.sol
// Implement:
// A function to deposit ETH.
// A function to withdraw ETH.
// A mapping to store balances.
// Events for deposit/withdrawals.

contract DexTraderX {
    address public owner;
    mapping(address => uint256) public balances; // User balances

    uint256 public ownerBalance; // Accumulated fees

    // Events
    event DepositETH(address indexed user, uint256 amount);
    event WithdrawETH(address indexed user, uint256 amount);
    event FeeCollected(address indexed user, uint256 feeAmount);
    event OwnerWithdrawal(uint256 amount);

    // ✅ Modifier for owner-only functions
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner allowed");
        _;
    }

    constructor() {
        owner = msg.sender; // Set contract deployer as owner
    }

    // ✅ Deposit ETH
    function depositETH() external payable {
        require(msg.value > 0, "Must send ETH");
        balances[msg.sender] += msg.value;
        emit DepositETH(msg.sender, msg.value);
    }

    // ✅ Withdraw ETH with 0.1% fee
    function withdrawETH(uint256 amount) external {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        uint256 fee = (amount * 10) / 10000; // 0.1% fee
        uint256 amountAfterFee = amount - fee;

        balances[msg.sender] -= amount;
        ownerBalance += fee; // Store fee separately

        payable(msg.sender).transfer(amountAfterFee);

        emit WithdrawETH(msg.sender, amountAfterFee);
        emit FeeCollected(msg.sender, fee);
    }

    // ✅ Get user's balance
    function getUserBalance() external view returns (uint256) {
        return balances[msg.sender];
    }

    // ✅ Get owner's accumulated fees (corrected function name)
    function getOwnerBalance() external view onlyOwner returns (uint256) {
        return ownerBalance;
    }

    // ✅ Owner withdraws collected fees (only fees, not user deposits)
    function withdrawFees() external onlyOwner {
        require(ownerBalance > 0, "No fees available");
        
        uint256 amount = ownerBalance;
        ownerBalance = 0; // Reset fee balance
        
        payable(owner).transfer(amount);
        emit OwnerWithdrawal(amount);
    }

    // ✅ Allow contract to receive ETH
    receive() external payable {
        balances[msg.sender] += msg.value;
        emit DepositETH(msg.sender, msg.value);
    }
}



// msg.sender → Who sent the transaction
// msg.value → Amount of ETH sent
// address(this).balance → ETH balance of contract



// require(condition, "Error Message") → Checks if condition is met.
// payable(msg.sender).transfer(amount); → Transfers ETH to sender.





