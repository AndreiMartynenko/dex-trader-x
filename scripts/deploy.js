async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("🚀 Deploying contracts with account:", deployer.address);

    const balance = await ethers.provider.getBalance(deployer.address);
    console.log("💰 Account balance:", ethers.formatEther(balance), "ETH");

    const Arbitrage = await ethers.getContractFactory("FlashArbitrage");
    const arbitrage = await Arbitrage.deploy();

    await arbitrage.waitForDeployment();

    console.log("✅ FlashArbitrage deployed to:", arbitrage.target);
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
