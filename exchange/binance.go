package exchange

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/AndreiMartynenko/dex-trader-x/config"
	"github.com/adshao/go-binance/v2"
)

var binanceClient *binance.Client

// InitBinance initializes the Binance API client
func InitBinance() {
	fmt.Println("✅ Using Binance API Key:", config.BinanceAPIKey[:6]+"********")
	binanceClient = binance.NewClient(config.BinanceAPIKey, config.BinanceAPISecret)
}

// FetchPricesForSelectedSymbols fetches prices for predefined symbols from Binance
func FetchPricesForSelectedSymbols() (map[string]float64, error) {
	if binanceClient == nil {
		log.Println("❌ Binance client is not initialized. Call InitBinance() first.")
		return nil, fmt.Errorf("binance client not initialized")
	}

	priceMap := make(map[string]float64)

	// Fetch all market prices
	prices, err := binanceClient.NewListPricesService().Do(context.Background())
	if err != nil {
		log.Println("❌ Error fetching Binance prices:", err)
		return nil, err
	}

	fmt.Println("🚀 Fetching Prices for Selected Symbols from Binance...\n")

	// Loop through all Binance prices and match with selected symbols
	for _, p := range prices {
		for _, symbol := range config.Symbols {
			if p.Symbol == symbol {
				price, _ := strconv.ParseFloat(p.Price, 64)
				priceMap[symbol] = price
				// fmt.Printf("✅ Binance Price for %s: %.2f USDT\n", symbol, price)
			}
		}
	}

	if len(priceMap) == 0 {
		fmt.Println("❌ No matching symbols found in Binance API response!")
	}

	return priceMap, nil
}
