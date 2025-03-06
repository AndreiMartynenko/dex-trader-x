package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Global Variables
var (
	InfuraURL    string
	PrivateKey   string
	ContractAddr string
	SwapRouter   string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	InfuraURL = os.Getenv("INFURA_URL")
	PrivateKey = os.Getenv("PRIVATE_KEY")
	ContractAddr = os.Getenv("CONTRACT_ADDRESS")
	SwapRouter = os.Getenv("SWAP_ROUTER")
}
