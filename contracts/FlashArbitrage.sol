// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@uniswap/v2-core/contracts/interfaces/IUniswapV2Pair.sol";
import "@uniswap/v2-periphery/contracts/interfaces/IUniswapV2Router02.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract FlashArbitrage is ReentrancyGuard, Ownable {
    event ArbitrageExecuted(
        address indexed token1,
        address indexed token2,
        uint256 profit,
        uint256 timestamp
    );

    // Execute arbitrage between two routers
    function executeTrade(
        address router1,
        address router2,
        address token1,
        address token2,
        uint256 amountIn,
        uint256 minProfit // ✅ Set minimum acceptable profit to avoid sandwich attacks
    ) external onlyOwner nonReentrant {
        require(amountIn > 0, "Invalid amount");

        // Transfer token1 from sender to this contract
        IERC20(token1).transferFrom(msg.sender, address(this), amountIn);
        IERC20(token1).approve(router1, amountIn);

        address[] memory path = new address[](2);
        path[0] = token1;
        path[1] = token2;

        // Step 1: Swap token1 → token2 on router1
        uint256[] memory amountsOut1 = IUniswapV2Router02(router1).swapExactTokensForTokens(
            amountIn,
            0, // Accept any amountOut for now
            path,
            address(this),
            block.timestamp
        );

        uint256 token2Received = amountsOut1[1];
        require(token2Received > 0, "Swap 1 failed");

        IERC20(token2).approve(router2, token2Received);

        // Step 2: Swap token2 → token1 on router2
        address[] memory reversePath = new address[](2);
        reversePath[0] = token2;
        reversePath[1] = token1;

        uint256[] memory amountsOut2 = IUniswapV2Router02(router2).swapExactTokensForTokens(
            token2Received,
            0, // Accept any amountOut for now
            reversePath,
            address(this),
            block.timestamp
        );

        uint256 finalToken1Balance = amountsOut2[1];
        require(finalToken1Balance > amountIn, "No Profit!");

        uint256 profit = finalToken1Balance - amountIn;
        require(profit >= minProfit, "Profit too low");

        // ✅ Transfer profit to the owner
        IERC20(token1).transfer(owner(), finalToken1Balance);

        emit ArbitrageExecuted(token1, token2, profit, block.timestamp);
    }

    // ✅ Emergency withdrawal
    function withdrawTokens(address token) external onlyOwner {
        uint256 balance = IERC20(token).balanceOf(address(this));
        require(balance > 0, "Nothing to withdraw");
        IERC20(token).transfer(owner(), balance);
    }

    function rescueETH() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }

    receive() external payable {}
}
