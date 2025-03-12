package main

import (
	"fmt"

	"github.com/AndreiMartynenko/dex-trader-x/config"
)

func main() {
	config.LoadConfig()

	// Print API keys to check if they load (DO NOT DO THIS IN PRODUCTION)
	fmt.Println("Binance API Key:", config.BinanceAPIKey)
	fmt.Println("Infura API Key:", config.InfuraAPIKey)
	fmt.Println("Wallet Private Key:", config.WalletPrivateKey[:10]+"********") // Hide most of the private key
}
