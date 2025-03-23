const { ethers } = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();

  const flashArb = await ethers.getContract("FlashArbitrage");

  const router1 = "0x..."; // Uniswap
  const router2 = "0x..."; // SushiSwap
  const token1 = "0x..."; // WETH
  const token2 = "0x..."; // USDT

  const amountIn = ethers.utils.parseEther("1"); // 1 WETH
  const minProfit = ethers.utils.parseUnits("1", 6); // 1 USDT

  // Approve contract to spend token1
  const tokenContract = await ethers.getContractAt("IERC20", token1);
  await tokenContract.approve(flashArb.address, amountIn);

  // Execute simulated arbitrage
  const tx = await flashArb.executeTrade(
    router1,
    router2,
    token1,
    token2,
    amountIn,
    minProfit
  );

  console.log("✅ Transaction submitted:", tx.hash);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
