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
	binanceClient = binance.NewClient(config.BinanceAPIKey, config.BinanceAPISecret)
}

// GetBinancePrice fetches the latest price for a given symbol (e.g., BTCUSDT)
func GetBinancePrice(symbol string) (float64, error) {
	prices, err := binanceClient.NewListPricesService().Do(context.Background())
	if err != nil {
		log.Println("Error fetching Binance prices:", err)
		return 0, err
	}

	for _, p := range prices {
		if p.Symbol == symbol {
			price, _ := strconv.ParseFloat(p.Price, 64)
			fmt.Println("🔹 Binance Price:", price)
			return price, nil
		}
	}
	return 0, fmt.Errorf("Symbol not found")
}
