package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var (
	AlchemyURL       string
	InfuraURL        string
	UniswapFactory   string
	SushiSwapFactory string
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env file not found, using system environment variables.")
	}

	AlchemyURL = os.Getenv("ALCHEMY_URL")
	InfuraURL = os.Getenv("INFURA_URL")

	// Load Factory Addresses
	UniswapFactory = os.Getenv("UNISWAP_FACTORY")
	SushiSwapFactory = os.Getenv("SUSHISWAP_FACTORY")

	if UniswapFactory == "" || SushiSwapFactory == "" {
		log.Fatal("🚨 Missing UNISWAP_FACTORY or SUSHISWAP_FACTORY in environment variables!")
	}

	if AlchemyURL == "" && InfuraURL == "" {
		log.Fatal("🚨 Missing ALCHEMY_URL and INFURA_URL in environment variables. At least one is required!")
	}
	return nil
}

// GetEthereumClient tries Alchemy first, then Infura if Alchemy fails
func GetEthereumClient() (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error

	// Try Alchemy First
	if AlchemyURL != "" {
		client, err = ethclient.Dial(AlchemyURL)
		if err == nil {
			fmt.Println("✅ Connected to Ethereum via Alchemy")
			return client, nil
		}
		log.Println("⚠️ Alchemy connection failed. Trying Infura...")
	}

	// Try Infura If Alchemy Fails
	if InfuraURL != "" {
		client, err = ethclient.Dial(InfuraURL)
		if err == nil {
			fmt.Println("✅ Connected to Ethereum via Infura")
			return client, nil
		}
		log.Println("❌ Both Alchemy and Infura failed to connect!")
	}

	return nil, fmt.Errorf("unable to connect to Ethereum via Alchemy or Infura")
}
