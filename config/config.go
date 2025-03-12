package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BinanceAPIKey    string
	BinanceAPISecret string
	InfuraAPIKey     string
	WalletPrivateKey string
)

// LoadConfig reads API keys from the .env file
func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	BinanceAPIKey = os.Getenv("BINANCE_API_KEY")
	BinanceAPISecret = os.Getenv("BINANCE_API_SECRET")
	InfuraAPIKey = os.Getenv("INFURA_API_KEY")
	WalletPrivateKey = os.Getenv("WALLET_PRIVATE_KEY")
}
