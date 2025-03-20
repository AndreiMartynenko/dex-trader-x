package main

import (
	"fmt"
	"log"

	"dex-trader-x/config"
	"dex-trader-x/pairs"
	"dex-trader-x/reserves"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to Ethereum (Auto-switch between Alchemy & Infura)
	fmt.Println("🚀 Connecting to Ethereum...")
	client, err := config.GetEthereumClient()
	if err != nil {
		log.Fatalf("Error connecting to Ethereum: %v", err)
	}
	defer client.Close()

	// Fetch Uniswap and SushiSwap pairs
	fmt.Println("Fetching Uniswap & SushiSwap Pairs...")
	uniswapPairs, err := pairs.FetchPairs(config.UniswapFactory, client, "Uniswap")
	if err != nil {
		log.Fatalf("Error fetching Uniswap pairs: %v", err)
	}

	_, err = pairs.FetchPairs(config.SushiSwapFactory, client, "SushiSwap")
	if err != nil {
		log.Printf("Error fetching SushiSwap pairs: %v", err)
	}

	// Fetch reserves for pairs
	fmt.Println("\nFetching Liquidity Reserves...")
	for _, pair := range uniswapPairs {
		reserves.GetReserves(pair.Hex(), client)
	}
}
