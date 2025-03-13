package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	BinanceAPIKey    string
	BinanceAPISecret string
	InfuraAPIKey     string
	WalletPrivateKey string

	// List of trading pairs to fetch dynamically from Binance
	Symbols      []string
	UniswapPairs []string // List of Uniswap trading pairs
)

// LoadConfig reads API keys and trading pairs from .env
func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	BinanceAPIKey = os.Getenv("BINANCE_API_KEY")
	BinanceAPISecret = os.Getenv("BINANCE_API_SECRET")
	InfuraAPIKey = os.Getenv("INFURA_API_KEY")
	WalletPrivateKey = os.Getenv("WALLET_PRIVATE_KEY")

	// Read Binance trading pairs
	symbolsEnv := os.Getenv("SYMBOLS")
	if symbolsEnv != "" {
		Symbols = strings.Split(symbolsEnv, ",")
	} else {
		Symbols = []string{"BTCUSDT", "ETHUSDT"} // Default Binance symbols
	}

	// Read Uniswap trading pairs
	uniswapEnv := os.Getenv("UNISWAP_PAIRS")
	if uniswapEnv != "" {
		UniswapPairs = strings.Split(uniswapEnv, ",")
	} else {
		UniswapPairs = []string{"ETH/USDC"} // Default Uniswap pair
	}
}
