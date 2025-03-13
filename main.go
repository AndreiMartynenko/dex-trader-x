package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/AndreiMartynenko/dex-trader-x/config"
	"github.com/AndreiMartynenko/dex-trader-x/exchange"
)

func main() {
	// Load API keys and trading pairs
	config.LoadConfig()

	// Initialize Binance client
	exchange.InitBinance()

	// Fetch Binance Prices
	fmt.Println("\n🚀 Fetching Prices from Binance...")
	binancePrices, err := exchange.FetchPricesForSelectedSymbols()
	if err != nil {
		log.Println("❌ Error fetching Binance prices:", err)
	} else {
		// Use a map to track printed symbols to avoid duplicates
		printedSymbols := make(map[string]bool)
		for symbol, price := range binancePrices {
			if !printedSymbols[symbol] {
				fmt.Printf("✅ Binance Price for %s: %.2f USDT\n", symbol, price)
				printedSymbols[symbol] = true
			}
		}
	}

	// Fetch Uniswap Prices
	fmt.Println("\n🚀 Fetching Prices from Uniswap V3...")
	for _, pair := range config.UniswapPairs {
		pair = strings.ReplaceAll(pair, "/", "")
		fmt.Printf("\n🔍 Fetching %s Price from Uniswap...\n", pair)

		sqrtPriceX96, err := exchange.FetchSqrtPriceX96(pair)
		if err != nil {
			log.Printf("❌ Error fetching sqrtPriceX96 for %s: %v", pair, err)
			continue
		}

		ethPrice := exchange.ConvertSqrtPriceToETHPrice(sqrtPriceX96, pair)
		fmt.Printf("✅ Uniswap Price for %s: %.6f USDT\n", pair, ethPrice)
	}
}
