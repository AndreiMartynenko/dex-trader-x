const { ethers } = require("hardhat");

async function main() {
  // Manually create and connect your wallet
  const wallet = new ethers.Wallet(process.env.WALLET_PRIVATE_KEY, ethers.provider);

  // Log the wallet address being used
  console.log("Deploying with address:", wallet.address); // Should output your MetaMask address

  // Fetch the wallet balance
  const balance = await ethers.provider.getBalance(wallet.address);
  console.log("Account balance:", ethers.utils.formatEther(balance), "ETH");

  // Ensure sufficient funds for deployment
  if (balance.lt(ethers.utils.parseEther("0.01"))) { // Adjust based on estimated gas fees
    throw new Error("❌ Insufficient funds for deployment. Please ensure your wallet is funded.");
  }

  // Deploy the contract using your manually created wallet
  const FlashArbitrage = await ethers.getContractFactory("FlashArbitrage", wallet);
  const contract = await FlashArbitrage.deploy();

  await contract.deployed();
  console.log("✅ Contract successfully deployed at:", contract.target);
}

// Global error handling
main().catch((error) => {
  console.error("❌ Error deploying contract:", error);
  process.exitCode = 1;
});
