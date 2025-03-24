const { ethers } = require("hardhat");

async function main() {
  const wallet = new ethers.Wallet(process.env.WALLET_PRIVATE_KEY, ethers.provider);

  console.log("Deploying with address:", wallet.address); // Should output MetaMask address

  const balance = await ethers.provider.getBalance(wallet.address);
  console.log("Account balance:", ethers.utils.formatEther(balance), "ETH");

  if (balance.lt(ethers.utils.parseEther("0.01"))) { // Adjust based on estimated gas fees
    throw new Error("Insufficient funds for deployment. Please ensure your wallet is funded.");
  }

  const FlashArbitrage = await ethers.getContractFactory("FlashArbitrage", wallet);
  const contract = await FlashArbitrage.deploy();

  await contract.deployed();
  console.log("✅ Contract successfully deployed at:", contract.target);
}

main().catch((error) => {
  console.error("Error deploying contract:", error);
  process.exitCode = 1;
});
