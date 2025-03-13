package main

import (
	"fmt"
	"log"

	"github.com/AndreiMartynenko/dex-trader-x/config"
	"github.com/AndreiMartynenko/dex-trader-x/exchange"
)

func main() {
	// Load API keys and symbol list
	config.LoadConfig()

	// Initialize Binance API
	exchange.InitBinance()

	// Fetch selected prices
	prices, err := exchange.FetchPricesForSelectedSymbols()
	if err != nil {
		log.Fatal("❌ Error fetching Binance prices:", err)
	}

	fmt.Println("\n✅ Final Binance Prices:")
	for symbol, price := range prices {
		fmt.Printf("   %s: %.2f USDT\n", symbol, price)
	}
}
