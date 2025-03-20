package main

import (
	"fmt"
	"log"

	"dex-trader-x/config"
	"dex-trader-x/pairs"
	"dex-trader-x/reserves"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to Ethereum
	fmt.Println("🚀 Connecting to Ethereum Mainnet...")
	client, err := ethclient.Dial(config.InfuraURL)
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
	// Fetch SushiSwap pairs (currently unused, but kept for future use)
	// Fetch SushiSwap pairs (currently unused, but kept for future use)
	_, err = pairs.FetchPairs(config.SushiSwapFactory, client, "SushiSwap")
	if err != nil {
		log.Printf("Error fetching SushiSwap pairs: %v", err)
	}

	// Fetch reserves for pairs
	fmt.Println("\nFetching Liquidity Reserves...")
	for _, pair := range uniswapPairs {
		reserves.GetReserves(pair, client)
	}
}
