package main

import (
	"dex-trader-x/arbitrage"
	"dex-trader-x/pairs"
	"dex-trader-x/reserves"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Fetch Ethereum RPC URL from environment variable
	ethRPC := os.Getenv("ALCHEMY_URL")
	if ethRPC == "" {
		log.Fatal("❌ Missing ALCHEMY_URL in .env file")
	}

	fmt.Println("🚀 Connecting to Ethereum...")

	client, err := ethclient.Dial(ethRPC)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Ethereum: %v", err)
	}
	fmt.Println("✅ Connected to Ethereum!")

	// Fetch Uniswap and SushiSwap pairs
	fmt.Println("\nFetching Uniswap & SushiSwap Pairs...")
	uniswapPairs := pairs.FetchUniswapPairs(client)
	sushiswapPairs := pairs.FetchSushiSwapPairs(client)

	// Find common pairs
	commonPairs := pairs.FindCommonPairs(uniswapPairs, sushiswapPairs)
	if len(commonPairs) == 0 {
		fmt.Println("❌ No common pairs found. Exiting.")
		return
	}

	// Fetch liquidity reserves
	fmt.Println("\nFetching Liquidity Reserves...")
	reserves.FetchReserves(commonPairs, client)

	// Check for arbitrage opportunities
	fmt.Println("\n🔍 Checking Arbitrage Opportunities...")
	arbitrage.CheckArbitrageOpportunities(commonPairs, client, uniswapPairs, sushiswapPairs)
}
