# dex-trader-x
# DExTraderX – Decentralized Arbitrage Trading Bot

**DExTraderX** is a decentralized arbitrage trading bot built using **Solidity** and **Golang** that automatically detects and executes profitable token price differences across two decentralized exchanges (DEXes), such as **Uniswap** and **Sushiswap**. It uses smart contracts to perform token swaps and a Golang backend to handle logic, logging, and automation.

---

## Tech Stack

| Layer        | Technology                             |
|--------------|-----------------------------------------|
| Smart Contract | Solidity (v0.8.x), OpenZeppelin        |
| Dev Tools    | Hardhat, Ethers.js, dotenv              |
| Backend      | Golang (Go-Ethereum + Web3 integration) |
| Deployment   | EVM Testnets (Goerli, Sepolia)          |
| Monitoring   | Telegram Bot (planned), JSON logs       |
| UI (optional) | React + Ethers.js (future roadmap)     |

---

## Key Features

- Smart contract-based arbitrage execution
- Compare token prices between two routers (Uniswap V2 & Sushiswap)
- Gas-efficient flash arbitrage logic
- Golang backend to monitor price differences in real-time
- Automated execution when profit exceeds threshold
- Supports ERC20 token pairs
- Easily extendable to support more DEXs and pairs
- Built-in reentrancy guard & safety checks

---

## Architecture Diagram

![DeXTraderX drawio](https://github.com/user-attachments/assets/cf4b5695-0e22-45ea-a1b8-4e0b0e72766c)

