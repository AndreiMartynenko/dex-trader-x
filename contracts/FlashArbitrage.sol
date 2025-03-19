// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@uniswap/v2-core/contracts/interfaces/IUniswapV2Pair.sol";
import "@uniswap/v2-periphery/contracts/interfaces/IUniswapV2Router02.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract FlashArbitrage {
    address owner;

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner {
        require(msg.sender == owner, "Not owner");
        _;
    }

    // Execute arbitrage
    function executeTrade(
        address router1,
        address router2,
        address token1,
        address token2,
        uint amountIn
    ) external onlyOwner {
        IERC20(token1).transferFrom(msg.sender, address(this), amountIn);
        IERC20(token1).approve(router1, amountIn);

        address[] memory path = new address[](2);
        path[0] = token1;
        path[1] = token2;

        uint[] memory amounts = IUniswapV2Router02(router1).swapExactTokensForTokens(
            amountIn,
            0,
            path,
            address(this),
            block.timestamp
        );

        IERC20(token2).approve(router2, amounts[1]);
        
        address[] memory pathReverse = new address[](2);
        pathReverse[0] = token2;
        pathReverse[1] = token1;

        IUniswapV2Router02(router2).swapExactTokensForTokens(
            amounts[1],
            0,
            pathReverse,
            address(this),
            block.timestamp
        );

        uint finalBalance = IERC20(token1).balanceOf(address(this));
        require(finalBalance > amountIn, "No Profit!");

        IERC20(token1).transfer(owner, finalBalance);
    }
}
