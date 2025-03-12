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

func InitBinance() {
	binanceClient = binance.NewClient(config.BinanceAPIKey, config.BinanceAPISecret)
}

func GetBinancePrice(symbol string) (float64, error) {
	prices, err := binanceClient.NewListPricesService().Do(context.Background())
	if err != nil {
		log.Println("Error fetching Binance prices:", err)
		return 0, err
	}

	for _, p := range prices {
		if p.Symbol == symbol {
			fmt.Println("Binance Price:", p.Price)
			price, err := strconv.ParseFloat(p.Price, 64)
			if err != nil {
				return 0, err
			}
			return price, nil
		}
	}
	return 0, fmt.Errorf("Symbol not found")
}
